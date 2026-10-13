//go:build dbmod

package drivermodule

import (
	"database/sql"
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"github.com/linlay/cli-dbx/internal/conn"
	format "github.com/linlay/cli-dbx/internal/dbmodformat"
	"github.com/pkujhd/goloader"
)

const Enabled = true

var loadMu sync.Mutex

// Loaded is deliberately never unloaded: runtime and database/sql retain module
// types and function pointers. Workers terminate to reclaim all module memory.
type Loaded struct {
	code    *goloader.CodeModule
	Prepare func(*conn.Spec) error
}

func Load(v *Verified) (loaded *Loaded, err error) {
	loadMu.Lock()
	defer loadMu.Unlock()
	if runtime.Version() != "go1.26.8" {
		return nil, invalid("module host requires Go 1.26.8")
	}
	defer func() {
		if recover() != nil {
			loaded = nil
			err = invalid("module initialization failed")
		}
	}()
	symbols := map[string]uintptr{}
	if err = goloader.RegSymbol(symbols); err != nil {
		return nil, invalid("cannot read host symbols; rebuild DBX without -s")
	}
	registerSharedHostTypes(symbols)
	graph, err := format.Decode(v.File)
	v.File.Close()
	if err != nil {
		return nil, invalid("cannot decode driver module: " + err.Error())
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	code, err := goloader.Load((*goloader.Linker)(graph), symbols)
	if err != nil {
		return nil, invalid("driver module linking failed")
	}
	entry := code.Syms["main.Register"]
	if entry == 0 {
		return nil, invalid("driver module registration entry is missing")
	}
	// Go function values point at a descriptor whose first word is the entry PC.
	descriptor := unsafe.Pointer(&entry)
	register := *(*func())(unsafe.Pointer(&descriptor))
	register()
	loaded = &Loaded{code: code}
	if prepare := code.Syms["main.Prepare"]; prepare != 0 {
		ptr := unsafe.Pointer(&prepare)
		loaded.Prepare = *(*func(*conn.Spec) error)(unsafe.Pointer(&ptr))
	}
	name := v.Manifest.Engine
	if name == "postgres" {
		name = "pgx"
	}
	found := false
	for _, driver := range sql.Drivers() {
		if driver == name {
			found = true
		}
	}
	// Oracle uses its exported connector; other modules register database/sql.
	if !found && !(v.Manifest.Engine == "oracle" && loaded.Prepare != nil) {
		return nil, invalid(fmt.Sprintf("module did not register %s", name))
	}
	runtime.KeepAlive(code)
	return loaded, nil
}
