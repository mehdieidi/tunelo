package graceful

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func ShutdownHandler() {
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
