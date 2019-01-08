package observerPattern

import "fmt"

type (
	Event struct {
		Date int64
	}

	Observer interface {
		OnNotify(Event)
	}

	Notifier interface {
		Register(Observer)
		Deregister(Observer)
		Notify(Event)
	}
)

type (
	EventObserver struct {
		id int
	}

	EventNotifier struct {
		observers map[Observer]struct{}
	}
)

func NewEventObserver(id int) *EventObserver {
	return &EventObserver{id: id}
}

func NewEventNotifier() EventNotifier {
	return EventNotifier{
		observers: make(map[Observer]struct{}),
	}
}

func (eo *EventObserver) OnNotify(e Event) {
	fmt.Printf("*** Observer %d received: %d\n", eo.id, e.Date)
}

func (en *EventNotifier) Register(o Observer) {
	en.observers[o] = struct{}{}
}

func (en *EventNotifier) Deregister(o Observer) {
	delete(en.observers, o)
}

func (en *EventNotifier) Notify(e Event) {
	for o := range en.observers {
		o.OnNotify(e)
	}
}
