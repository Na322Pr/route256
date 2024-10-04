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
	ctx := context.Background()
	ctxWithCancel, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	syncChan := make(chan struct{})
	defer cancel()

	pool, err := pgxpool.New(ctxWithCancel, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	facade := repository.NewFacade(pool)
	orderUseCase := usecase.NewOrderUseCase(facade)
	cli := cli.NewCLI(orderUseCase)

	go cli.Run(ctxWithCancel, syncChan)

	for range syncChan {
		fmt.Println("All goroutines are done")
	}
	fmt.Println("Exiting...")
	os.Exit(0)
}
