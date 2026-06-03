package domain

import "context"

type OrderEventPublisher interface {
	PublishOrderCreated(ctx context.Context, order Order) error
}
