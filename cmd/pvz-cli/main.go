package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Na322Pr/route256/internal/cli"
	"github.com/Na322Pr/route256/internal/config"
)

func main() {
	cfg := config.MustLoad()
	serviceURL := fmt.Sprintf("http://%s", cfg.HTTP.Host)

	ctx := context.Background()
	ctxWithCancel, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	syncChan := make(chan struct{})
	defer cancel()

	cli := cli.NewCLI(serviceURL)

	go cli.Run(ctxWithCancel, syncChan)

	for range syncChan {
		fmt.Println("All goroutines are done")
	}
	fmt.Println("Exiting...")
	os.Exit(0)
}
