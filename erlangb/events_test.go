package erlangb

import (
	"container/heap"
	"testing"
)

func TestEvents_PopsInTimeOrder(t *testing.T) {
	events := &Events{}
	heap.Push(events, Event{Time: 3.5, Kind: Arrival})
	heap.Push(events, Event{Time: 1.0, Kind: Departure})
	heap.Push(events, Event{Time: 2.2, Kind: Arrival})

	want := []float64{1.0, 2.2, 3.5}
	for i, w := range want {
		if events.Len() == 0 {
			t.Fatalf("pop %d: expected an event, queue was empty", i)
		}
		got := heap.Pop(events).(Event)
		if got.Time != w {
			t.Errorf("pop %d: got time %v, want %v", i, got.Time, w)
		}
	}
}

func TestEvents_TiesBreakArrivalBeforeDeparture(t *testing.T) {
	events := &Events{}
	heap.Push(events, Event{Time: 5.0, Kind: Departure})
	heap.Push(events, Event{Time: 5.0, Kind: Arrival})

	first := heap.Pop(events).(Event)
	if first.Kind != Arrival {
		t.Fatalf("got %+v, want an Arrival event first on a time tie", first)
	}

	second := heap.Pop(events).(Event)
	if second.Kind != Departure {
		t.Fatalf("got %+v, want a Departure event second", second)
	}
}

func TestEvents_Len(t *testing.T) {
	events := &Events{}
	if events.Len() != 0 {
		t.Fatalf("new queue: got Len() = %d, want 0", events.Len())
	}

	heap.Push(events, Event{Time: 1.0, Kind: Arrival})
	heap.Push(events, Event{Time: 2.0, Kind: Departure})
	if events.Len() != 2 {
		t.Fatalf("after 2 pushes: got Len() = %d, want 2", events.Len())
	}

	heap.Pop(events)
	if events.Len() != 1 {
		t.Fatalf("after 1 pop: got Len() = %d, want 1", events.Len())
	}
}
