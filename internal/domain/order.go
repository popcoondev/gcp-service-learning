package domain

import (
	"errors"
	"strings"
	"time"
)

type OrderStatus string

const OrderStatusAccepted OrderStatus = "ACCEPTED"

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrInvalidOrder     = errors.New("invalid order")
	ErrInvalidOrderID   = errors.New("invalid order id")
	ErrInvalidCustomer  = errors.New("invalid customer id")
	ErrInvalidOrderItem = errors.New("invalid order item")
)

type Order struct {
	ID         string
	CustomerID string
	Items      []OrderItem
	Status     OrderStatus
	CreatedAt  time.Time
}

type OrderItem struct {
	ProductID string
	Quantity  int
}

func NewOrder(id string, customerID string, items []OrderItem, now time.Time) (Order, error) {
	if strings.TrimSpace(id) == "" {
		return Order{}, ErrInvalidOrderID
	}
	if strings.TrimSpace(customerID) == "" {
		return Order{}, ErrInvalidCustomer
	}
	if len(items) == 0 {
		return Order{}, ErrInvalidOrder
	}

	normalizedItems := make([]OrderItem, 0, len(items))
	for _, item := range items {
		if strings.TrimSpace(item.ProductID) == "" || item.Quantity <= 0 {
			return Order{}, ErrInvalidOrderItem
		}
		normalizedItems = append(normalizedItems, OrderItem{
			ProductID: strings.TrimSpace(item.ProductID),
			Quantity:  item.Quantity,
		})
	}

	return Order{
		ID:         strings.TrimSpace(id),
		CustomerID: strings.TrimSpace(customerID),
		Items:      normalizedItems,
		Status:     OrderStatusAccepted,
		CreatedAt:  now.UTC(),
	}, nil
}
