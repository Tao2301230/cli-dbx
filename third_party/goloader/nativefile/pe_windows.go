//go:build windows

// Symbol classification follows Go src/cmd/internal/objfile/pe.go.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in LICENSE.go.txt.
package nativefile

import (
	"debug/pe"
	"fmt"
)

type Sym struct {
	Name string
	Addr uint64
	Code rune
}
type File struct{ f *pe.File }

func Open(path string) (*File, error) {
	f, e := pe.Open(path)
	if e != nil {
		return nil, e
	}
	return &File{f}, nil
}
func (f *File) Close() error { return f.f.Close() }
func (f *File) Symbols() ([]Sym, error) {
	var base uint64
	switch h := f.f.OptionalHeader.(type) {
	case *pe.OptionalHeader64:
		base = h.ImageBase
	case *pe.OptionalHeader32:
		base = uint64(h.ImageBase)
	default:
		return nil, fmt.Errorf("unsupported PE header")
	}
	if len(f.f.Symbols) == 0 {
		return nil, fmt.Errorf("host has no symbols; build without -s")
	}
	out := make([]Sym, 0, len(f.f.Symbols))
	for _, s := range f.f.Symbols {
		sym := Sym{Name: s.Name, Addr: uint64(s.Value), Code: '?'}
		switch s.SectionNumber {
		case 0:
			sym.Code = 'U'
		case -1:
			sym.Code = 'C'
		case -2:
		default:
			if s.SectionNumber < 0 || int(s.SectionNumber) > len(f.f.Sections) {
				return nil, fmt.Errorf("invalid PE symbol section")
			}
			section := f.f.Sections[s.SectionNumber-1]
			ch := section.Characteristics
			switch {
			case ch&0x20 != 0:
				sym.Code = 'T'
			case ch&0x40 != 0:
				if ch&0x80000000 == 0 {
					sym.Code = 'R'
				} else {
					sym.Code = 'D'
				}
			case ch&0x80 != 0:
				sym.Code = 'B'
			}
			sym.Addr += base + uint64(section.VirtualAddress)
		}
		out = append(out, sym)
	}
	return out, nil
}
