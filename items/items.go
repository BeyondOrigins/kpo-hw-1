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
	Name() string
}

type Easiness interface {
	Easiness()
}

type Transport interface {
	Identifiable
	EnergyConsumer
	Kind() string
}

type base struct {
	name     string
	id       ID
	easiness int
}

type Vehicle struct {
	base
	kind     string
	easiness int
}

type ElectroVehicle struct {
	Vehicle
	kwh int
}

func (obj *base) ID() ID {
	return obj.id
}

func (obj *base) Name() string {
	return obj.name
}

func (obj *base) Easiness() int {
	return obj.easiness
}
