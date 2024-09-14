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
		fmt.Println(err)
		return
	}

	orderUserCase := usecase.NewOrderUseCase(*orderRepository)

	go func() {
		cli.Run(orderUserCase)
	}()

	<-stop
	fmt.Println("\nExiting...")
}
