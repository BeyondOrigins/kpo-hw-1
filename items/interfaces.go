package items

type Kind string
type ID int

const (
	KindBicycle Kind = "bicycle"
	KindScooter Kind = "scooter"
)

type EnergyConsumer interface {
	KWh() int
}

type Identifiable interface {
	ID() ID
}

type Named interface {
	Name() string
}

type Easiness interface {
	Easiness()
}

type Item interface {
	Identifiable
	Easiness
	Named
}
