package erlangb

// EventKind identifies what kind of event occurs at a scheduled time: an
// arrival or a departure.
type EventKind int

const (
	Arrival EventKind = iota
	Departure
)

// Event is something scheduled to happen at a point in simulated time.
type Event struct {
	Time float64
	Kind EventKind
}

// Events is a min-heap of Event ordered by Time, ties broken by Kind
// (Arrival before Departure, matching "arrival" < "departure" in the
// Python version). Use it with heap.Push/heap.Pop directly, the same way
// heapq.heappush/heapq.heappop work on a plain list in Python.
type Events []Event

func (e Events) Len() int { return len(e) }

func (e Events) Less(i, j int) bool {
	if e[i].Time != e[j].Time {
		return e[i].Time < e[j].Time
	}
	return e[i].Kind < e[j].Kind
}

func (e Events) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

func (e *Events) Push(x any) { *e = append(*e, x.(Event)) }

func (e *Events) Pop() any {
	old := *e
	n := len(old)
	item := old[n-1]
	*e = old[:n-1]
	return item
}
