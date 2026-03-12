// Package main is the entry point for the mycelium daemon.
package main

import (
	"log/slog"
	"os"
)

var version = "dev"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("starting mycelium daemon", "version", version)

	// TODO: implement daemon
}
