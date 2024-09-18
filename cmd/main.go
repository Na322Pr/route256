package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/marchenkosasha2/homework/internal/cli"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/repository"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/usecase"
)

const storage_path = "storage/data.json"

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	orderRepository, err := repository.NewOrderRepository(storage_path)
	if err != nil {
		fmt.Printf("error creating new order repository: %s", err)
		os.Exit(1)
		return
	}

	orderUserCase := usecase.NewOrderUseCase(*orderRepository)

	go cli.Run(orderUserCase)

	<-stop
	fmt.Println("\nExiting...")
}
