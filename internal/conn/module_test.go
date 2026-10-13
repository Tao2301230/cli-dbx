package conn

import (
	"context"
	"testing"
)

func TestResolveModuleKeepsDatabaseDriver(t *testing.T) {
	for _, tc := range []struct{ engine, driver string }{{"postgres", "pgx"}, {"oracle", "oracle"}, {"dm", "dm"}, {"sqlserver", "sqlserver"}} {
		t.Run(tc.engine, func(t *testing.T) {
			dir := writeConfig(t, "remote", "[connection]\nengine='"+tc.engine+"'\ndriver='module'\nhost='localhost'\nuser='app'\ndatabase='app'\nallow_actions=['query']\n")
			s, err := Resolve(context.Background(), ResolveInput{ConfigPath: dir, Name: "remote"})
			if err != nil {
				t.Fatal(err)
			}
			if s.Driver != tc.driver {
				t.Fatalf("driver = %q", s.Driver)
			}
			if !s.Module {
				t.Fatal("module selection was lost")
			}
		})
	}
}

func TestResolveRejectsModuleForBuiltinEngine(t *testing.T) {
	dir := writeConfig(t, "local", "[connection]\nengine='sqlite'\ndriver='module'\npath=':memory:'\nallow_actions=['query']\n")
	if _, err := Resolve(context.Background(), ResolveInput{ConfigPath: dir, Name: "local"}); err == nil {
		t.Fatal("accepted module for builtin engine")
	}
}
