module github.com/linlay/cli-dbx

go 1.26.8

require (
	filippo.io/edwards25519 v1.1.0
	gitee.com/chunanyong/dm v1.8.22
	github.com/BurntSushi/toml v1.4.0
	github.com/danieljoos/wincred v1.2.3
	github.com/dustin/go-humanize v1.0.1
	github.com/go-sql-driver/mysql v1.8.1
	github.com/godbus/dbus/v5 v5.2.2
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9
	github.com/golang-sql/sqlexp v0.1.0
	github.com/golang/snappy v0.0.1
	github.com/google/uuid v1.6.0
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/jackc/pgpassfile v1.0.0
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761
	github.com/jackc/pgx/v5 v5.7.2
	github.com/jackc/puddle/v2 v2.2.2
	github.com/mattn/go-isatty v0.0.20
	github.com/microsoft/go-mssqldb v1.8.0
	github.com/ncruces/go-strftime v0.1.9
	github.com/oracle/go-oracledb/v26 v26.0.1-beta
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec
	github.com/spf13/cobra v1.10.1
	github.com/zalando/go-keyring v0.2.8
	golang.org/x/crypto v0.54.0
	golang.org/x/sync v0.22.0
	golang.org/x/sys v0.47.0
	golang.org/x/term v0.45.0
	golang.org/x/text v0.40.0
	golang.org/x/tools v0.47.0
	modernc.org/gc/v3 v3.0.0-20240107210532-573471604cb6
	modernc.org/libc v1.55.3
	modernc.org/mathutil v1.6.0
	modernc.org/memory v1.8.0
	modernc.org/sqlite v1.34.2
	modernc.org/strutil v1.2.0
	modernc.org/token v1.1.0
)

replace github.com/spf13/cobra => ./third_party/cobra

require github.com/pkujhd/goloader v0.0.0
replace github.com/pkujhd/goloader => ./third_party/goloader
