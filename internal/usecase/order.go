package usecase

import (
	"fmt"
	"time"

	"gitlab.ozon.dev/marchenkosasha2/homework/internal/domain"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
)

type orderRepository interface {
	AddOrder(newOrder *domain.Order) error
	GetOrderByID(id int) (*domain.Order, error)
	GetOrdersByID(ids []int) ([]*domain.Order, error)
	GetClientOrdersList(clientID int) ([]*domain.Order, error)
	GetRefundsList(limit, offset int) ([]*domain.Order, error)
	Update() error
}

type OrderUseCase struct {
	repo orderRepository
}

func NewOrderUseCase(repo orderRepository) *OrderUseCase {
	return &OrderUseCase{repo: repo}
}

func (uc *OrderUseCase) ReceiveOrderFromCourier(req dto.AddOrder) error {
	op := "OrderUseCase.ReceiveOrderFromCourier"

	order, err := domain.NewOrder(
		req,
		domain.OrderPackageOptions[domain.OrderPackageStringMap[req.Packages[0]]],
		domain.OrderPackageOptions[domain.OrderPackageStringMap[req.Packages[1]]],
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = uc.repo.AddOrder(order)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := uc.repo.Update(); err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	return nil
}

func (uc *OrderUseCase) ReturnOrderToCourier(orderID int) error {
	op := "OrderUseCase.ReturnOrderToCourier"

	order, err := uc.repo.GetOrderByID(orderID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	orderStatus := order.GetOrderStatus()

	if orderStatus == "pickedUp" {
		return ErrOrderPickedUp
	}

	if orderStatus == "deleted" {
		return ErrOrderDeleted
	}

	if orderStatus == "received" && order.GetOrderStoreUntil().After(time.Now()) {
		return ErrOrderStoreTimeNotExpired
	}

	order.SetStatus(domain.OrderStatusDelete)
	if err := uc.repo.Update(); err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	return nil
}

func (uc *OrderUseCase) GiveOrderToClient(orderIDs []int) error {
	op := "OrderUseCase.GiveOrderToClient"

	if len(orderIDs) == 0 {
		return fmt.Errorf("%s: %s", op, "no order IDs")
	}

	orders, err := uc.repo.GetOrdersByID(orderIDs)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	clientID := orders[0].GetOrderClientID()

	for _, order := range orders {
		if order.GetOrderClientID() != clientID {
			return fmt.Errorf("%s: %s", op, "orders belong to several clients")
		}
	}

	for i := 0; i < len(orders); i++ {
		orders[i].SetStatus(domain.OrderStatusPickedUp)
		orders[i].SetPickUpTime()
	}

	if err := uc.repo.Update(); err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	return nil
}

func (uc *OrderUseCase) OrderList(clientID int) (*dto.ListOrdersDTO, error) {
	op := "OrderUseCase.OrderList"

	orders, err := uc.repo.GetClientOrdersList(clientID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ordersDTO := dto.ListOrdersDTO{}
	for _, order := range orders {
		ordersDTO.Orders = append(ordersDTO.Orders, order.ToDTO())
	}

	return &ordersDTO, nil
}

func (uc *OrderUseCase) GetRefundFromСlient(clientID, orderID int) error {
	op := "OrderUseCase.GetRefundFromСlient"

	order, err := uc.repo.GetOrderByID(orderID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if order.GetOrderClientID() != clientID {
		return fmt.Errorf("%s: %w", op, ErrOrderClientMismatch)
	}

	if order.GetOrderStatus() != "pickedUp" {
		return fmt.Errorf("%s: %w", op, ErrOrderIsNotRefundable)
	}

	if time.Now().After(order.GetOrderPickUpTime().AddDate(0, 0, 2)) {
		return fmt.Errorf("%s: %s", op, "refund time expired")
	}

	order.SetStatus(domain.OrderStatusRefunded)

	if err := uc.repo.Update(); err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	return nil
}

func (uc *OrderUseCase) RefundList(limit, offset int) (*dto.ListOrdersDTO, error) {
	op := "OrderUseCase.RefundList"

	refunds, err := uc.repo.GetRefundsList(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	refundsDTO := dto.ListOrdersDTO{}
	for _, order := range refunds {
		refundsDTO.Orders = append(refundsDTO.Orders, order.ToDTO())
	}

	return &refundsDTO, nil
}
