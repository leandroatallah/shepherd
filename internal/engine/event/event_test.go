package event

import "testing"

type testEvt struct{ k string }

func (e testEvt) Type() string { return e.k }

func TestEventSubscribeAndPublish(t *testing.T) {
	m := NewManager()

	count := 0
	m.Subscribe("X", func(e Event) {
		count++
		if e.Type() != "X" {
			t.Fatalf("wrong type")
		}
	})

	m.Publish(testEvt{"X"})
	m.Publish(testEvt{"Y"})
	m.Publish(testEvt{"X"})

	if count != 2 {
		t.Fatalf("expected 2 invocations, got %d", count)
	}
}

func TestGenericEvent(t *testing.T) {
	e := GenericEvent{EventType: "A", Payload: map[string]interface{}{"v": 1}}
	if e.Type() != "A" || e.Payload["v"].(int) != 1 {
		t.Fatalf("unexpected generic event")
	}
}
