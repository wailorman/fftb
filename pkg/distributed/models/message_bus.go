package models

// MessageBus _
type MessageBus struct {
}

// NewMessageBus _
func NewMessageBus() *MessageBus {
	return &MessageBus{}
}

// Subscribe _
func (mb *MessageBus) Subscribe() Subscriber {
	return nil
}

// Publish _
func (mb *MessageBus) Publish(IProgress) {
}
