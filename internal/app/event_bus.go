package app

import (
	"context"
	"cookaholic/internal/interfaces"
	"sync"
)

// eventBus implements the EventBus interface, providing a thread-safe event publishing and subscription system.
type eventBus struct {
	handlers map[string][]interfaces.EventHandler
	mu       sync.RWMutex
}

// NewEventBus creates and initializes a new event bus instance.
func NewEventBus() *eventBus {
	return &eventBus{
		handlers: make(map[string][]interfaces.EventHandler),
	}
}

// Publish sends an event to all registered handlers for that event type.
// Uses read lock for thread-safe access to handlers map.
// Returns error if any handler fails.
func (b *eventBus) Publish(ctx context.Context, event interfaces.Event) error {
	b.mu.RLock()
	handlers := b.handlers[event.Type()]
	b.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// Subscribe registers a handler for a specific event type.
// Uses write lock for thread-safe handler registration.
func (b *eventBus) Subscribe(eventType string, handler interfaces.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}
