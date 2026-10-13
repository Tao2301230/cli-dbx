//go:build dbmod

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkujhd/goloader/link"
)

func TestHelpExplainsCanonicalBuildPrecondition(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"--help"}, &out); err != flag.ErrHelp {
		t.Fatalf("help error = %v", err)
	}
	for _, phrase := range []string{"internal build tool", "host, object, and dependency exports", "same canonical build", "scripts/build-modules.py"} {
		if !strings.Contains(out.String(), phrase) {
			t.Fatalf("help missing %q: %s", phrase, out.String())
		}
	}
}

func TestInvalidArgumentsLeaveNoArtifacts(t *testing.T) {
	for _, engine := range []string{"", "mysql", "../postgres"} {
		t.Run(engine, func(t *testing.T) {
			out := filepath.Join(t.TempDir(), "output")
			err := run([]string{"--engine", engine, "--host", "missing", "--object", "missing", "--deps", "missing", "--out", out}, &bytes.Buffer{})
			if err == nil || !strings.Contains(err.Error(), "engine") {
				t.Fatalf("expected engine error, got %v", err)
			}
			if _, err := os.Stat(out); !os.IsNotExist(err) {
				t.Fatalf("invalid arguments created output: %v", err)
			}
		})
	}
}

func TestDependencyListRejectsConflictsAndMalformedRows(t *testing.T) {
	for _, input := range []string{"fmt\n", "fmt\t\n", "\t/cache/fmt.a\n", "fmt\t/cache/a\nfmt\t/cache/b\n"} {
		if _, _, err := readDependencies(strings.NewReader(input)); err == nil {
			t.Fatalf("accepted invalid dependency closure %q", input)
		}
	}
	files, pkgs, err := readDependencies(strings.NewReader("fmt\t/cache/fmt.a\n\nerrors\t/cache/errors.a\n"))
	if err != nil || len(files) != 2 || pkgs[0] != "fmt" || pkgs[1] != "errors" {
		t.Fatalf("lost dependency mapping: %v %v %v", files, pkgs, err)
	}
}

func TestPackageBindsExactHostAndPayloadWithoutPathLeak(t *testing.T) {
	root := t.TempDir()
	host := filepath.Join(root, "private-build-host")
	if err := os.WriteFile(host, []byte("host version one"), 0600); err != nil {
		t.Fatal(err)
	}
	var first []byte
	for _, output := range []string{"one", "two"} {
		out := filepath.Join(root, output)
		if err := writePackage(host, "postgres", out, &link.Linker{}); err != nil {
			t.Fatal(err)
		}
		payload, err := os.ReadFile(filepath.Join(out, "postgres.dbmod"))
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(filepath.Join(out, "postgres.dbmod.json"))
		if err != nil {
			t.Fatal(err)
		}
		if first != nil && !bytes.Equal(first, data) {
			t.Fatal("output path changed manifest bytes")
		}
		first = data
		var got manifest
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		hostSum := sha256.Sum256([]byte("host version one"))
		moduleSum := sha256.Sum256(payload)
		if got.HostSHA256 != hex.EncodeToString(hostSum[:]) || got.SHA256 != hex.EncodeToString(moduleSum[:]) {
			t.Fatalf("manifest does not bind exact content: %+v", got)
		}
		if got.Module != "postgres.dbmod" || bytes.Contains(data, []byte(root)) {
			t.Fatalf("manifest leaks local paths: %s", data)
		}
	}
}

func TestMissingHostDoesNotPublishModule(t *testing.T) {
	out := filepath.Join(t.TempDir(), "output")
	if err := writePackage(filepath.Join(t.TempDir(), "missing"), "postgres", out, &link.Linker{}); err == nil {
		t.Fatal("accepted missing host")
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Fatalf("published a module without a host: %v", err)
	}
}
