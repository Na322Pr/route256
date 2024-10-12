package order_suite

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/domain"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/repository"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/usecase"
)

type OrderSuite struct {
	suite.Suite
	pool *pgxpool.Pool
	repo usecase.Facade
}

const psqlDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

func (s *OrderSuite) SetupSuite() {
	var err error
	e := exec.Command("make", "compose-up")
	e.Stdout = os.Stdout
	e.Stderr = os.Stderr

	if err = e.Run(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	s.pool, err = pgxpool.New(ctx, psqlDSN)
	s.Require().NoError(err)

	s.repo = repository.NewFacade(s.pool)
}

func (s *OrderSuite) TearDownSuite() {
	s.pool.Close()

	var err error

	e := exec.Command("make", "compose-down")
	e.Stdout = os.Stdout
	e.Stderr = os.Stderr

	if err = e.Run(); err != nil {
		log.Fatal(err)
	}
}

func (s *OrderSuite) SetupTest() {
	var err error
	e := exec.Command("make", "goose-up")
	e.Stdout = os.Stdout
	e.Stderr = os.Stderr

	if err = e.Run(); err != nil {
		log.Fatal(err)
	}
}

func (s *OrderSuite) TearDownTest() {
	var err error

	e := exec.Command("make", "goose-down")
	e.Stdout = os.Stdout
	e.Stderr = os.Stderr

	if err = e.Run(); err != nil {
		log.Fatal(err)
	}
}

func (s *OrderSuite) TestAddOrderSuccess() {
	order := dto.OrderDTO{
		ID:         10,
		ClientID:   10,
		StoreUntil: time.Now().AddDate(0, 0, 2),
		Cost:       1000,
		Weight:     5,
		Packages:   []string{"unknown", "unknown"},
	}

	err := s.repo.AddOrder(context.Background(), order)
	s.Require().NoError(err)
}

func (s *OrderSuite) TestAddOrderFailed() {
	order := dto.OrderDTO{
		ID:         10,
		ClientID:   10,
		StoreUntil: time.Now().AddDate(0, 0, 2),
		Cost:       1000,
		Weight:     5,
		Packages:   []string{"unknown", "unknown"},
	}

	err := s.repo.AddOrder(context.Background(), order)
	s.Require().NoError(err)

	order = dto.OrderDTO{
		ID:         10,
		ClientID:   10,
		StoreUntil: time.Now().AddDate(0, 0, 2),
		Cost:       1000,
		Weight:     5,
		Packages:   []string{"unknown", "unknown"},
	}

	err = s.repo.AddOrder(context.Background(), order)
	s.Require().Error(err)
}

func (s *OrderSuite) TestGetOrderByIDSuccess() {
	order := dto.OrderDTO{
		ID:         10,
		ClientID:   10,
		StoreUntil: time.Now().AddDate(0, 0, 2),
		Cost:       1000,
		Weight:     5,
		Packages:   []string{"unknown", "unknown"},
	}
	err := s.repo.AddOrder(context.Background(), order)
	s.Require().NoError(err)

	_, err = s.repo.GetOrderByID(context.Background(), 10)
	s.Require().NoError(err)
}

func (s *OrderSuite) TestGetOrderByIDFailed() {
	_, err := s.repo.GetOrderByID(context.Background(), 10)
	s.Require().Error(err)
}

func (s *OrderSuite) TestGetOrdersByIDSuccess() {
	order := dto.OrderDTO{
		ID:         10,
		ClientID:   10,
		StoreUntil: time.Now().AddDate(0, 0, 2),
		Cost:       1000,
		Weight:     5,
		Packages:   []string{"unknown", "unknown"},
	}

	err := s.repo.AddOrder(context.Background(), order)
	s.Require().NoError(err)

	order = dto.OrderDTO{
		ID:         11,
		ClientID:   10,
		StoreUntil: time.Now().AddDate(0, 0, 2),
		Cost:       100,
		Weight:     7,
		Packages:   []string{"unknown", "unknown"},
	}

	err = s.repo.AddOrder(context.Background(), order)
	s.Require().NoError(err)

	_, err = s.repo.GetOrdersByIDs(context.Background(), []int64{10, 11})
	s.Require().NoError(err)
}

func (s *OrderSuite) TestGetClientOrdersListSuccess() {
	order := dto.OrderDTO{
		ID:         10,
		ClientID:   10,
		StoreUntil: time.Now().AddDate(0, 0, 2),
		Cost:       1000,
		Weight:     5,
		Packages:   []string{"unknown", "unknown"},
	}

	err := s.repo.AddOrder(context.Background(), order)
	s.Require().NoError(err)

	order = dto.OrderDTO{
		ID:         11,
		ClientID:   10,
		StoreUntil: time.Now().AddDate(0, 0, 2),
		Cost:       100,
		Weight:     7,
		Packages:   []string{"unknown", "unknown"},
	}

	err = s.repo.AddOrder(context.Background(), order)
	s.Require().NoError(err)

	_, err = s.repo.GetClientOrdersList(context.Background(), 10)
	s.Require().NoError(err)
}

func (s *OrderSuite) TestGetRefundsListSuccess() {
	pickUpTime := time.Now()
	order := dto.OrderDTO{
		ID:         10,
		ClientID:   10,
		StoreUntil: time.Now().Add(24 * time.Hour),
		Cost:       1000,
		Weight:     7,
		PickUpTime: sql.NullTime{Time: pickUpTime, Valid: true},
		Status:     domain.OrderStatusMap[domain.OrderStatusRefunded],
	}

	err := s.repo.AddOrder(context.Background(), order)
	s.Require().NoError(err)

	_, err = s.repo.GetRefundsList(context.Background(), 0, 0)
	s.Require().NoError(err)
}
