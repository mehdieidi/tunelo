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
	defer stop()

	<-ctx.Done()

	fmt.Println("\n[-] shutdown signal received")

	os.Exit(0)
}
