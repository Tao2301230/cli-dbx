package driverprocess

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/linlay/cli-dbx/internal/conn"
	"github.com/linlay/cli-dbx/internal/drivermodule"
)

func TestModuleCommandReusesHostAndKeepsCredentialsOffArguments(t *testing.T) {
	host, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(host)
	if err != nil {
		t.Fatal(err)
	}
	hostHash := sha256.Sum256(data)
	moduleHash := sha256.Sum256([]byte("fixture"))
	path := filepath.Join(t.TempDir(), "postgres.dbmod")
	if err = os.WriteFile(path, []byte("fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	manifest := drivermodule.Manifest{SchemaVersion: 1, Engine: "postgres", GoVersion: runtime.Version(), OS: runtime.GOOS, Arch: runtime.GOARCH, Module: filepath.Base(path), SHA256: hex.EncodeToString(moduleHash[:]), HostSHA256: hex.EncodeToString(hostHash[:])}
	raw, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path+".json", raw, 0600); err != nil {
		t.Fatal(err)
	}
	executable, args, err := command(conn.Spec{Engine: "postgres", Driver: "pgx", Module: true, DriverPath: path, DSN: "postgres://app:private-password@localhost/db"})
	if !drivermodule.Enabled {
		if err == nil {
			t.Fatal("non-module host accepted module")
		}
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	if executable != host || len(args) != 3 || args[0] != "--internal-module-worker" || args[1] != "postgres" || args[2] != path {
		t.Fatalf("unexpected command: %q %q", executable, args)
	}
	if strings.Contains(strings.Join(args, " "), "private-password") {
		t.Fatal("credentials exposed in worker arguments")
	}
}
