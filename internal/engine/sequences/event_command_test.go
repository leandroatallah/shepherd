package sequences

import (
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/event"
)

func TestEventCommand(t *testing.T) {
	// Setup
	evtManager := event.NewManager()
	appContext := &app.AppContext{
		EventManager: evtManager,
	}

	receivedEvent := false
	var receivedPayload map[string]interface{}

	evtManager.Subscribe("TEST_EVENT", func(e event.Event) {
		receivedEvent = true
		if genEvt, ok := e.(event.GenericEvent); ok {
			receivedPayload = genEvt.Payload
		}
	})

	// Create command
	payload := map[string]interface{}{
		"foo": "bar",
		"baz": 123,
	}
	cmd := &EventCommand{
		EventType: "TEST_EVENT",
		Payload:   payload,
	}

	// Execute
	cmd.Init(appContext)

	// Verify
	if !receivedEvent {
		t.Error("Expected event listener to be called")
	}

	if receivedPayload["foo"] != "bar" {
		t.Errorf("Expected payload 'foo' to be 'bar', got %v", receivedPayload["foo"])
	}

	if receivedPayload["baz"] != 123 {
		t.Errorf("Expected payload 'baz' to be 123, got %v", receivedPayload["baz"])
	}

	if !cmd.Update() {
		t.Error("Expected Update to return true")
	}
}
