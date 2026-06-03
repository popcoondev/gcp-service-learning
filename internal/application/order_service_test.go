package application

import (
	"context"
	"testing"
	"time"

	"github.com/popcoondev/gcp-service-learning/internal/domain"
)

type stubClock struct {
	now time.Time
}

func (s stubClock) Now() time.Time {
	return s.now
}

type stubIDGenerator struct {
	id string
}

func (s stubIDGenerator) NewID() string {
	return s.id
}

type inMemoryRepo struct {
	orders map[string]domain.Order
}

func (r *inMemoryRepo) Save(_ context.Context, order domain.Order) error {
	r.orders[order.ID] = order
	return nil
}

func (r *inMemoryRepo) FindByID(_ context.Context, orderID string) (domain.Order, error) {
	order, ok := r.orders[orderID]
	if !ok {
		return domain.Order{}, domain.ErrOrderNotFound
	}
	return order, nil
}

func TestCreateOrder(t *testing.T) {
	repo := &inMemoryRepo{orders: map[string]domain.Order{}}
	service := NewOrderService(
		repo,
		stubClock{now: time.Date(2026, 6, 3, 0, 0, 0, 0, time.UTC)},
		stubIDGenerator{id: "ord-test"},
	)

	order, err := service.CreateOrder(context.Background(), CreateOrderCommand{
		CustomerID: "cust-1",
		Items: []OrderItemInput{
			{ProductID: "prod-1", Quantity: 2},
		},
	})
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}
	if order.OrderID != "ord-test" {
		t.Fatalf("unexpected order id: %s", order.OrderID)
	}
	if order.Status != "ACCEPTED" {
		t.Fatalf("unexpected status: %s", order.Status)
	}
}

func TestGetOrderNotFound(t *testing.T) {
	repo := &inMemoryRepo{orders: map[string]domain.Order{}}
	service := NewOrderService(repo, stubClock{now: time.Now()}, stubIDGenerator{id: "ord-x"})

	_, err := service.GetOrder(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}
