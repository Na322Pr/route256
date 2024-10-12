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
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/config"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/repository"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/usecase"
)

func main() {
	cfg := config.MustLoad()

	psqlDSN := getPsqlDSN(cfg)

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

func getPsqlDSN(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PG.User,
		cfg.PG.Password,
		cfg.PG.Host,
		cfg.PG.Port,
		cfg.PG.DB,
	)
}
