package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestModulePolicyPrecedesDiscovery(t *testing.T) {
	dir, _ := driverConfig(t, `"query", "update"`)
	path := filepath.Join(dir, "remote.toml")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	raw = []byte(strings.Replace(string(raw), "[connection]", "[connection]\ndriver = 'module'", 1))
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		args []string
		code string
	}{
		{[]string{"schema", "--config", dir, "remote", "create table users (id int)"}, "action_blocked"},
		{[]string{"update", "--config", dir, "remote", "delete from users"}, "missing_where"},
		{[]string{"query", "--config", dir, "remote", "select * from users; select * from users"}, "multiple_statements_blocked"},
	} {
		p := runDriverCLI(t, tc.args...)
		if p["code"] != tc.code {
			t.Fatalf("expected %s, got %#v", tc.code, p)
		}
	}
}
