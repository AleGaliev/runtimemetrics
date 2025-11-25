package observer

import (
	"github.com/AleGaliev/runtimemetrics/internal/audit"
)

// интерфейсы
type Publisher interface {
	register(Observer)
	deregister(Observer)
	notify(audit.Audit)
}

type Observer interface {
	SendAudit(audit.Audit) error
	GetID() string
}

// реализация publisher
type Event struct {
	observers map[string]Observer
}

func NewEvent() *Event {
	return &Event{}
}

func (e *Event) Register(o ...Observer) {
	if e.observers == nil {
		e.observers = make(map[string]Observer)
	}
	for _, o := range o {
		e.observers[o.GetID()] = o
	}
}

func (e *Event) Deregister(o Observer) {
	delete(e.observers, o.GetID())
}

func (e *Event) Notify(message audit.Audit) {
	for _, observer := range e.observers {
		observer.SendAudit(message)
	}
}
