package app

import (
	"strings"
	"testing"

	"hw-1/inspection"
	"hw-1/items"
	"hw-1/repo"
	"hw-1/service"
	"hw-1/ui"
)

func newTestApp(t *testing.T, input string) (*App, *strings.Builder) {
	t.Helper()
	reg := items.NewRegistry()
	factories := map[string]items.Factory{
		items.KindBike:    items.VehicleFactory(items.KindBike),
		items.KindScooter: items.ElectricVehicleFactory(items.KindScooter),
		items.KindEBike:   items.ElectricVehicleFactory(items.KindEBike),
		items.KindHelmet:  items.NewDevice,
		items.KindDock:    items.NewElectricDevice,
		items.KindCharger: items.NewElectricDevice,
	}
	for kind, f := range factories {
		if err := reg.Register(kind, f); err != nil {
			t.Fatal(err)
		}
	}

	var out strings.Builder
	cli := ui.NewCLI(strings.NewReader(input), &out)
	svc := service.NewService(&repo.SimpleRepo{}, inspection.NewBasicServiceCenter())
	app, err := NewApp(svc, cli, reg)
	if err != nil {
		t.Fatal(err)
	}
	return app, &out
}

func TestDemo(t *testing.T) {
	app, out := newTestApp(t, "demo\nexit\n")
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
	t.Log(out.String())

	for _, want := range []string{
		"Принято в парк: Ninebot Max G30 (№1)",
		"Техосмотр не пройден: Самокат без тормозов (№4)",
		"инвентарный номер уже занят",
		`неизвестный вид "rocket"`,
		"Принято в парк: Зарядный шкаф (№21)",
		"этот вид не потребляет энергию",
		"Суммарный расход: 7 кВт·ч/сут",
		"Всего на балансе: 6",
		"Списано: №3",
		"Всего на балансе: 5",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("в выводе нет %q", want)
		}
	}
}

func TestAdd_BadArgs(t *testing.T) {
	app, out := newTestApp(t, "add scooter\nadd scooter x \"A\"\nadd scooter 1 \"A\" 5 abc\nremove x\nremove 99\nexit\n")
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(out.String(), "Ошибка:"); n != 5 {
		t.Errorf("ожидалось 5 ошибок, получено %d:\n%s", n, out.String())
	}
}

func TestEmptyFleet(t *testing.T) {
	app, out := newTestApp(t, "inventory\nbeginners\nexit\n")
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"На балансе ничего нет.", "Нет устройств, подходящих для новичков."} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("в выводе нет %q", want)
		}
	}
}
