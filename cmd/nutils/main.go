package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cmd := &cobra.Command{
		Use:          "nutils",
		Short:        "Notation Utilities",
		SilenceUsage: true,
	}
	cmd.AddCommand(
		annotationsCmd(),
	)
	if err := cmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
