# ODBC as an optional independent process

Add an ODBC route for Oracle, Dameng and SQL Server alongside their native Go adapters. The core continues to link only MySQL and SQLite vendor packages. This extends the user's flat-path experiment; it does not merge the older main-branch ODBC implementation wholesale.

Connections explicitly select `driver = "odbc"`, a versioned ODBC executable via `driver_path`, and the separately installed vendor library via `odbc_driver`. `odbc_options` supplies vendor settings but cannot override credentials, DRIVER/DSN, or structured target fields. Raw ODBC DSNs are rejected. Passwords continue through the existing encrypted-config resolver and private request pipe.

One `drivers/odbc` module links the Go ODBC adapter; separate binaries are built with a fixed engine (oracle/dm/sqlserver). Parent policy and SQL classification run before child startup. Child loads the vendor driver and builds its ODBC connection string. Transactions and imports remain one process/transaction; ODBC uses question-mark parameter markers with the engine's SQL quoting and metadata rules. Redact vendor details conservatively and retain uncertain-write guidance.

ODBC executable names are `dbx-driver-<engine>-odbc-<version-without-v>[.exe]`; manifest names append `.json`. Manifest `adapter:"odbc"` must agree with the connection route; absent/native adapter continues to support existing Go-driver packages. Optional `files` maps flat companion names to SHA-256. Reject escaping paths, symlinks/nonregular resources, missing resources and invalid digests.

Mac ODBC packages bundle copied unixODBC and libltdl libraries, flattened and prefixed with the DBX release version. Rewrite copied load paths to @loader_path and ad-hoc sign them; never alter system libraries. Include flat versioned license/author/notice files. All files have manifest hashes. Runtime files of the same release can be shared by the three packages; different releases use distinct names. Vendor libraries remain separately installed with their original dependencies and licensing; their hardcoded dependencies are not rewritten by DBX.

Windows x64 uses the system ODBC Driver Manager and requires the vendor DLL to be registered for that architecture; DBX does not register DLLs. Linux requires a native CGO build and locally installed unixODBC. Current ODBC release validation covers Mac ARM64 and Windows x64 cross compilation; no Windows ARM64 support is claimed.

Do not touch existing user-owned untracked ODBC installation notes. Add `docs/odbc-driver-process.md`, keep native documentation current, increment release to v0.2.0-driver-process.3. Test core dependency isolation, config/options/redaction, adapter/resource rejection, ODBC placeholder/rollback behavior, flat deterministic packages and actual relocated Mac ODBC calls using a local test vendor. Real vendor compatibility remains subject to live-database validation.
