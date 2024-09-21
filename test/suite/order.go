package order_suite

import (
	"os"
	"time"

	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/domain"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/repository"
)

type OrderSuite struct {
	suite.Suite
	repo *repository.OrderRepository
}

const test_storage_path = "../../storage/test_data.json"

func (s *OrderSuite) SetupTest() {
	var err error
	s.repo, err = repository.NewOrderRepository(test_storage_path)
	s.Require().NoError(err)
}

func (s *OrderSuite) TearDownTest() {
	os.Remove(test_storage_path)
}

func (s *OrderSuite) TestAddOrderSuccess() {
	order, _ := domain.NewOrder(dto.AddOrder{
		ID:         10,
		ClientID:   10,
		StoreUntil: time.Now().AddDate(0, 0, 2),
		Cost:       1000,
		Weight:     5,
		Packages:   []string{"unknown", "unknown"},
	})

	err := s.repo.AddOrder(order)
	s.Require().NoError(err)
}

func (s *OrderSuite) TestAddOrderFailed() {
	order, _ := domain.NewOrder(dto.AddOrder{
		ID:         10,
		ClientID:   10,
		StoreUntil: time.Now().AddDate(0, 0, 2),
		Cost:       1000,
		Weight:     5,
		Packages:   []string{"unknown", "unknown"},
	})

	err := s.repo.AddOrder(order)
	s.Require().NoError(err)

	order, _ = domain.NewOrder(dto.AddOrder{
		ID:         10,
		ClientID:   10,
		StoreUntil: time.Now().AddDate(0, 0, 2),
		Cost:       1000,
		Weight:     5,
		Packages:   []string{"unknown", "unknown"},
	})

	err = s.repo.AddOrder(order)
	s.Require().Error(err)
}

func (s *OrderSuite) TestGetOrderByIDSuccess() {
	order, _ := domain.NewOrder(dto.AddOrder{
		ID:         10,
		ClientID:   10,
		StoreUntil: time.Now().AddDate(0, 0, 2),
		Cost:       1000,
		Weight:     5,
		Packages:   []string{"unknown", "unknown"},
	})

	err := s.repo.AddOrder(order)
	s.Require().NoError(err)

	_, err = s.repo.GetOrderByID(10)
	s.Require().NoError(err)
}

func (s *OrderSuite) TestGetOrderByIDFailed() {
	_, err := s.repo.GetOrderByID(10)
	s.Require().Error(err)
}

func (s *OrderSuite) TestGetOrdersByIDSuccess() {
	order, _ := domain.NewOrder(dto.AddOrder{
		ID:         10,
		ClientID:   10,
		StoreUntil: time.Now().AddDate(0, 0, 2),
		Cost:       1000,
		Weight:     5,
		Packages:   []string{"unknown", "unknown"},
	})

	err := s.repo.AddOrder(order)
	s.Require().NoError(err)

	order, _ = domain.NewOrder(dto.AddOrder{
		ID:         11,
		ClientID:   10,
		StoreUntil: time.Now().AddDate(0, 0, 2),
		Cost:       100,
		Weight:     7,
		Packages:   []string{"unknown", "unknown"},
	})

	err = s.repo.AddOrder(order)
	s.Require().NoError(err)

	_, err = s.repo.GetOrdersByID([]int{10, 11})
	s.Require().NoError(err)
}

func (s *OrderSuite) TestGetClientOrdersListSuccess() {
	order, _ := domain.NewOrder(dto.AddOrder{
		ID:         10,
		ClientID:   10,
		StoreUntil: time.Now().AddDate(0, 0, 2),
		Cost:       1000,
		Weight:     5,
		Packages:   []string{"unknown", "unknown"},
	})

	err := s.repo.AddOrder(order)
	s.Require().NoError(err)

	order, _ = domain.NewOrder(dto.AddOrder{
		ID:         11,
		ClientID:   10,
		StoreUntil: time.Now().AddDate(0, 0, 2),
		Cost:       100,
		Weight:     7,
		Packages:   []string{"unknown", "unknown"},
	})

	err = s.repo.AddOrder(order)
	s.Require().NoError(err)

	_, err = s.repo.GetClientOrdersList(10)
	s.Require().NoError(err)
}

func (s *OrderSuite) TestGetRefundsListSuccess() {
	order := &domain.Order{}
	order.SetID(10)
	order.SetClientID(10)
	order.SetStoreUntil(time.Now().AddDate(0, 0, 1))
	order.SetCost(1000)
	order.SetWeight(7)
	order.SetPickUpTime()
	order.SetStatus(domain.OrderStatusRefunded)

	err := s.repo.AddOrder(order)
	s.Require().NoError(err)

	_, err = s.repo.GetRefundsList(0, 0)
	s.Require().NoError(err)
}
