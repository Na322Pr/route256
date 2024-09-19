package domain

import (
	"fmt"
	"time"

	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
)

// Status
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

// Package
type OrderPackage int

const (
	OrderPackageUnknown OrderPackage = iota
	OrderPackageBox
	OrderPackageBag
	OrderPackageTape
)

type OrderPackageEntry struct {
	Package OrderPackage
	Name    string
}

var orderPackageEntries = []OrderPackageEntry{
	{OrderPackageUnknown, "unknown"},
	{OrderPackageBox, "box"},
	{OrderPackageBag, "bag"},
	{OrderPackageTape, "tape"},
}

var OrderPackageMap = make(map[OrderPackage]string)
var OrderPackageStringMap = make(map[string]OrderPackage)

// Set up Status and Package convertion maps
func init() {
	for _, entry := range orderStatusEntries {
		OrderStatusMap[entry.Status] = entry.Name
		OrderStatusStringMap[entry.Name] = entry.Status
	}

	for _, entry := range orderPackageEntries {
		OrderPackageMap[entry.Package] = entry.Name
		OrderPackageStringMap[entry.Name] = entry.Package
	}
}

// Package Options Builder
type PackageOption func(*Order) error

func PackBag() PackageOption {
	return func(o *Order) error {
		packageCost := 5
		packageMaxWeight := 10

		if len(o.GetOrderPackages()) != 0 {
			return ErrAlreadyPackaged
		}

		if o.GetOrderWeight() > packageMaxWeight {
			return ErrPackageTooHeavy
		}

		o.AddPackage(OrderPackageBag)
		o.SetCost(o.GetOrderCost() + packageCost)
		return nil
	}
}

func PackBox() PackageOption {
	return func(o *Order) error {
		packageCost := 20
		packageMaxWeight := 30

		if len(o.GetOrderPackages()) != 0 {
			return ErrAlreadyPackaged
		}

		if o.GetOrderWeight() > packageMaxWeight {
			return ErrPackageTooHeavy
		}

		o.AddPackage(OrderPackageBox)
		o.SetCost(o.GetOrderCost() + packageCost)
		return nil
	}
}

func PackTape() PackageOption {
	return func(o *Order) error {
		packageCost := 1
		o.AddPackage(OrderPackageTape)
		o.SetCost(o.GetOrderCost() + packageCost)
		return nil
	}
}

var OrderPackageOptions = map[OrderPackage]PackageOption{
	OrderPackageUnknown: nil,
	OrderPackageBag:     PackBag(),
	OrderPackageBox:     PackBox(),
	OrderPackageTape:    PackTape(),
}

// Order
type Order struct {
	id         int
	clientID   int
	storeUntil time.Time
	status     OrderStatus
	cost       int
	weight     int
	packages   []OrderPackage
	pickUpTime time.Time
}

func NewOrder(orderDTO dto.AddOrder, packOpts ...PackageOption) (*Order, error) {
	op := "Order.NewOrder"

	order := Order{}

	if err := order.SetID(orderDTO.ID); err != nil {
		return nil, err
	}

	if err := order.SetClientID(orderDTO.ClientID); err != nil {
		return nil, err
	}

	if orderDTO.StoreUntil.Before(time.Now()) {
		return nil, ErrStoreTimeExpired
	}
	order.SetStoreUntil(orderDTO.StoreUntil)

	if err := order.SetCost(orderDTO.Cost); err != nil {
		return nil, err
	}

	if err := order.SetWeight(orderDTO.Weight); err != nil {
		return nil, err
	}

	for _, opt := range packOpts {
		if opt == nil {
			continue
		}

		if err := opt(&order); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	order.SetStatus(OrderStatusReceived)
	return &order, nil
}

// Setters
func (o *Order) SetID(id int) error {
	if id < 0 {
		return ErrInvalidID
	}

	o.id = id
	return nil
}

func (o *Order) SetClientID(clientID int) error {
	if clientID < 0 {
		return ErrInvalidClientID
	}

	o.clientID = clientID
	return nil
}

func (o *Order) SetStoreUntil(storeUntil time.Time) error {
	o.storeUntil = storeUntil
	return nil
}

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

func (o *Order) AddPackage(packageType OrderPackage) error {
	o.packages = append(o.packages, packageType)
	return nil
}

func (o *Order) SetPickUpTime() {
	o.pickUpTime = time.Now()
}

func (o *Order) SetCost(cost int) error {
	if cost < 0 {
		return ErrInvalidCost
	}

	o.cost = cost
	return nil
}

func (o *Order) SetWeight(weight int) error {
	if weight < 0 {
		return ErrInvalidWeight
	}

	o.weight = weight
	return nil
}

// Getters
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

func (o *Order) GetOrderCost() int {
	return o.cost
}

func (o *Order) GetOrderWeight() int {
	return o.weight
}

func (o *Order) GetOrderPackages() []string {
	var packagesType []string

	for _, packageType := range o.packages {
		packagesType = append(packagesType, OrderPackageMap[packageType])
	}

	return packagesType
}

func (o *Order) GetOrderPickUpTime() time.Time {
	return o.pickUpTime
}

// DTO Conversion
func (o *Order) ToDTO() dto.OrderDto {
	orderDTO := dto.OrderDto{
		ID:         o.id,
		ClientID:   o.clientID,
		StoreUntil: o.storeUntil,
		Status:     OrderStatusMap[o.status],
		Cost:       o.cost,
		Weight:     o.weight,
		PickUpTime: o.pickUpTime,
	}

	for _, packageType := range o.packages {
		orderDTO.Packages = append(orderDTO.Packages, OrderPackageMap[packageType])
	}

	return orderDTO
}

func (o *Order) FromDTO(orderDTO dto.OrderDto) error {
	if err := o.SetID(orderDTO.ID); err != nil {
		return err
	}

	if err := o.SetClientID(orderDTO.ClientID); err != nil {
		return err
	}

	if err := o.SetCost(orderDTO.Cost); err != nil {
		return err
	}

	if err := o.SetWeight(orderDTO.Weight); err != nil {
		return err
	}

	o.storeUntil = orderDTO.StoreUntil
	o.pickUpTime = orderDTO.PickUpTime

	orderStatus, ok := OrderStatusStringMap[orderDTO.Status]
	if ok {
		o.status = orderStatus
	}

	for _, packageType := range orderDTO.Packages {
		orderPackage, ok := OrderPackageStringMap[packageType]
		if ok {
			o.packages = append(o.packages, orderPackage)
		}
	}

	return nil
}
