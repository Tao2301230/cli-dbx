import importlib.util
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest


SCRIPT = Path(__file__).resolve().parents[1] / "build-modules.py"


class BuildModulesTests(unittest.TestCase):
    def load_builder(self):
        self.assertTrue(SCRIPT.is_file(), "module build entrypoint is missing")
        spec = importlib.util.spec_from_file_location("build_modules", SCRIPT)
        module = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(module)
        return module

    def test_prepare_sdk_does_not_modify_installed_sdk_and_reuses_prepared_copy(self):
        builder = self.load_builder()
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            sdk = root / "installed"
            (sdk / "src/cmd/internal/obj").mkdir(parents=True)
            (sdk / "src/cmd/internal/obj/obj.go").write_text("compiler internals")
            (sdk / "VERSION").write_text("go1.26.8\n")
            private = builder.prepare_sdk(sdk, root / "private")
            self.assertFalse((sdk / "src/cmd/objfile").exists())
            self.assertEqual((private / "src/cmd/objfile/obj/obj.go").read_text(), "compiler internals")
            self.assertEqual(builder.prepare_sdk(private, root / "unused"), private)
            self.assertFalse((root / "unused").exists())

    def test_dependency_export_is_sorted_and_rejects_incomplete_or_conflicting_packages(self):
        builder = self.load_builder()
        stream = '\n'.join(json.dumps(p) for p in [
            {"ImportPath": "z", "Export": "/cache/z.a"},
            {"ImportPath": "unsafe"},
            {"ImportPath": "a", "Export": "/cache/a.a"},
        ])
        self.assertEqual(builder.dependency_exports(stream), [("a", "/cache/a.a"), ("z", "/cache/z.a")])
        for packages in [
            [{"ImportPath": "broken", "Error": {"Err": "missing module"}}],
            [{"ImportPath": "broken"}],
            [{"ImportPath": "a", "Export": "/one"}, {"ImportPath": "a", "Export": "/two"}],
        ]:
            with self.subTest(packages=packages), self.assertRaises(ValueError):
                builder.dependency_exports(''.join(json.dumps(p) for p in packages))

    def test_cross_target_requires_host_only(self):
        builder = self.load_builder()
        with self.assertRaisesRegex(ValueError, "host-only"):
            builder.validate_target("windows", "amd64", "darwin", "arm64", False)
        builder.validate_target("windows", "amd64", "darwin", "arm64", True)

    @unittest.skipIf(os.name == "nt", "fixture SDK uses a POSIX executable script")
    def test_existing_host_requires_identical_canonical_build_even_with_same_go_version(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            sdk = root / "sdk"
            (sdk / "bin").mkdir(parents=True)
            (sdk / "src/cmd/objfile/obj").mkdir(parents=True)
            go = sdk / "bin/go"
            canonical = b"go1.26.8 same-toolchain canonical-package-ABI"
            go.write_text(f'''#!{sys.executable}
import json, pathlib, sys
args = sys.argv[1:]
if args[0] == "env":
    print(json.dumps({{"GOROOT": {str(sdk)!r}, "GOVERSION": "go1.26.8", "GOHOSTOS": "darwin", "GOHOSTARCH": "arm64"}}))
elif args[0] == "build" and args[-1] == "./cmd/dbx":
    pathlib.Path(args[args.index("-o") + 1]).write_bytes({canonical!r})
else:
    sys.exit("PACKAGING_REACHED")
''')
            go.chmod(0o755)
            selected = root / "external-dbx"
            for original, rejected in [(b"go1.26.8 same-toolchain different-package-ABI", True), (canonical, False)]:
                with self.subTest(rejected=rejected):
                    selected.write_bytes(original)
                    output = root / ("rejected" if rejected else "accepted")
                    result = subprocess.run([sys.executable, "-B", str(SCRIPT), "--go", str(go),
                                             "--out", str(output), "--host", str(selected), "--engine", "postgres"],
                                            capture_output=True, text=True)
                    if rejected:
                        self.assertIn("differs from the canonical host", result.stderr)
                        self.assertNotIn("PACKAGING_REACHED", result.stderr)
                    else:
                        self.assertIn("PACKAGING_REACHED", result.stderr)
                    canonical_path = output / "_build/canonical-host"
                    self.assertTrue(canonical_path.is_file(), "canonical comparison host was not built")
                    self.assertEqual(canonical_path.read_bytes(), canonical)
                    self.assertEqual(selected.read_bytes(), original, "selected host was modified")
                    self.assertFalse(list(output.glob("*.dbmod")))

            output = root / "overlap"
            selected = output / "_build/canonical-host"
            selected.parent.mkdir(parents=True)
            selected.write_bytes(b"preserve selected host")
            result = subprocess.run([sys.executable, "-B", str(SCRIPT), "--go", str(go),
                                     "--out", str(output), "--host", str(selected)],
                                    capture_output=True, text=True)
            self.assertIn("canonical verification output", result.stderr)
            self.assertEqual(selected.read_bytes(), b"preserve selected host")

    def test_size_report_distinguishes_raw_and_gzip_and_is_reproducible(self):
        builder = self.load_builder()
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            host = root / "dbx"
            host.write_bytes(b"host" * 100)
            module = root / "postgres.dbmod"
            module.write_bytes(b"module" * 100)
            first = builder.size_report(host, [module])
            second = builder.size_report(host, [module])
            self.assertEqual(first, second)
            self.assertEqual(first["modules"][0]["raw_bytes"], 600)
            self.assertLess(first["modules"][0]["gzip_bytes"], 600)
            self.assertEqual(first["host"]["raw_bytes"], 400)


if __name__ == "__main__":
    unittest.main()
