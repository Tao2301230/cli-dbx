package main

import _ "github.com/microsoft/go-mssqldb"

// Register is called after the loader runs the driver's package initializers.
func Register() {}

func main() { Register() }
