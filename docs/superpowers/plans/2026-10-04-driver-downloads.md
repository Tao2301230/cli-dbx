# Driver Downloads Implementation Plan

> **For agentic workers:** Use independent agents for vendor acquisition and inline execution for catalog/build orchestration. User explicitly prohibits Git commits, staging, and push for this task.

**Goal:** Produce one flat, versioned download catalog with all independent Go drivers and the previously discussed ODBC packages.

**Architecture:** Existing build scripts output to a temporary stage; validated products and unchanged historical archives enter driver-downloads. Vendor original packages and complete installed bundles preserve their layout. A standard-library Python tool generates and verifies the catalog.

**Tech Stack:** Go cross compilation, Python standard library, vendor official downloads, macOS unixODBC packaging.

## Global Constraints

- Fixed root: cli-dbx/driver-downloads; no empty platform subdirectories.
- Native engines: postgres/oracle/dm/sqlserver. Native matrix: macOS/Linux/Windows × amd64/arm64.
- ODBC adapters: oracle/dm/sqlserver. Vendor catalog: Oracle/SQLServer/DM/OceanBase.
- Preserve existing staged and unstaged changes and all user files. Never commit or stage.
- Keep source, binary, vendor installer and availability statuses distinct. Hash every saved artifact.

### Task 1: Build and retain DBX packages
- [x] Snapshot original HEAD/index/files.
- [x] Build current Go packages with scripts/build-drivers.sh --output-dir <stage>; check package manifests/target headers.
- [x] Build Mac ARM64 and Windows x64 ODBC variants using scripts/build-odbc-drivers.sh --output-dir <stage>.
- [x] Retain historical .1/.2 archives with valid checksum sidecars; preserve prior .3 artifacts as a dated snapshot if bytes differ.
- [x] Check availability of Linux native container/cross compiler; build if available, otherwise ship native build material with precise status.

### Task 2: Collect complete vendor material
- [x] Preserve complete local Oracle 23.26.2, SQLServer 18.7.1.1, OceanBase 2.0.8.3 bundles and licenses.
- [x] Collect official matching Windows/Linux/macOS packages using verified primary source URLs; no global installation.
- [x] Add unavailable DM/platform entries with vendor acquisition instructions; do not claim support from unrelated Go binaries.
- [x] Record versions, download source URLs, signature/hash validation, dependencies and exact vendor file paths.

### Task 3: Generate and verify download catalog
- [x] Implement scripts/release/driver-downloads.py using Python standard library: index, SHA256SUMS, README, full verification.
- [x] Add source snapshot and portable build instructions for Go/ODBC/Oracle toolchain.
- [x] Verify archive paths, DBX manifest checksums/target headers, vendor bundle symlinks and catalog hashes.
- [x] Confirm original HEAD/index/user files unchanged and report saved packages plus remaining platform limits.

## Execution evidence

107 artifacts retained, current native24 + ODBC15, source and porting kits saved. Tests and byte/target checks passed; current MacARM native/ODBC smoke and MacIntel Rosetta checks passed. Windows/Linux live execution and vendor database acceptance are explicitly unverified. Git HEAD/index and all original source bytes preserved except the authorized .gitignore line. See driver-downloads/VERIFICATION.md for exact scope.
