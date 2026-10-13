// Package drivermodule validates independently distributed Go driver modules.
// Modules are trusted native code bound to an exact DBX host, not a sandbox.
package drivermodule

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/linlay/cli-dbx/internal/conn"
)

const MaxModuleBytes = 128 << 20

type Manifest struct {
	SchemaVersion int    `json:"schema_version"`
	Engine        string `json:"engine"`
	GoVersion     string `json:"go_version"`
	OS            string `json:"os"`
	Arch          string `json:"arch"`
	HostSHA256    string `json:"host_sha256"`
	SHA256        string `json:"sha256"`
	Module        string `json:"module"`
}

type Error struct{ Code, Message string }

func (e *Error) Error() string     { return e.Message }
func invalid(message string) error { return &Error{Code: "driver_invalid", Message: message} }

// Verified owns an immutable snapshot of the bytes whose hash was checked.
// Close File to release the snapshot after decoding or parent-side validation.
type Verified struct {
	Path     string
	File     io.ReadCloser
	Manifest Manifest
}

type payloadReader struct{ reader *bytes.Reader }

func (p *payloadReader) Read(b []byte) (int, error) {
	if p.reader == nil {
		return 0, os.ErrClosed
	}
	return p.reader.Read(b)
}
func (p *payloadReader) Close() error { p.reader = nil; return nil }

func modulePath(spec conn.Spec) (string, error) {
	switch spec.Engine {
	case "postgres", "oracle", "dm", "sqlserver":
	default:
		return "", invalid("unsupported module engine")
	}
	path := spec.DriverPath
	if path == "" {
		root := spec.DriverDir
		if root == "" {
			root = os.Getenv("DBX_DRIVER_DIR")
		}
		if root == "" {
			home, e := os.UserHomeDir()
			if e != nil {
				return "", e
			}
			root = filepath.Join(home, ".config", "dbx", "drivers")
		}
		path = filepath.Join(root, spec.Engine+".dbmod")
	}
	return filepath.Abs(path)
}

func openRegular(path string, limit int64) (*os.File, error) {
	info, e := os.Lstat(path)
	if os.IsNotExist(e) {
		return nil, &Error{Code: "driver_missing", Message: fmt.Sprintf("module file or manifest is missing at %s; install the matching module package", path)}
	}
	if e != nil || !info.Mode().IsRegular() {
		return nil, invalid("module files must be regular files, not symbolic links")
	}
	if info.Size() > limit {
		return nil, invalid("module file exceeds its size limit")
	}
	f, e := os.Open(path)
	if e != nil {
		return nil, invalid("cannot open module file")
	}
	opened, e := f.Stat()
	if e != nil || !os.SameFile(info, opened) {
		f.Close()
		return nil, invalid("module file changed during validation")
	}
	return f, nil
}

func checkHash(r io.Reader, digest string) error {
	raw, e := hex.DecodeString(digest)
	if e != nil || len(raw) != sha256.Size || strings.ToLower(digest) != digest {
		return invalid("invalid module SHA-256")
	}
	h := sha256.New()
	if _, e = io.Copy(h, r); e != nil {
		return invalid("cannot read module bytes")
	}
	if hex.EncodeToString(h.Sum(nil)) != digest {
		return &Error{Code: "driver_integrity_failed", Message: "module or host SHA-256 does not match its manifest; rebuild or install the matching module package"}
	}
	return nil
}

// Open verifies metadata, the host and payload before returning a rewound file.
func Open(spec conn.Spec, hostPath string) (*Verified, error) {
	path, e := modulePath(spec)
	if e != nil {
		return nil, e
	}
	if info, err := os.Lstat(filepath.Dir(path)); err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil, &Error{Code: "driver_missing", Message: fmt.Sprintf("module directory is missing or invalid at %s", filepath.Dir(path))}
	}
	mf, e := openRegular(path+".json", 64<<10)
	if e != nil {
		return nil, e
	}
	var m Manifest
	dec := json.NewDecoder(io.LimitReader(mf, (64<<10)+1))
	dec.DisallowUnknownFields()
	e = dec.Decode(&m)
	if e == nil {
		var extra any
		if dec.Decode(&extra) != io.EOF {
			e = fmt.Errorf("trailing manifest data")
		}
	}
	mf.Close()
	if e != nil {
		return nil, invalid("invalid module manifest")
	}
	if m.SchemaVersion != 1 || m.Engine != spec.Engine || m.GoVersion != runtime.Version() || m.OS != runtime.GOOS || m.Arch != runtime.GOARCH {
		return nil, invalid("module schema, engine, toolchain or platform does not match DBX")
	}
	if m.Module != filepath.Base(path) || strings.ContainsAny(m.Module, "/\\:") {
		return nil, invalid("module manifest filename must match the selected file")
	}
	host, e := os.Open(hostPath)
	if e != nil {
		return nil, invalid("cannot verify DBX host")
	}
	e = checkHash(host, m.HostSHA256)
	host.Close()
	if e != nil {
		return nil, e
	}
	f, e := openRegular(path, MaxModuleBytes)
	if e != nil {
		return nil, e
	}
	payload, e := io.ReadAll(io.LimitReader(f, MaxModuleBytes+1))
	f.Close()
	if e != nil {
		return nil, invalid("cannot read module bytes")
	}
	if len(payload) > MaxModuleBytes {
		return nil, invalid("module file exceeds its size limit")
	}
	if e = checkHash(bytes.NewReader(payload), m.SHA256); e != nil {
		return nil, e
	}
	return &Verified{Path: path, File: &payloadReader{reader: bytes.NewReader(payload)}, Manifest: m}, nil
}
