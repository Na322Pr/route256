package repository

import (
	"context"

	"github.com/Na322Pr/route256/internal/dto"
	"github.com/Na322Pr/route256/internal/repository/postgres"
	"github.com/Na322Pr/route256/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (s *StorageFacade) GetOrderByID(ctx context.Context, id int64) (*dto.OrderDTO, error) {
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

func (s *StorageFacade) GetOrdersByIDs(ctx context.Context, ids []int64) (*dto.ListOrdersDTO, error) {
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

func NewFacade(pool *pgxpool.Pool) usecase.OrderRepoFacade {
	txManager := postgres.NewTxManager(pool)
	pgOrderRepository := postgres.NewPgOrderRepository(txManager)
	return NewStorageFacade(txManager, pgOrderRepository)
}
