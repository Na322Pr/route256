package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-chi/chi"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/app/mw"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/app/pvz_service"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/config"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/kafka/event"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/kafka/producer"
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
	httpHost := cfg.HTTP.Host
	grpcHost := cfg.GRPC.Host
	adminHost := cfg.Admin.Host

	ctx := context.Background()
	ctxWithCancel, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	pool, err := pgxpool.New(ctxWithCancel, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	prod, err := producer.NewSyncProducer(cfg.Kafka,
		producer.WithRequiredAcks(sarama.WaitForLocal),
		producer.WithMaxOpenRequests(1),
		producer.WithMaxRetries(5),
		producer.WithRetryBackoff(10*time.Millisecond),
		producer.WithProducerPartitioner(sarama.NewHashPartitioner),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer prod.Close()

	eventLogProd, err := event.NewEventLogProducer(prod, "pvz.events-log", "pvz-service")
	if err != nil {
		log.Fatal(err)
	}

	facade := repository.NewFacade(pool)
	orderUseCase := usecase.NewOrderUseCase(facade, eventLogProd)
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
		log.Fatal("failed to register telephone service handler: %w", err)
	}

	fmt.Println("Starting http server...")
	go func() {
		if err := http.ListenAndServe(httpHost, mux); err != nil {
			log.Fatalf("failed to listen and serve http service: %v", err)
		}
	}()

	fmt.Println("Starting admin server...")
	go func() {
		adminServer := chi.NewMux()
		adminServer.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			b, _ := os.ReadFile("./pkg/pvz-service/v1/pvz_service.swagger.json")
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
		})

		adminServer.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://localhost:7002/swagger.json"),
		))

		if err := http.ListenAndServe(adminHost, adminServer); err != nil {
			log.Fatalf("failed to listen and server admin server: %v", err)
		}
	}()

	fmt.Println("Starting grpc server...")
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to listen and server grpc server: %v", err)
		}
	}()

	<-stop
	fmt.Println("\nShutting down servers...")
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
