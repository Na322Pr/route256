package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"gitlab.ozon.dev/marchenkosasha2/homework/internal/domain"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
)

func ReadOrdersFromFile(path string) ([]*domain.Order, error) {
	op := "ReadOrdersFromFile"

	var file *os.File
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer file.Close()

	var listOrdersDTO dto.ListOrdersDTO

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&listOrdersDTO)
	if err != nil {
		if err == io.EOF {
			return []*domain.Order{}, nil
		}
		return nil, fmt.Errorf("error decoding orders: %w", err)
	}

	orders := make([]*domain.Order, 0)

	for _, orderDTO := range listOrdersDTO.Orders {
		var order domain.Order
		order.FromDTO(orderDTO)
		orders = append(orders, &order)
	}

	return orders, nil
}

func WriteOrdersToFile(path string, orders []*domain.Order) error {
	op := "WriteOrdersToFile"

	var file *os.File
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer file.Close()

	var listOrdersDTO dto.ListOrdersDTO
	for _, order := range orders {
		listOrdersDTO.Orders = append(listOrdersDTO.Orders, order.ToDTO())
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(listOrdersDTO)
	if err != nil {
		return fmt.Errorf("error encoding orders: %w", err)
	}

	return nil
}
