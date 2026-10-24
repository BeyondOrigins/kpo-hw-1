package items

type Vehicle struct {
	itemBase
	kind string
}

func (v *Vehicle) Kind() string {
	return v.kind
}

type ElectricVehicle struct {
	Vehicle
	energy
}
