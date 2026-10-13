//go:build dbmod

// dbmod-pack is an internal build tool. The host, object, and dependency exports
// must originate from the same canonical build (source, SDK, dependencies, and
// flags). Use scripts/build-modules.py to establish that precondition; matching
// Go version and platform alone do not prove package ABI compatibility.
package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"debug/buildinfo"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	runtimeformat "github.com/linlay/cli-dbx/internal/dbmodformat"
	"github.com/pkujhd/goloader"
	"github.com/pkujhd/goloader/link"
)

type manifest struct {
	SchemaVersion int    `json:"schema_version"`
	Engine        string `json:"engine"`
	GoVersion     string `json:"go_version"`
	OS            string `json:"os"`
	Arch          string `json:"arch"`
	HostSHA256    string `json:"host_sha256"`
	SHA256        string `json:"sha256"`
	Module        string `json:"module"`
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil && err != flag.ErrHelp {
		fmt.Fprintln(os.Stderr, "dbmod-pack:", err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("dbmod-pack", flag.ContinueOnError)
	fs.SetOutput(out)
	fs.Usage = func() {
		fmt.Fprintln(out, "dbmod-pack is an internal build tool: host, object, and dependency exports must originate from the same canonical build (source, SDK, dependencies, and flags).")
		fmt.Fprintln(out, "Use scripts/build-modules.py to establish this precondition.")
		fs.PrintDefaults()
	}
	host := fs.String("host", "", "final DBX host executable")
	object := fs.String("object", "", "compiled module object")
	deps := fs.String("deps", "", "dependency exports as import-path TAB archive-path")
	engine := fs.String("engine", "", "postgres, oracle, dm, or sqlserver")
	dir := fs.String("out", "", "output directory")
	if err := fs.Parse(args); err != nil {
		return err
	}
	switch *engine {
	case "postgres", "oracle", "dm", "sqlserver":
	default:
		return fmt.Errorf("unsupported engine %q", *engine)
	}
	if fs.NArg() != 0 || *host == "" || *object == "" || *deps == "" || *dir == "" {
		return fmt.Errorf("required: --host HOST --object OBJECT --deps DEPS --engine ENGINE --out DIR")
	}
	if err := validateHost(*host); err != nil {
		return err
	}
	hostHash, err := hashFile(*host)
	if err != nil {
		return err
	}
	data, err := os.Open(*deps)
	if err != nil {
		return err
	}
	files, pkgs, err := readDependencies(data)
	data.Close()
	if err != nil {
		return err
	}
	symbols := map[string]uintptr{}
	if err := goloader.RegSymbolWithPath(symbols, *host); err != nil {
		return fmt.Errorf("read host symbols (build without -s): %w", err)
	}
	if len(symbols) == 0 {
		return fmt.Errorf("host has no symbols; build without -s")
	}
	hostSymbolCount := len(symbols)
	l, err := goloader.ReadObj(*object, "main")
	if err != nil {
		return fmt.Errorf("read module object: %w", err)
	}
	if l.SymMap["main.Register"] == nil {
		return fmt.Errorf("module has no main.Register entrypoint")
	}
	for pass := 0; ; pass++ {
		unresolved := goloader.UnresolvedSymbols(l, symbols)
		if len(unresolved) == 0 && pass > 0 {
			break
		}
		if pass > len(pkgs) {
			sort.Strings(unresolved)
			return fmt.Errorf("unresolved module symbols: %v", unresolved)
		}
		if err := goloader.ReadDependPackages(l, files, pkgs, unresolved, symbols); err != nil {
			return fmt.Errorf("resolve dependencies: %w", err)
		}
	}
	stats, err := link.Prelink((*link.Linker)(l))
	if err != nil {
		return fmt.Errorf("prelink module: %w", err)
	}
	if after, err := hashFile(*host); err != nil || after != hostHash {
		return fmt.Errorf("host changed during packaging; rebuild modules against the final host")
	}
	if err := writePackage(*host, *engine, *dir, (*link.Linker)(l)); err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(struct {
		Engine      string            `json:"engine"`
		HostSymbols int               `json:"host_symbols"`
		CodeBytes   int               `json:"code_bytes"`
		DataBytes   int               `json:"data_bytes"`
		Prelink     link.PrelinkStats `json:"prelink"`
	}{*engine, hostSymbolCount, len(l.Code), len(l.Data) + len(l.NoPtrData) + len(l.NoPtrTypeData) + len(l.NoPtrItabData), stats})
}

func validateHost(path string) error {
	info, err := buildinfo.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read host Go build information: %w", err)
	}
	if info.GoVersion != runtime.Version() {
		return fmt.Errorf("host Go version %q differs from packer %q", info.GoVersion, runtime.Version())
	}
	settings := map[string]string{}
	for _, value := range info.Settings {
		settings[value.Key] = value.Value
	}
	if settings["GOOS"] != runtime.GOOS || settings["GOARCH"] != runtime.GOARCH {
		return fmt.Errorf("module packaging requires a native host: got %s/%s, want %s/%s", settings["GOOS"], settings["GOARCH"], runtime.GOOS, runtime.GOARCH)
	}
	return nil
}

func readDependencies(r io.Reader) (files, pkgs []string, err error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	seen := map[string]string{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 2 || fields[0] == "" || fields[1] == "" {
			return nil, nil, fmt.Errorf("invalid dependency row %q", line)
		}
		if previous, ok := seen[fields[0]]; ok {
			if previous != fields[1] {
				return nil, nil, fmt.Errorf("conflicting export for %s", fields[0])
			}
			continue
		}
		seen[fields[0]] = fields[1]
		pkgs = append(pkgs, fields[0])
		files = append(files, fields[1])
	}
	return files, pkgs, scanner.Err()
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func writePackage(host, engine, dir string, l *link.Linker) error {
	hostHash, err := hashFile(host)
	if err != nil {
		return err
	}
	var payload bytes.Buffer
	if err := runtimeformat.Encode(l, &payload); err != nil {
		return err
	}
	sum := sha256.Sum256(payload.Bytes())
	name := engine + ".dbmod"
	m := manifest{1, engine, runtime.Version(), runtime.GOOS, runtime.GOARCH, hostHash, hex.EncodeToString(sum[:]), name}
	metadata, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := atomicWrite(filepath.Join(dir, name), payload.Bytes()); err != nil {
		return err
	}
	return atomicWrite(filepath.Join(dir, name+".json"), append(metadata, '\n'))
}

func atomicWrite(path string, data []byte) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".dbmod-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Chmod(0644); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}
