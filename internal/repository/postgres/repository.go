package postgres

import (
	"context"
	"fmt"

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
	tx := r.txManager.GetQueryEngine(ctx)

	_, err := tx.Exec(ctx, `
		insert into orders(order_id, client_id, store_until, status, cost, weight, packages)
		values ($1, $2, $3, $4, $5, $6, $7::varchar[])
	`,
		orderDTO.ID,
		orderDTO.ClientID,
		orderDTO.StoreUntil,
		orderDTO.Status,
		orderDTO.Cost,
		orderDTO.Weight,
		orderDTO.Packages,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgOrderRepository) UpdateOrder(ctx context.Context, orderDTO dto.OrderDTO) error {
	tx := r.txManager.GetQueryEngine(ctx)

	_, err := tx.Exec(ctx, `
		update orders
        set status = $2,
            pick_up_time = $3
        where order_id = $1
	`,
		orderDTO.ID,
		orderDTO.Status,
		orderDTO.PickUpTime,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgOrderRepository) GetOrderByID(ctx context.Context, id int) (*dto.OrderDTO, error) {
	tx := r.txManager.GetQueryEngine(ctx)

	listOrdersDTO := &dto.ListOrdersDTO{
		Orders: []dto.OrderDTO{},
	}

	err := pgxscan.Select(ctx, tx, &listOrdersDTO.Orders, `
		select * from orders where order_id = $1
	`, id)

	if len(listOrdersDTO.Orders) == 0 {
		return nil, ErrOrderNotFound
	}

	return &listOrdersDTO.Orders[0], err
}

func (r *PgOrderRepository) GetOrdersByID(ctx context.Context, ids []int) (*dto.ListOrdersDTO, error) {
	listOrdersDTO := &dto.ListOrdersDTO{
		Orders: []dto.OrderDTO{},
	}

	tx := r.txManager.GetQueryEngine(ctx)
	err := pgxscan.Select(ctx, tx, &listOrdersDTO.Orders, `
		select * from orders where order_id = any($1)
	`, ids)

	return listOrdersDTO, err
}

func (r *PgOrderRepository) GetClientOrdersList(ctx context.Context, clientID int) (*dto.ListOrdersDTO, error) {
	listOrdersDTO := &dto.ListOrdersDTO{
		Orders: []dto.OrderDTO{},
	}

	tx := r.txManager.GetQueryEngine(ctx)
	err := pgxscan.Select(ctx, tx, &listOrdersDTO.Orders, `
		select * from orders where client_id = $1 and status = $2
	`, clientID, domain.OrderStatusMap[domain.OrderStatusReceived])

	return listOrdersDTO, err
}

func (r *PgOrderRepository) GetRefundsList(ctx context.Context, limit, offset int) (*dto.ListOrdersDTO, error) {
	listOrdersDTO := &dto.ListOrdersDTO{
		Orders: []dto.OrderDTO{},
	}

	query := "select * from orders where status = $1 order by order_id "
	params := []interface{}{domain.OrderStatusMap[domain.OrderStatusRefunded]}

	if limit > 0 {
		query += "limit $2 "
		params = append(params, limit)
	}

	if offset > 0 {
		query += "offset $" + fmt.Sprint(len(params)+1)
		params = append(params, offset)
	}

	tx := r.txManager.GetQueryEngine(ctx)
	err := pgxscan.Select(ctx, tx, &listOrdersDTO.Orders, query, params...)

	return listOrdersDTO, err
}
