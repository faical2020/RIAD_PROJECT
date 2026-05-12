package eventbus

import "sync"

type Event struct {
	Type       string
	EntityID   string
	Data       interface{}
	SequenceID int64
}

type EventBus struct {
	mu          sync.RWMutex
	subscribers []chan Event
	sequence    int64
}

var GlobalBus = &EventBus{}

func (b *EventBus) Subscribe() chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan Event, 100)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *EventBus) Publish(eventType, entityID string, data interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.sequence++
	event := Event{
		Type:       eventType,
		EntityID:   entityID,
		Data:       data,
		SequenceID: b.sequence,
	}

	for _, ch := range b.subscribers {
		select {
		case ch <- event:
		default:
		}
	}
}

func (b *EventBus) GetCurrentSequence() int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.sequence
}
