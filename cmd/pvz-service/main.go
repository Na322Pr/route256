package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/app/mw"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/app/pvz_service"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/config"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/repository"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	desc "gitlab.ozon.dev/marchenkosasha2/homework/pkg/pvz-service/v1"
)

func main() {
	cfg := config.MustLoad()

	psqlDSN := getPsqlDSN(cfg)
	grpcHost := fmt.Sprintf("%s:%s", cfg.GRPC.Host, cfg.GRPC.Port)

	ctx := context.Background()
	ctxWithCancel, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pool, err := pgxpool.New(ctxWithCancel, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	facade := repository.NewFacade(pool)
	orderUseCase := usecase.NewOrderUseCase(facade)
	pvzService := pvz_service.NewImplementation(*orderUseCase)

	lis, err := net.Listen("tcp", grpcHost)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(mw.Logging),
	)
	reflection.Register(grpcServer)
	desc.RegisterPVZServiceServer(grpcServer, pvzService)

	mux := runtime.NewServeMux()
	err = desc.RegisterPVZServiceHandlerFromEndpoint(ctx, mux, grpcHost, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		log.Fatal("failed t register telephone service handler: %w", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

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
