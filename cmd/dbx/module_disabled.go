//go:build !dbmod

package main

import "context"

func internalWorker(context.Context, []string) (bool, int) { return false, 0 }
