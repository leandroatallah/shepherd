package event

// Event defines the interface for all events.
type Event interface {
	Type() string
}

// Listener is a function that handles an event.
type Listener func(e Event)

// Manager handles event subscription and dispatching.
type Manager struct {
	listeners map[string][]Listener
}

// NewManager creates a new event manager.
func NewManager() *Manager {
	return &Manager{
		listeners: make(map[string][]Listener),
	}
}

// Subscribe adds a listener for a given event type.
func (m *Manager) Subscribe(eventType string, listener Listener) {
	m.listeners[eventType] = append(m.listeners[eventType], listener)
}

// Publish dispatches an event to all registered listeners.
func (m *Manager) Publish(e Event) {
	if listeners, ok := m.listeners[e.Type()]; ok {
		for _, listener := range listeners {
			listener(e)
		}
	}
}
