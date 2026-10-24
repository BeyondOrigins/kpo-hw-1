package items

import (
	"errors"
	"fmt"
)

const (
	KindBike    = "bike"
	KindScooter = "scooter"
	KindEBike   = "ebike"
	KindHelmet  = "helmet"
	KindDock    = "dock"
	KindCharger = "charger"
)

var (
	ErrNoName   = errors.New("нет названия")
	ErrNoEnergy = errors.New("этот вид не потребляет энергию")
)

func NewDevice(s Spec) (Item, error) {
	base, err := newBase(s, false)
	if err != nil {
		return nil, err
	}
	return &Device{base}, nil
}

func NewElectricDevice(s Spec) (Item, error) {
	base, err := newBase(s, true)
	if err != nil {
		return nil, err
	}
	return &ElectricDevice{base, energy{s.KWh}}, nil
}

func VehicleFactory(kind string) Factory {
	return func(s Spec) (Item, error) {
		base, err := newBase(s, false)
		if err != nil {
			return nil, err
		}
		return &Vehicle{base, kind}, nil
	}
}

func ElectricVehicleFactory(kind string) Factory {
	return func(s Spec) (Item, error) {
		base, err := newBase(s, true)
		if err != nil {
			return nil, err
		}
		return &ElectricVehicle{Vehicle{base, kind}, energy{s.KWh}}, nil
	}
}

func newBase(s Spec, electric bool) (itemBase, error) {
	if s.Name == "" {
		return itemBase{}, ErrNoName
	}
	if !electric && s.KWh != 0 {
		return itemBase{}, fmt.Errorf("%w: %q", ErrNoEnergy, s.Name)
	}
	return itemBase{name: s.Name, id: s.ID, easiness: s.Easiness}, nil
}
