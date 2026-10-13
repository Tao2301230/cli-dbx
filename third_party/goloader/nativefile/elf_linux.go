//go:build linux

// Symbol classification follows Go src/cmd/internal/objfile/elf.go.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in LICENSE.go.txt.
package nativefile

import "debug/elf"

type Sym struct {
	Name string
	Addr uint64
	Code rune
}
type File struct{ f *elf.File }

func Open(path string) (*File, error) {
	f, e := elf.Open(path)
	if e != nil {
		return nil, e
	}
	return &File{f}, nil
}
func (f *File) Close() error { return f.f.Close() }
func (f *File) Symbols() ([]Sym, error) {
	symbols, e := f.f.Symbols()
	if e != nil {
		return nil, e
	}
	out := make([]Sym, 0, len(symbols))
	for _, s := range symbols {
		sym := Sym{Name: s.Name, Addr: s.Value, Code: '?'}
		switch s.Section {
		case elf.SHN_UNDEF:
			sym.Code = 'U'
		case elf.SHN_COMMON:
			sym.Code = 'B'
		default:
			if int(s.Section) < len(f.f.Sections) {
				sect := f.f.Sections[s.Section]
				switch sect.Flags & (elf.SHF_WRITE | elf.SHF_ALLOC | elf.SHF_EXECINSTR) {
				case elf.SHF_ALLOC | elf.SHF_EXECINSTR:
					sym.Code = 'T'
				case elf.SHF_ALLOC:
					sym.Code = 'R'
				case elf.SHF_ALLOC | elf.SHF_WRITE:
					sym.Code = 'D'
				}
			}
		}
		out = append(out, sym)
	}
	return out, nil
}
