package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type IOData struct {
	Orders []Order `json:"orders"`
}

type Order struct {
	ID         int       `json:"id"`
	ClientID   int       `json:"clientID"`
	StoreUntil time.Time `json:"storeUntil"`
	IsPickedUp bool      `json:"pickedUp"`
	PickUpTime time.Time `json:"pickUpTime"`
	IsRefund   bool      `json:"refund"`
}

type Store struct {
	Orders []Order
	Path   string
}

func NewStore(path string) (*Store, error) {
	f, err := os.OpenFile(path, os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	f.Close()

	Store := &Store{Orders: make([]Order, 0), Path: path}
	err = Store.readDataFromFile()
	if err != nil {
		return nil, err
	}

	return Store, nil
}

func (s *Store) readDataFromFile() error {
	var file *os.File
	file, err := os.OpenFile(s.Path, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var ioData IOData

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ioData)
	if err != nil {
		return err
	}

	s.Orders = ioData.Orders

	return nil
}

func (s *Store) writeDataToFile() error {
	var file *os.File
	file, err := os.OpenFile(s.Path, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(IOData{s.Orders})
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetOrderFromСourier(orderID, clientID int, storeUntil time.Time) (err error) {
	for _, order := range s.Orders {
		if orderID == order.ID {
			return fmt.Errorf("Order already exist")
		}
	}

	if storeUntil.Before(time.Now()) {
		return fmt.Errorf("Order store time expired")
	}

	s.Orders = append(s.Orders, Order{
		ID:         orderID,
		ClientID:   clientID,
		StoreUntil: storeUntil,
	})

	if err = s.writeDataToFile(); err != nil {
		return err
	}

	return nil
}

func (s *Store) GiveOrderToCourier(orderID int) (err error) {
	for i, order := range s.Orders {
		if orderID != order.ID {
			continue
		}

		if order.IsPickedUp {
			return fmt.Errorf("Order is picked up by client")
		}

		if order.StoreUntil.After(time.Now()) {
			return fmt.Errorf("Order store time is not expired yet")
		}

		s.Orders = append(s.Orders[:i], s.Orders[i+1:]...)
		break
	}

	if err = s.writeDataToFile(); err != nil {
		return err
	}

	return nil
}

func (s *Store) GiveOrderToClient(orderIDs []int) (err error) {
	uniqueClientID := -1

	var orderIndexes []int

	for _, orderID := range orderIDs {
		for i, order := range s.Orders {
			if orderID != order.ID {
				continue
			}

			if order.StoreUntil.Before(time.Now()) || order.IsRefund || order.IsPickedUp {
				continue
			}

			if uniqueClientID == -1 {
				uniqueClientID = order.ClientID
			}

			if uniqueClientID != order.ClientID {
				return fmt.Errorf("Orders does not belong to one person")
			}

			orderIndexes = append(orderIndexes, i)

			break
		}
	}

	for _, ind := range orderIndexes {
		s.Orders[ind].IsPickedUp = true
		s.Orders[ind].PickUpTime = time.Now()
	}

	if err = s.writeDataToFile(); err != nil {
		return err
	}

	return nil
}

func (s *Store) OrderList(clientID int, params ...int) (err error) {
	fmt.Println(s.Orders)

	var orders []Order

	for _, order := range s.Orders {
		if order.ClientID == clientID && !order.IsPickedUp && !order.IsRefund {
			orders = append(orders, order)
		}
	}

	offset := len(orders)

	if len(params) != 0 && params[0] < offset {
		offset = params[0]
	}

	fmt.Println("Список заказов: ")

	for _, order := range orders[len(orders)-offset:] {
		fmt.Println(order.ID)
	}

	return nil
}

func (s *Store) GetRefundFromСlient(clientID, orderID int) (err error) {
	for i, order := range s.Orders {
		if clientID != order.ClientID || orderID != order.ID {
			continue
		}

		if !order.IsPickedUp {
			return fmt.Errorf("Order is not pick up yet")
		}

		if time.Now().After(order.PickUpTime.AddDate(0, 0, 2)) {
			return fmt.Errorf("Time to refund already expired")
		}

		s.Orders[i].IsRefund = true

		if err = s.writeDataToFile(); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("Order was picked up in another point")
}

func (s *Store) RefundList(limit, offset int) (err error) {

	cur_offset := 0
	cur_limit := 0

	for _, order := range s.Orders {
		if !order.IsRefund {
			continue
		}

		if cur_offset != offset {
			cur_offset++
			continue
		}

		fmt.Println(order.ID)
		cur_limit++

		if cur_limit == limit {
			break
		}
	}

	return nil
}
