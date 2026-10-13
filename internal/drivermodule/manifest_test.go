package drivermodule

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/linlay/cli-dbx/internal/conn"
)

func fixture(t *testing.T) (conn.Spec, string, Manifest) {
	t.Helper()
	dir := t.TempDir()
	host := filepath.Join(dir, "dbx")
	path := filepath.Join(dir, "postgres.dbmod")
	if err := os.WriteFile(host, []byte("host executable"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("module payload"), 0600); err != nil {
		t.Fatal(err)
	}
	digest := func(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }
	m := Manifest{SchemaVersion: 1, Engine: "postgres", GoVersion: runtime.Version(), OS: runtime.GOOS, Arch: runtime.GOARCH, Module: "postgres.dbmod", SHA256: digest("module payload"), HostSHA256: digest("host executable")}
	writeManifest(t, path, m)
	return conn.Spec{Engine: "postgres", Module: true, DriverDir: dir}, host, m
}
func writeManifest(t *testing.T, path string, m Manifest) {
	t.Helper()
	b, e := json.Marshal(m)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(path+".json", b, 0600); e != nil {
		t.Fatal(e)
	}
}

func TestOpenVerifiedModule(t *testing.T) {
	spec, host, _ := fixture(t)
	v, err := Open(spec, host)
	if err != nil {
		t.Fatal(err)
	}
	defer v.File.Close()
	b, err := io.ReadAll(v.File)
	if err != nil || string(b) != "module payload" {
		t.Fatalf("verified file not rewound: %q %v", b, err)
	}
}

func TestVerifiedPayloadCannotChangeAfterHashing(t *testing.T) {
	spec, host, _ := fixture(t)
	v, err := Open(spec, host)
	if err != nil {
		t.Fatal(err)
	}
	defer v.File.Close()
	if err = os.WriteFile(v.Path, []byte("changed after hashing"), 0600); err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(v.File)
	if err != nil || string(b) != "module payload" {
		t.Fatalf("verified payload changed: %q %v", b, err)
	}
}
func TestRejectIncompatibleModule(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(*Manifest)
	}{
		{"schema", func(m *Manifest) { m.SchemaVersion++ }},
		{"engine", func(m *Manifest) { m.Engine = "dm" }},
		{"toolchain", func(m *Manifest) { m.GoVersion = "go0.0.0" }},
		{"platform", func(m *Manifest) { m.OS = "other" }},
		{"arch", func(m *Manifest) { m.Arch = "other" }},
		{"host", func(m *Manifest) { m.HostSHA256 = m.SHA256 }},
		{"payload", func(m *Manifest) { m.SHA256 = m.HostSHA256 }},
		{"filename", func(m *Manifest) { m.Module = "../postgres.dbmod" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s, h, m := fixture(t)
			tc.change(&m)
			writeManifest(t, filepath.Join(s.DriverDir, "postgres.dbmod"), m)
			if v, e := Open(s, h); e == nil {
				v.File.Close()
				t.Fatal("accepted incompatible module")
			}
		})
	}
}
func TestModuleDiscoveryAndSymlinkRejection(t *testing.T) {
	s, h, _ := fixture(t)
	path := filepath.Join(s.DriverDir, "postgres.dbmod")
	t.Setenv("DBX_DRIVER_DIR", s.DriverDir)
	s.DriverDir = ""
	v, e := Open(s, h)
	if e != nil {
		t.Fatal(e)
	}
	v.File.Close()
	s.DriverPath = path
	t.Setenv("DBX_DRIVER_DIR", t.TempDir())
	v, e = Open(s, h)
	if e != nil {
		t.Fatal(e)
	}
	v.File.Close()
	if e = os.Rename(path, path+".real"); e != nil {
		t.Fatal(e)
	}
	if e = os.Symlink(path+".real", path); e != nil {
		t.Skip(e)
	}
	if v, e = Open(s, h); e == nil {
		v.File.Close()
		t.Fatal("accepted symlink")
	}
}
func TestModuleRejectsManifestTrailingData(t *testing.T) {
	s, h, _ := fixture(t)
	p := filepath.Join(s.DriverDir, "postgres.dbmod.json")
	f, e := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, 0600)
	if e != nil {
		t.Fatal(e)
	}
	f.WriteString("{}")
	f.Close()
	if v, e := Open(s, h); e == nil {
		v.File.Close()
		t.Fatal("accepted trailing JSON")
	}
}
