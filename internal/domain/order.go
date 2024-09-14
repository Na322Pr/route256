package domain

import (
	"fmt"
	"time"

	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
)

type OrderStatus int

const (
	OrderStatusUnknown OrderStatus = iota
	OrderStatusReceived
	OrderStatusPickedUp
	OrderStatusRefunded
	OrderStatusDelete
)

type OrderStatusEntry struct {
	Status OrderStatus
	Name   string
}

var orderStatusEntries = []OrderStatusEntry{
	{OrderStatusUnknown, "unknown"},
	{OrderStatusReceived, "received"},
	{OrderStatusPickedUp, "pickedUp"},
	{OrderStatusRefunded, "refunded"},
	{OrderStatusDelete, "deleted"},
}

var OrderStatusMap = make(map[OrderStatus]string)
var OrderStatusStringMap = make(map[string]OrderStatus)

func init() {
	for _, entry := range orderStatusEntries {
		OrderStatusMap[entry.Status] = entry.Name
		OrderStatusStringMap[entry.Name] = entry.Status
	}
}

// Можно передвавать pickUpTime по ссылке
type Order struct {
	id         int
	clientID   int
	storeUntil time.Time
	status     OrderStatus
	pickUpTime time.Time
}

// var OrderStatusDescriptions = map[OrderStatus]string{}

func NewOrder(id, clientID int, storeUntil time.Time) (*Order, error) {

	if storeUntil.Before(time.Now()) {
		return nil, fmt.Errorf("order store time expired")
	}

	order := &Order{
		id:         id,
		clientID:   clientID,
		storeUntil: storeUntil,
		status:     OrderStatusReceived,
	}

	return order, nil

}

func (o *Order) SetID(id int) error {
	if id < 0 {
		return fmt.Errorf("invalid ID")
	}

	o.id = id
	return nil
}

// Можно проверить, что новый статус идет после старого
func (o *Order) SetStatus(status OrderStatus) error {
	o.status = status
	return nil
}

func (o *Order) UpdateStatus(status OrderStatus) error {
	if status <= o.status {
		return ErrStatusStageDone
	}

	o.status = status

	return nil
}

func (o *Order) SetPickUpTime() {
	o.pickUpTime = time.Now()
}

func (o *Order) GetOrderID() int {
	return o.id
}

func (o *Order) GetOrderClientID() int {
	return o.clientID
}

func (o *Order) GetOrderStatus() string {
	return OrderStatusMap[o.status]
}

func (o *Order) GetOrderStoreUntil() time.Time {
	return o.storeUntil
}

func (o *Order) GetOrderPickUpTime() time.Time {
	return o.pickUpTime
}

func (o *Order) ToDTO() dto.OrderDto {
	return dto.OrderDto{
		ID:         o.id,
		ClientID:   o.clientID,
		StoreUntil: o.storeUntil,
		Status:     OrderStatusMap[o.status],
		PickUpTime: o.pickUpTime,
	}
}

func (o *Order) FromDTO(orderDTO dto.OrderDto) {
	o.id = orderDTO.ID
	o.clientID = orderDTO.ClientID
	o.storeUntil = orderDTO.StoreUntil
	o.pickUpTime = orderDTO.PickUpTime

	orderStatus, ok := OrderStatusStringMap[orderDTO.Status]
	if ok {
		o.status = orderStatus
	}

}
