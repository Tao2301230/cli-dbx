//go:build darwin

// Mach-O symbol classification follows Go src/cmd/internal/objfile/macho.go.
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE.go file.
package nativefile
import("debug/macho";"slices";"cmp")
type Sym struct { Name string; Addr uint64; Code rune }
type File struct { f *macho.File }
func Open(path string)(*File,error){ f,e:=macho.Open(path); if e!=nil{return nil,e};return &File{f},nil }
func(f *File)Close()error{return f.f.Close()}
func(f *File)Symbols()([]Sym,error){
 if f.f.Symtab==nil{return nil,nil}
 out:=make([]Sym,0,len(f.f.Symtab.Syms))
 for _,s:=range f.f.Symtab.Syms{
  if s.Type&0xe0!=0{continue}
  v:=Sym{s.Name,s.Value,'?'}
  if s.Sect==0{v.Code='U'} else if int(s.Sect)<=len(f.f.Sections){
   sect:=f.f.Sections[s.Sect-1]
   switch sect.Seg {case "__TEXT","__DATA_CONST":v.Code='R';case "__DATA":v.Code='D'}
   switch sect.Seg+" "+sect.Name {case "__TEXT __text":v.Code='T';case "__DATA __bss","__DATA __noptrbss":v.Code='B'}
  }
  out=append(out,v)
 }
 slices.SortFunc(out,func(a,b Sym)int{return cmp.Compare(a.Addr,b.Addr)}); return out,nil
}
