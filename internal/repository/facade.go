package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/repository/postgres"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/usecase"
)

type StorageFacade struct {
	txManager         postgres.TransactionManager
	pgOrderRepository postgres.PgOrderRepository
}

func NewStorageFacade(
	txManager postgres.TransactionManager,
	pgOrderRepository *postgres.PgOrderRepository,
) *StorageFacade {
	return &StorageFacade{
		txManager:         txManager,
		pgOrderRepository: *pgOrderRepository,
	}
}

func (s *StorageFacade) AddOrder(ctx context.Context, orderDTO dto.OrderDTO) error {
	return s.txManager.RunSerializable(ctx, func(ctxTx context.Context) error {
		err := s.pgOrderRepository.AddOrder(ctx, orderDTO)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *StorageFacade) UpdateOrder(ctx context.Context, orderDTO dto.OrderDTO) error {
	return s.txManager.RunSerializable(ctx, func(ctxTx context.Context) error {
		err := s.pgOrderRepository.UpdateOrder(ctx, orderDTO)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *StorageFacade) GetOrderByID(ctx context.Context, id int) (*dto.OrderDTO, error) {
	var orderDTO *dto.OrderDTO

	err := s.txManager.RunReadCommitted(ctx, func(ctxTx context.Context) error {
		c, err := s.pgOrderRepository.GetOrderByID(ctx, id)
		if err != nil {
			return err
		}

		orderDTO = c
		return nil
	})

	return orderDTO, err
}

func (s *StorageFacade) GetOrdersByIDs(ctx context.Context, ids []int) (*dto.ListOrdersDTO, error) {
	var listOrdersDTO *dto.ListOrdersDTO

	err := s.txManager.RunReadCommitted(ctx, func(ctxTx context.Context) error {
		c, err := s.pgOrderRepository.GetOrdersByIDs(ctx, ids)
		if err != nil {
			return err
		}

		listOrdersDTO = c
		return nil
	})

	return listOrdersDTO, err
}

func (s *StorageFacade) GetClientOrdersList(ctx context.Context, clientID int) (*dto.ListOrdersDTO, error) {
	return s.pgOrderRepository.GetClientOrdersList(ctx, clientID)
}

func (s *StorageFacade) GetRefundsList(ctx context.Context, limit, offset int) (*dto.ListOrdersDTO, error) {
	return s.pgOrderRepository.GetRefundsList(ctx, limit, offset)
}

func NewFacade(pool *pgxpool.Pool) usecase.Facade {
	txManager := postgres.NewTxManager(pool)
	pgOrderRepository := postgres.NewPgOrderRepository(txManager)
	return NewStorageFacade(txManager, pgOrderRepository)
}
