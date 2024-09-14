package repository

import (
	"fmt"

	"gitlab.ozon.dev/marchenkosasha2/homework/internal/domain"
	"gitlab.ozon.dev/marchenkosasha2/homework/storage"
)

type OrderRepository struct {
	orders []*domain.Order
	path   string
}

func NewOrderRepository(path string) (*OrderRepository, error) {
	op := "NewOrderRepository"

	orders, err := storage.ReadOrdersFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &OrderRepository{
		orders: orders,
		path:   path,
	}, nil
}

func (r *OrderRepository) AddOrder(newOrder *domain.Order) error {
	for _, order := range r.orders {
		if newOrder.GetOrderID() == order.GetOrderID() {
			return fmt.Errorf("order already exist")
		}
	}

	r.orders = append(r.orders, newOrder)
	storage.WriteOrdersToFile(r.path, r.orders)

	return nil
}

func (r *OrderRepository) GetOrderByID(id int) (*domain.Order, error) {
	for _, order := range r.orders {
		if order.GetOrderID() == id {
			return order, nil
		}
	}

	return nil, fmt.Errorf("order %d not found", id)
}

func (r *OrderRepository) GetOrdersByID(ids []int) ([]*domain.Order, error) {
	var orders []*domain.Order

	for i, id := range ids {
		for _, order := range r.orders {
			if order.GetOrderID() == id {
				orders = append(orders, order)
				break
			}
		}

		if len(orders) <= i {
			return nil, fmt.Errorf("order %d not found", id)
		}
	}

	return orders, nil
}

func (r *OrderRepository) GetClientOrdersList(clientID int) ([]*domain.Order, error) {
	var orders []*domain.Order

	for _, order := range r.orders {
		if order.GetOrderClientID() != clientID {
			continue
		}

		if order.GetOrderStatus() != "received" {
			continue
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) GetRefundsList(limit, offset int) ([]*domain.Order, error) {
	curLimit, curOffset := 0, 0

	var orders []*domain.Order

	for _, order := range r.orders {
		if order.GetOrderStatus() != "refunded" {
			continue
		}

		if curOffset != offset {
			curOffset++
			continue
		}

		orders = append(orders, order)
		curLimit++

		if curLimit == limit {
			break
		}
	}

	return orders, nil
}

func (r *OrderRepository) Update() error {
	op := "OrderRepository.Update"
	err := storage.WriteOrdersToFile(r.path, r.orders)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
