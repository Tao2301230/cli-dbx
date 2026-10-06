# Flat driver installation

The user approved replacing per-engine/version directories with a flat directory of executables and matching JSON sidecars. Each connection selects an exact executable using `connection.driver_path`.

- A packaged executable is `dbx-driver-<engine>-<version-without-v>`; Windows adds `.exe`. The manifest is the complete executable filename plus `.json`, including `.exe.json` on Windows. Archives contain these two files at their root. Existing archive formats and platform-specific package names remain unchanged.
- `driver_path` overrides all driver-directory settings. `~` expands to the user home; relative paths are resolved relative to the connection TOML file. It selects one file, without scanning, latest-version selection, PATH search, or fallback.
- If no `driver_path` is supplied, retain the directory precedence `driver_dir` > `DBX_DRIVER_DIR` > `~/.config/dbx/drivers`, looking for the flat unversioned `dbx-driver-<engine>[.exe]` and its sidecar. Published versioned packages require `driver_path`.
- Validate protocol, engine, platform, version presence, SHA-256, regular files, and manifest executable equality with the chosen basename. A sidecar cannot redirect execution to another program. Reject executable/manifest symlinks and nonregular files. Keep the existing process and transaction protocol unchanged.
- Different versions can coexist and separate connections can select different executables. Missing or invalid selected files fail without falling back to another installed version.
- Old `<engine>/driver.json` installations need reinstallation using the new flat packages; no compatibility search is added to this experimental branch.
- Bump experimental release to `v0.2.0-driver-process.2`; update CLI help and installation documentation. No installer command, ODBC change, automatic download, or unified ZIP conversion is in scope.

Validation: regression tests for configuration paths, explicit path precedence and exact version selection, sibling-manifest binding, invalid files, default flat lookup, process calls/rollback/cancellation, deterministic flat archives and multiple package coexistence. Run root/module tests and rebuild Mac ARM64 and Windows x64 artifacts; smoke test Mac drivers against a closed local port (not a real database compatibility test).
