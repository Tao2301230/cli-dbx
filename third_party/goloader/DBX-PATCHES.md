# DBX loader provenance

Upstream: https://github.com/pkujhd/goloader, commit `2695cfc5350631cc3becc1409162b12a885471b5`.

Imported from the DBX 2026-10-09 source evidence archive. Upstream license is retained in `LICENSE`; the native executable symbol reader retains the Go license in `nativefile/LICENSE.go.txt`.

The archived changes add a small native Mach-O reader (avoiding runtime linkage of Go object-file tooling), refactor loader package methods into free functions to improve dead-code elimination, resolve local relative wrapper references, pin Darwin JIT writes to one OS thread, guard invalid Oracle itab offsets, and rebase interface method slice pointers during offline host type inspection.

This branch pins Go 1.26.8. Compiler-internal imports require a private SDK copy with `src/cmd/internal` copied to `src/cmd/objfile`; no globally installed SDK is modified. The module host keeps runtime metadata alive until process exit. Native modules are trusted code and are not unloadable during an operation.

New offline prelink changes are documented in `link/prelink.go` and tested beside it. A native Windows PE symbol reader is added for host cross-builds; cross-compilation does not establish runtime correctness.
