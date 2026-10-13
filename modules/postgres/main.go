package main

import _ "github.com/jackc/pgx/v5/stdlib"

// Register is called after the loader runs the driver's package initializers.
func Register() {}

func main() { Register() }
