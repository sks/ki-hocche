package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/sks/kihocche/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	err := cmd.Execute(ctx)
	if err != nil {
		slog.Error("Error executing command", "error", err)
	}
}
