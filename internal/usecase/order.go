package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gitlab.ozon.dev/marchenkosasha2/homework/internal/domain"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
)

type Facade interface {
	AddOrder(ctx context.Context, orderDTO dto.OrderDTO) error
	UpdateOrder(ctx context.Context, orderDTO dto.OrderDTO) error
	GetOrderByID(ctx context.Context, id int64) (*dto.OrderDTO, error)
	GetOrdersByIDs(ctx context.Context, ids []int64) (*dto.ListOrdersDTO, error)
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

func (uc *OrderUseCase) ReturnOrderToCourier(ctx context.Context, orderID int64) error {
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

func (uc *OrderUseCase) GiveOrderToClient(ctx context.Context, orderIDs []int64) error {
	op := "OrderUseCase.GiveOrderToClient"

	if len(orderIDs) == 0 {
		return fmt.Errorf("%s: %s", op, "no order IDs")
	}

	listOrdersDTO, err := uc.repo.GetOrdersByIDs(ctx, orderIDs)
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
		orders[i].SetPickUpTime(time.Now())
	}

	uc.giveClientPool(ctx, orders)

	return nil
}

func (uc *OrderUseCase) giveClientPool(ctx context.Context, orders []*domain.Order) {
	const numWorkers = 4
	numOrders := len(orders)

	wg := sync.WaitGroup{}
	resChan := make(chan string, numOrders)
	orderChan := make(chan *domain.Order, numOrders)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go uc.giveClientWorker(ctx, &wg, orderChan, resChan)
	}

	go func() {
		for _, order := range orders {
			orderChan <- order
		}
		close(orderChan)
	}()

	go func() {
		wg.Wait()
		close(resChan)
	}()

	for res := range resChan {
		fmt.Println(res)
	}
}

func (uc *OrderUseCase) giveClientWorker(ctx context.Context, wg *sync.WaitGroup, orders <-chan *domain.Order, result chan<- string) {
	defer wg.Done()

	for order := range orders {
		if err := uc.repo.UpdateOrder(ctx, order.ToDTO()); err != nil {
			result <- fmt.Sprintf("Order %d issue failed", order.GetOrderID())
			continue
		}

		result <- fmt.Sprintf("Order %d issued successfully", order.GetOrderID())
	}
}

func (uc *OrderUseCase) OrderList(ctx context.Context, clientID int) (*dto.ListOrdersDTO, error) {
	op := "OrderUseCase.OrderList"

	listOrdersDTO, err := uc.repo.GetClientOrdersList(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return listOrdersDTO, nil
}

func (uc *OrderUseCase) GetRefundFromСlient(ctx context.Context, clientID int, orderID int64) error {
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

	fmt.Println(limit, offset)

	refundsDTO, err := uc.repo.GetRefundsList(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	fmt.Println(refundsDTO)

	return refundsDTO, nil
}
