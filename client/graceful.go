package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func handleGracefulShutdown() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-ctx.Done()

		fmt.Println("\n[-] shutdown signal received")

		defer func() {
			stop()
		}()

		os.Exit(0)
	}()
}
