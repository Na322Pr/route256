package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/cli"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/repository"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/usecase"
)

const psqlDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	facade := repository.NewFacade(pool)
	orderUseCase := usecase.NewOrderUseCase(facade)

	go cli.Run(orderUseCase)

	<-stop
	fmt.Println("\nExiting...")
}
