package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func main() {
	storage, err := NewStorage("data.json")
	if err != nil {
		fmt.Println(err)
	}

	storage.readDataFromFile()

	// storage.Orders = append(storage.Orders, Order{
	// 	ID:          11,
	// 	ClientID:    10,
	// 	StoredUntil: time.Now(),
	// 	PickUpTime:  time.Now(),
	// })

	// if err = storage.GetOrderFromСourier(1, 1, time.Now().AddDate(0, 0, 5)); err != nil {
	// 	fmt.Println(err)
	// }

	// if err = storage.GetOrderFromСourier(2, 1, time.Now().AddDate(0, 0, 5)); err != nil {
	// 	fmt.Println(err)
	// }

	storage.writeDataToFile()
}

type IOData struct {
	Orders []Order `json:"orders"`
}

type Order struct {
	ID          int       `json:"id"`
	ClientID    int       `json:"clientID"`
	StoredUntil time.Time `json:"storedUntil"`
	IsPickedUp  bool      `json:"pickedUp"`
	PickUpTime  time.Time `json:"pickUpTime"`
	IsRefund    bool      `json:"refund"`
}

type Storage struct {
	Orders []Order
	Path   string
}

func NewStorage(path string) (*Storage, error) {
	f, err := os.OpenFile(path, os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	f.Close()

	storage := &Storage{Orders: make([]Order, 0), Path: path}
	storage.readDataFromFile()

	return storage, nil
}

func (s *Storage) readDataFromFile() error {
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

func (s *Storage) writeDataToFile() error {
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

func (s *Storage) GetOrderFromСourier(orderID, clientID int, storedUntil time.Time) (err error) {
	for _, order := range s.Orders {
		if orderID == order.ID {
			return fmt.Errorf("Order already exist")
		}
	}

	if storedUntil.Before(time.Now()) {
		return fmt.Errorf("Order store time expired")
	}

	s.Orders = append(s.Orders, Order{
		ID:          orderID,
		ClientID:    clientID,
		StoredUntil: storedUntil,
	})

	if err = s.writeDataToFile(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GiveOrderToCourier(orderID int) (err error) {
	for i, order := range s.Orders {
		if orderID != order.ID {
			continue
		}

		if order.IsPickedUp {
			return fmt.Errorf("Order is picked up by client")
		}

		if order.StoredUntil.After(time.Now()) {
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

func (s *Storage) GiveOrderToClient(orderIDs []int) (err error) {
	uniqueClientID := -1
	for i, orderID := range orderIDs {
		for _, order := range s.Orders {
			if orderID != order.ID {
				continue
			}

			if uniqueClientID != order.ClientID {
				return fmt.Errorf("Orders does not belong to one person")
			}

			if order.StoredUntil.Before(time.Now()) {
				continue
			}

			if uniqueClientID == -1 {
				uniqueClientID = order.ClientID
			}

			s.Orders[i].IsPickedUp = true
			s.Orders[i].PickUpTime = time.Now()
			break
		}
	}

	if err = s.writeDataToFile(); err != nil {
		return err
	}

	return nil
}

// TODO: Add params handler
func (s *Storage) OrderList(clientID int, params ...int) (err error) {
	for _, order := range s.Orders {
		fmt.Println(order)
	}

	return nil
}

func (s *Storage) GetRefundFromСlient(clientID, orderID int) (err error) {
	for i, order := range s.Orders {
		if orderID != order.ID || clientID != order.ClientID {
			continue
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

// TODO: Add pagination
func (s *Storage) RefundList() (err error) {
	for _, order := range s.Orders {
		if order.IsRefund {
			fmt.Println(order)
		}
	}

	return nil
}
