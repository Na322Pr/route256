package usecase

import (
	"fmt"
	"time"

	"gitlab.ozon.dev/marchenkosasha2/homework/internal/domain"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/repository"
)

type OrderUseCase struct {
	repo repository.OrderRepository
}

func NewOrderUseCase(repo repository.OrderRepository) *OrderUseCase {
	return &OrderUseCase{repo: repo}
}

func (uc *OrderUseCase) ReceiveOrderFromCourier(req *dto.AddOrderRequest) error {

	order, err := domain.NewOrder(req.ID, req.ClientID, req.StoreUntil)
	if err != nil {
		return fmt.Errorf("ReceiveOrderFromCourier: %w", err)
	}

	err = uc.repo.AddOrder(order)
	if err != nil {
		return fmt.Errorf("ReceiveOrderFromCourier: %w", err)
	}

	uc.repo.Update()
	return nil
}

func (uc *OrderUseCase) ReturnOrderToCourier(orderID int) error {
	order, err := uc.repo.GetOrderByID(orderID)
	if err != nil {
		return fmt.Errorf("GiveOrderToCourier: %w", err)
	}

	orderStatus := order.GetOrderStatus()

	if orderStatus == "pickedUp" {
		return fmt.Errorf("order picked up")
	}

	if orderStatus == "deleted" {
		return fmt.Errorf("order deleted")
	}

	if orderStatus == "received" && order.GetOrderStoreUntil().After(time.Now()) {
		return fmt.Errorf("order store time is not expired yet")
	}

	order.SetStatus(domain.OrderStatusDelete)
	uc.repo.Update()

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

	uc.repo.Update()

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
		return fmt.Errorf("%s: %s", op, "order belong to different client")
	}

	if order.GetOrderStatus() != "pickedUp" {
		return fmt.Errorf("%s: %s", op, "order is non-refundable in its current state")
	}

	if time.Now().After(order.GetOrderPickUpTime().AddDate(0, 0, 2)) {
		return fmt.Errorf("%s: %s", op, "refund time expired")
	}

	order.SetStatus(domain.OrderStatusRefunded)
	uc.repo.Update()

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
