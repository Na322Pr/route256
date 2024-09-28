package usecase

import (
	"context"
	"fmt"
	"time"

	"gitlab.ozon.dev/marchenkosasha2/homework/internal/domain"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
)

type Facade interface {
	AddOrder(ctx context.Context, orderDTO *dto.OrderDTO) error
	UpdateOrder(ctx context.Context, orderDTO *dto.OrderDTO) error
	GetOrderByID(ctx context.Context, id int) (*dto.OrderDTO, error)
	GetOrdersByID(ctx context.Context, ids []int) (*dto.ListOrdersDTO, error)
	GetClientOrdersList(ctx context.Context, clientID int) (*dto.ListOrdersDTO, error)
	GetRefundsList(ctx context.Context, limit, offset int) (*dto.ListOrdersDTO, error) // Update() error
}

type OrderUseCase struct {
	repo Facade
}

func NewOrderUseCase(repo Facade) *OrderUseCase {
	return &OrderUseCase{repo: repo}
}

func (uc *OrderUseCase) ReceiveOrderFromCourier(ctx context.Context, req dto.AddOrder) error {
	op := "OrderUseCase.ReceiveOrderFromCourier"

	order, err := domain.NewOrder(
		req,
		domain.OrderPackageOptions[domain.OrderPackageStringMap[req.Packages[0]]],
		domain.OrderPackageOptions[domain.OrderPackageStringMap[req.Packages[1]]],
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = uc.repo.AddOrder(ctx, order.ToDTO())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *OrderUseCase) ReturnOrderToCourier(ctx context.Context, orderID int) error {
	op := "OrderUseCase.ReturnOrderToCourier"

	orderDTO, err := uc.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var order domain.Order
	order.FromDTO(*orderDTO)

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

	if err := uc.repo.UpdateOrder(ctx, order.ToDTO()); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *OrderUseCase) GiveOrderToClient(ctx context.Context, orderIDs []int) error {
	op := "OrderUseCase.GiveOrderToClient"

	if len(orderIDs) == 0 {
		return fmt.Errorf("%s: %s", op, "no order IDs")
	}

	listOrdersDTO, err := uc.repo.GetOrdersByID(ctx, orderIDs)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if len(listOrdersDTO.Orders) == 0 {
		return fmt.Errorf("%s: %s", op, "no orders")
	}

	var orders []*domain.Order
	for i := 0; i < len(listOrdersDTO.Orders); i++ {
		var order domain.Order
		order.FromDTO(listOrdersDTO.Orders[i])
		orders = append(orders, &order)
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

	// Multiple updates
	for _, order := range orders {
		if err := uc.repo.UpdateOrder(ctx, order.ToDTO()); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (uc *OrderUseCase) OrderList(ctx context.Context, clientID int) (*dto.ListOrdersDTO, error) {
	op := "OrderUseCase.OrderList"

	listOrdersDTO, err := uc.repo.GetClientOrdersList(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return listOrdersDTO, nil
}

func (uc *OrderUseCase) GetRefundFromСlient(ctx context.Context, clientID, orderID int) error {
	op := "OrderUseCase.GetRefundFromСlient"

	orderDTO, err := uc.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var order domain.Order
	order.FromDTO(*orderDTO)

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
	if err := uc.repo.UpdateOrder(ctx, order.ToDTO()); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *OrderUseCase) RefundList(ctx context.Context, limit, offset int) (*dto.ListOrdersDTO, error) {
	op := "OrderUseCase.RefundList"

	refundsDTO, err := uc.repo.GetRefundsList(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return refundsDTO, nil
}
