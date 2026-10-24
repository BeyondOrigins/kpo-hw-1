package items

type ID int

type EnergyConsumer interface {
	KWh() int
}

type Item interface {
	ID() ID
	Name() string
	Easiness() int
}

type Transport interface {
	Item
	Kind() string
}
