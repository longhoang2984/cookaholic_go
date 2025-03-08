package interfaces

import "context"

type Event interface {
	Type() string
}

type UserCreatedEvent struct {
	UserID uint
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
