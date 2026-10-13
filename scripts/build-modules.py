#!/usr/bin/env python3
"""Build the pinned DBX host and raw driver modules without changing a global SDK."""

import argparse
import gzip
import hashlib
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys


ROOT = Path(__file__).resolve().parents[1]
GO_VERSION = "go1.26.8"
ENGINES = ("postgres", "oracle", "dm", "sqlserver")


def prepare_sdk(source, destination):
    """Expose compiler internals in a private SDK, or reuse a prepared SDK."""
    source, destination = Path(source).resolve(), Path(destination).resolve()
    if (source / "src/cmd/objfile/obj").is_dir():
        return source
    if source == destination or source in destination.parents:
        raise ValueError("private SDK must be outside the installed SDK")
    if destination.exists():
        if ((destination / "VERSION").read_bytes() != (source / "VERSION").read_bytes()
                or not (destination / "src/cmd/objfile/obj").is_dir()):
            raise ValueError("private SDK is incomplete or has a different version; remove it and retry")
        return destination
    temporary = destination.with_name(destination.name + ".preparing")
    if temporary.exists():
        raise ValueError(f"unfinished SDK preparation exists: {temporary}")
    try:
        shutil.copytree(source, temporary, symlinks=True)
        shutil.copytree(temporary / "src/cmd/internal", temporary / "src/cmd/objfile")
        temporary.rename(destination)
    except BaseException:
        shutil.rmtree(temporary, ignore_errors=True)
        raise
    return destination


def dependency_exports(stream):
    """Read the concatenated JSON documents emitted by go list -deps -export."""
    decoder = json.JSONDecoder()
    exports = {}
    position = 0
    while position < len(stream):
        while position < len(stream) and stream[position].isspace():
            position += 1
        if position == len(stream):
            break
        package, position = decoder.raw_decode(stream, position)
        name, archive = package.get("ImportPath"), package.get("Export")
        if package.get("Error") or package.get("DepsErrors"):
            raise ValueError(f"dependency export failed for {name}")
        if name == "unsafe":
            continue
        if not name or not archive or any(c in name + archive for c in "\t\r\n"):
            raise ValueError(f"missing or invalid export for {name}")
        if name in exports and exports[name] != archive:
            raise ValueError(f"conflicting dependency exports for {name}")
        exports[name] = archive
    if not exports:
        raise ValueError("empty dependency export closure")
    return sorted(exports.items())


def validate_target(goos, goarch, hostos, hostarch, host_only):
    if (goos, goarch) != (hostos, hostarch) and not host_only:
        raise ValueError("module packaging currently requires a native target; use --host-only for a cross-build")


def size_report(host, modules):
    def sizes(path):
        raw = Path(path).read_bytes()
        compressed = gzip.compress(raw, compresslevel=9, mtime=0)
        return {"file": Path(path).name, "sha256": hashlib.sha256(raw).hexdigest(),
                "raw_bytes": len(raw), "raw_mib": round(len(raw) / 1048576, 6),
                "gzip_bytes": len(compressed), "gzip_mib": round(len(compressed) / 1048576, 6)}
    return {"schema_version": 1, "host": sizes(host), "modules": [sizes(p) for p in modules]}


def command(args, env, capture=False):
    print("+ " + " ".join(str(arg) for arg in args), flush=True)
    return subprocess.run([str(arg) for arg in args], cwd=ROOT, env=env, check=True,
                          text=True, stdout=subprocess.PIPE if capture else None).stdout


def build(args):
    query_env = dict(os.environ, GOTOOLCHAIN="local", GOWORK="off")
    query_env.pop("GOROOT", None)
    info = json.loads(command([args.go, "env", "-json", "GOROOT", "GOVERSION", "GOHOSTOS", "GOHOSTARCH"], query_env, True))
    if info["GOVERSION"] != GO_VERSION:
        raise ValueError(f"Go {GO_VERSION} is required; --go selected {info['GOVERSION']}")
    goos, goarch = args.goos or info["GOHOSTOS"], args.goarch or info["GOHOSTARCH"]
    validate_target(goos, goarch, info["GOHOSTOS"], info["GOHOSTARCH"], args.host_only)
    if args.host and args.host_only:
        raise ValueError("--host and --host-only cannot be combined")
    if args.host and not args.host.is_file():
        raise ValueError(f"host executable does not exist: {args.host}")
    output = args.out.resolve()
    work = output / "_build"
    work.mkdir(parents=True, exist_ok=True)
    sdk = prepare_sdk(Path(info["GOROOT"]), work / GO_VERSION)
    go = sdk / "bin" / ("go.exe" if os.name == "nt" else "go")
    env = dict(os.environ, GOROOT=str(sdk), GOTOOLCHAIN="local", GOWORK="off",
               GOOS=goos, GOARCH=goarch, GOFLAGS="-mod=readonly")
    env["PATH"] = str(sdk / "bin") + os.pathsep + env.get("PATH", "")
    env.setdefault("GOCACHE", str(work / "gocache"))
    # Darwin arm64's MAP_JIT implementation requires the local C toolchain.
    env["CGO_ENABLED"] = "1" if goos == "darwin" and goarch == "arm64" else "0"
    if args.offline:
        env.update(GOPROXY="off", GOSUMDB="off")
    host = args.host.resolve() if args.host else output / ("dbx.exe" if goos == "windows" else "dbx")
    build_target = host
    if args.host:
        if host == work or work in host.parents:
            raise ValueError("--host cannot be inside the canonical verification output directory")
        reserved = {output / "sizes.json"}
        for engine in ENGINES:
            reserved.update((output / f"{engine}.dbmod", output / f"{engine}.dbmod.json"))
        if host in reserved:
            raise ValueError("--host conflicts with a generated module output")
        build_target = work / ("canonical-host.exe" if goos == "windows" else "canonical-host")
    command([go, "build", "-tags=dbmod", "-trimpath", "-buildvcs=false", "-ldflags=-w", "-o", build_target, "./cmd/dbx"], env)
    if args.host and hashlib.sha256(host.read_bytes()).digest() != hashlib.sha256(build_target.read_bytes()).digest():
        raise ValueError("selected --host differs from the canonical host built from this source tree, SDK, dependencies, and flags; rebuild the host or use its matching source tree")
    if args.host_only:
        report = size_report(host, [])
        report.update(go_version=GO_VERSION, os=goos, arch=goarch, validation="build only; not executed")
        (output / "sizes.json").write_text(json.dumps(report, indent=2) + "\n")
        return
    packer = work / ("dbmod-pack.exe" if goos == "windows" else "dbmod-pack")
    command([go, "build", "-tags=dbmod", "-trimpath", "-buildvcs=false", "-ldflags=-w", "-o", packer, "./cmd/dbmod-pack"], env)
    modules = []
    prelink = {}
    selected = set(args.engine or ENGINES)
    for engine in ENGINES:
        if engine not in selected:
            continue
        module_dir = work / engine
        module_dir.mkdir(exist_ok=True)
        stream = command([go, "list", "-deps", "-export", "-json", "-trimpath", "-buildvcs=false", f"./modules/{engine}"], env, True)
        exports = dependency_exports(stream)
        deps, importcfg = module_dir / "deps.tsv", module_dir / "importcfg"
        deps.write_text("".join(f"{name}\t{archive}\n" for name, archive in exports))
        importcfg.write_text("".join(f"packagefile {name}={archive}\n" for name, archive in exports))
        obj = module_dir / "module.o"
        command([go, "tool", "compile", "-p", "main", "-trimpath", f"{ROOT}=>github.com/linlay/cli-dbx",
                 "-importcfg", importcfg, "-o", obj, ROOT / "modules" / engine / "main.go"], env)
        packed = command([packer, "--host", host, "--object", obj, "--deps", deps, "--engine", engine, "--out", output], env, True)
        print(packed, end="")
        prelink[engine] = json.loads(packed)
        modules.append(output / f"{engine}.dbmod")
    report = size_report(host, modules)
    report.update(go_version=GO_VERSION, os=goos, arch=goarch, prelink=prelink)
    (output / "sizes.json").write_text(json.dumps(report, indent=2) + "\n")
    print(json.dumps(report, indent=2))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--go", default="go", help="Go 1.26.8 executable (stock or privately prepared SDK)")
    parser.add_argument("--out", type=Path, default=ROOT / "dist/modules", help="host, raw modules, manifests and sizes.json output directory")
    parser.add_argument("--host", type=Path, help="use an existing host only if byte-identical to a canonical build from this source tree and SDK; leaves the selected file unchanged")
    parser.add_argument("--engine", action="append", choices=ENGINES, help="repeat to select modules; default: all four")
    parser.add_argument("--goos", help="target GOOS; defaults to build host")
    parser.add_argument("--goarch", help="target GOARCH; defaults to build host")
    parser.add_argument("--host-only", action="store_true", help="only build the host, including cross-target compilation")
    parser.add_argument("--offline", action="store_true", help="require all pinned Go dependencies to be cached")
    args = parser.parse_args()
    try:
        build(args)
    except (ValueError, OSError, subprocess.CalledProcessError) as err:
        parser.exit(1, f"build-modules: {err}\n")


if __name__ == "__main__":
    main()
