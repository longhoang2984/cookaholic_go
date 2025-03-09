package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type Event interface {
	Type() string
}

type UserCreatedEvent struct {
	UserID uuid.UUID
	Email  string
}

func (e UserCreatedEvent) Type() string {
	return "user.created"
}

type EventHandler interface {
	Handle(ctx context.Context, event Event) error
}

type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventType string, handler EventHandler)
}
