package items

import (
	"errors"
	"slices"
	"testing"
)

func TestFactories(t *testing.T) {
	spec := Spec{ID: 1, Name: "X", Easiness: 7, KWh: 2}
	tests := []struct {
		name      string
		factory   Factory
		kind      string
		transport bool
		electric  bool
	}{
		{"Device", NewDevice, "", false, false},
		{"ElectricDevice", NewElectricDevice, "", false, true},
		{"Vehicle", VehicleFactory(KindBike), KindBike, true, false},
		{"ElectricVehicle", ElectricVehicleFactory(KindScooter), KindScooter, true, true},
	}
	for _, tt := range tests {
		s := spec
		if !tt.electric {
			s.KWh = 0
		}
		it, err := tt.factory(s)
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		if it.ID() != 1 || it.Name() != "X" || it.Easiness() != 7 {
			t.Errorf("%s: поля не заполнены: %d %q %d", tt.name, it.ID(), it.Name(), it.Easiness())
		}

		tr, isTransport := it.(Transport)
		if isTransport != tt.transport {
			t.Errorf("%s: Transport = %v, want %v", tt.name, isTransport, tt.transport)
		}
		if isTransport && tr.Kind() != tt.kind {
			t.Errorf("%s: Kind = %q, want %q", tt.name, tr.Kind(), tt.kind)
		}

		ec, isElectric := it.(EnergyConsumer)
		if isElectric != tt.electric {
			t.Errorf("%s: EnergyConsumer = %v, want %v", tt.name, isElectric, tt.electric)
		}
		if isElectric && ec.KWh() != 2 {
			t.Errorf("%s: KWh = %d, want 2", tt.name, ec.KWh())
		}
	}
}

func TestFactories_Errors(t *testing.T) {
	for _, f := range []Factory{NewDevice, NewElectricDevice, VehicleFactory(KindBike), ElectricVehicleFactory(KindEBike)} {
		if _, err := f(Spec{ID: 1}); !errors.Is(err, ErrNoName) {
			t.Errorf("без названия: got %v", err)
		}
	}
	for _, f := range []Factory{NewDevice, VehicleFactory(KindBike)} {
		if _, err := f(Spec{ID: 1, Name: "X", KWh: 1}); !errors.Is(err, ErrNoEnergy) {
			t.Errorf("расход у неэлектрического вида: got %v", err)
		}
	}
}

func TestRegistry(t *testing.T) {
	r := NewRegistry()
	if err := r.Register(KindHelmet, NewDevice); err != nil {
		t.Fatal(err)
	}
	if err := r.Register(KindBike, VehicleFactory(KindBike)); err != nil {
		t.Fatal(err)
	}
	if err := r.Register(KindHelmet, NewDevice); err == nil {
		t.Error("повторная регистрация должна вернуть ошибку")
	}
	if err := r.Register("x", nil); err == nil {
		t.Error("nil-фабрика должна вернуть ошибку")
	}

	if got := r.Kinds(); !slices.Equal(got, []string{KindBike, KindHelmet}) {
		t.Errorf("Kinds = %v", got)
	}
	if it, err := r.New(KindBike, Spec{ID: 1, Name: "B", Easiness: 5}); err != nil || it.Name() != "B" {
		t.Errorf("New = %v, %v", it, err)
	}
	if _, err := r.New("rocket", Spec{}); !errors.Is(err, ErrUnknownKind) {
		t.Error("неизвестный вид должен вернуть ошибку")
	}
}
