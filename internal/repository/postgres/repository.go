package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/georgysavva/scany/v2/pgxscan"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/domain"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
)

type PgOrderRepository struct {
	txManager TransactionManager
}

func NewPgOrderRepository(txManager TransactionManager) *PgOrderRepository {
	return &PgOrderRepository{txManager: txManager}
}

func (r *PgOrderRepository) AddOrder(ctx context.Context, orderDTO dto.OrderDTO) error {
	const (
		op = "PgOrderRepository.AddOrder"

		sqlQuery = `insert into orders(order_id, client_id, store_until, status, cost, weight, packages)
		values ($1, $2, $3, $4, $5, $6, $7::varchar[])`
	)

	tx := r.txManager.GetQueryEngine(ctx)

	_, err := tx.Exec(ctx, sqlQuery,
		orderDTO.ID,
		orderDTO.ClientID,
		orderDTO.StoreUntil,
		orderDTO.Status,
		orderDTO.Cost,
		orderDTO.Weight,
		orderDTO.Packages,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *PgOrderRepository) UpdateOrder(ctx context.Context, orderDTO dto.OrderDTO) error {
	const (
		op = "PgOrderRepository.UpdateOrder"

		sqlQuery = `update orders
        set status = $2, pick_up_time = $3
        where order_id = $1`
	)

	tx := r.txManager.GetQueryEngine(ctx)

	_, err := tx.Exec(ctx, sqlQuery, orderDTO.ID, orderDTO.Status, orderDTO.PickUpTime)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *PgOrderRepository) GetOrderByID(ctx context.Context, id int) (*dto.OrderDTO, error) {
	const (
		op = "PgOrderRepository.GetOrderByID"

		sqlQuery = `select * from orders where order_id = $1`
	)

	tx := r.txManager.GetQueryEngine(ctx)

	orders := make([]dto.OrderDTO, 1)

	err := pgxscan.Select(ctx, tx, &orders, sqlQuery, id)

	if len(orders) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrOrderNotFound)
	}

	return &orders[0], err
}

func (r *PgOrderRepository) GetOrdersByIDs(ctx context.Context, ids []int) (*dto.ListOrdersDTO, error) {
	const (
		op = "PgOrderRepository.GetOrdersByIDs"

		sqlQuery = `select * from orders where order_id = any($1)`
	)

	orders := make([]dto.OrderDTO, 1)

	tx := r.txManager.GetQueryEngine(ctx)
	err := pgxscan.Select(ctx, tx, &orders, sqlQuery, ids)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.ListOrdersDTO{Orders: orders}, nil
}

func (r *PgOrderRepository) GetClientOrdersList(ctx context.Context, clientID int) (*dto.ListOrdersDTO, error) {
	const (
		op = "PgOrderRepository.GetClientOrdersList"

		sqlQuery = `select * from orders where client_id = $1 and status = $2`
	)

	orders := make([]dto.OrderDTO, 1)

	tx := r.txManager.GetQueryEngine(ctx)
	err := pgxscan.Select(ctx, tx, &orders, sqlQuery, clientID, domain.OrderStatusMap[domain.OrderStatusReceived])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.ListOrdersDTO{Orders: orders}, err
}

func (r *PgOrderRepository) GetRefundsList(ctx context.Context, limit, offset int) (*dto.ListOrdersDTO, error) {
	const op = "PgOrderRepository.GetRefundsList"

	orders := make([]dto.OrderDTO, 1)

	query := "select * from orders where status = $1 order by order_id "
	params := []any{domain.OrderStatusMap[domain.OrderStatusRefunded]}

	if limit > 0 {
		query += "limit $2 "
		params = append(params, limit)
	}

	if offset > 0 {
		query += "offset $" + strconv.Itoa(len(params)+1)
		params = append(params, offset)
	}

	tx := r.txManager.GetQueryEngine(ctx)
	err := pgxscan.Select(ctx, tx, &orders, query, params...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.ListOrdersDTO{Orders: orders}, err
}
