package application

import (
	"context"
	"fmt"
	"time"

	"github.com/popcoondev/gcp-service-learning/internal/domain"
)

type Clock interface {
	Now() time.Time
}

type IDGenerator interface {
	NewID() string
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}

type StaticIDGenerator struct{}

func (StaticIDGenerator) NewID() string {
	return fmt.Sprintf("ord-%d", time.Now().UTC().UnixNano())
}

type CreateOrderCommand struct {
	CustomerID string
	Items      []OrderItemInput
}

type OrderItemInput struct {
	ProductID string
	Quantity  int
}

type OrderDTO struct {
	OrderID    string
	CustomerID string
	Items      []OrderItemDTO
	Status     string
	CreatedAt  time.Time
}

type OrderItemDTO struct {
	ProductID string
	Quantity  int
}

type OrderService struct {
	repo  domain.OrderRepository
	clock Clock
	idGen IDGenerator
}

func NewOrderService(repo domain.OrderRepository, clock Clock, idGen IDGenerator) OrderService {
	return OrderService{
		repo:  repo,
		clock: clock,
		idGen: idGen,
	}
}

func (s OrderService) CreateOrder(ctx context.Context, cmd CreateOrderCommand) (OrderDTO, error) {
	items := make([]domain.OrderItem, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		items = append(items, domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	order, err := domain.NewOrder(s.idGen.NewID(), cmd.CustomerID, items, s.clock.Now())
	if err != nil {
		return OrderDTO{}, err
	}

	if err := s.repo.Save(ctx, order); err != nil {
		return OrderDTO{}, err
	}

	return toDTO(order), nil
}

func (s OrderService) GetOrder(ctx context.Context, orderID string) (OrderDTO, error) {
	order, err := s.repo.FindByID(ctx, orderID)
	if err != nil {
		return OrderDTO{}, err
	}
	return toDTO(order), nil
}

func toDTO(order domain.Order) OrderDTO {
	items := make([]OrderItemDTO, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, OrderItemDTO{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return OrderDTO{
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		Items:      items,
		Status:     string(order.Status),
		CreatedAt:  order.CreatedAt,
	}
}
