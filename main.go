package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/marchenkosasha2/homework/cmd"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		cmd.Execute()
	}()

	<-stop
	fmt.Println("\nExiting...")
}
