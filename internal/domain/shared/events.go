package shared

// DomainEvent is implemented by all domain events.
type DomainEvent interface {
	EventName() string
}
