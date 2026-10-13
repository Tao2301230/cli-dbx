//go:build dbmod

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/linlay/cli-dbx/internal/conn"
	"github.com/linlay/cli-dbx/internal/dbcore"
	"github.com/linlay/cli-dbx/internal/drivermodule"
	"github.com/linlay/cli-dbx/internal/driverprocess"
)

func internalWorker(ctx context.Context, args []string) (bool, int) {
	if len(args) == 0 || (args[0] != "--internal-module-worker" && args[0] != "--internal-module-check") {
		return false, 0
	}
	fail := func(err error) (bool, int) { fmt.Fprintln(os.Stderr, err); return true, 1 }
	if len(args) != 3 {
		return fail(fmt.Errorf("invalid internal module arguments"))
	}
	// Vendor initialization must not parse DBX's private worker arguments.
	// Oracle's logger calls flag.Parse when the global set is still unparsed.
	if err := flag.CommandLine.Parse(nil); err != nil {
		return fail(err)
	}
	host, err := os.Executable()
	if err != nil {
		return fail(err)
	}
	v, err := drivermodule.Open(conn.Spec{Engine: args[1], DriverPath: args[2]}, host)
	if err != nil {
		return fail(err)
	}
	defer v.File.Close()
	loaded, err := drivermodule.Load(v)
	if err != nil {
		return fail(err)
	}
	if args[0] == "--internal-module-check" {
		spec := conn.Spec{Engine: args[1], Driver: args[1], DSN: args[1] + "://app:test@127.0.0.1:1/app"}
		if spec.Engine == "postgres" {
			spec.Driver = "pgx"
		}
		if spec.Engine == "dm" {
			spec.DSN = "dm://app:test@127.0.0.1:1"
		}
		if loaded.Prepare != nil {
			if err = loaded.Prepare(&spec); err != nil {
				return fail(err)
			}
		}
		db, e := dbcore.Open(spec)
		if e != nil {
			return fail(e)
		}
		if e = db.Close(); e != nil {
			return fail(e)
		}
		runtime.GC()
		err = json.NewEncoder(os.Stdout).Encode(map[string]any{"engine": args[1], "loaded": true, "drivers": sql.Drivers(), "database_connected": false})
	} else {
		err = driverprocess.Serve(ctx, os.Stdin, os.Stdout, args[1], loaded.Prepare)
	}
	runtime.KeepAlive(loaded)
	if err != nil {
		return fail(err)
	}
	return true, 0
}
