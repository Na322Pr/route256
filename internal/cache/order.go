package cache

import (
	"time"

	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
)

type OrderCache struct {
	cli *CacheClient[int64, *dto.OrderDTO]
}

func NewOrderCache(ttl time.Duration) *OrderCache {
	return &OrderCache{
		cli: NewCacheClient[int64, *dto.OrderDTO](ttl),
	}
}

func (c *OrderCache) Get(orderID int64) (*dto.OrderDTO, bool) {
	return c.cli.Get(orderID)
}

func (c *OrderCache) Set(orderDTO *dto.OrderDTO, now time.Time) error {
	c.cli.Set(orderDTO.ID, orderDTO, now)
	return nil
}
