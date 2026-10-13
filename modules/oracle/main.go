package main

import (
	"fmt"
	"net/url"

	"github.com/linlay/cli-dbx/internal/conn"
	"github.com/oracle/go-oracledb/v26/oracle"
)

// Register is called after the loader runs the driver's package initializers.
func Register() {}

// Prepare retains the existing native driver's structured Oracle connection.
// It constructs a connector; the shared executor opens the database later.
func Prepare(spec *conn.Spec) error {
	u, err := url.Parse(spec.DSN)
	if err != nil || u.Scheme != "oracle" || u.User == nil || u.Host == "" || u.Path == "" {
		return fmt.Errorf("invalid Oracle connection; use structured connection fields")
	}
	cfg := oracle.NewOracleDriverConfig()
	cfg.ConnectDescriptor = "tcp://" + u.Host + u.EscapedPath()
	cfg.Credentials.User = u.User.Username()
	cfg.Credentials.Password, _ = u.User.Password()
	connector, err := oracle.NewOracleConnector(cfg)
	if err != nil {
		return fmt.Errorf("invalid Oracle connector configuration")
	}
	spec.Connector = connector
	return nil
}

func main() { Register() }
