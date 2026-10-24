package main

import (
	"fmt"
	"io"
	"os"

	"go.uber.org/dig"

	"hw-1/app"
	"hw-1/inspection"
	"hw-1/items"
	"hw-1/repo"
	"hw-1/service"
	"hw-1/ui"
)

func main() {
	c, err := buildContainer(os.Stdin, os.Stdout)
	if err == nil {
		err = c.Invoke(func(a *app.App) error { return a.Run() })
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "ошибка:", err)
		os.Exit(1)
	}
}

func buildContainer(in io.Reader, out io.Writer) (*dig.Container, error) {
	c := dig.New()
	providers := []struct {
		constructor any
		opts        []dig.ProvideOption
	}{
		{func() *ui.CLI { return ui.NewCLI(in, out) }, nil},
		{repo.NewSimpleRepo, []dig.ProvideOption{dig.As(new(service.Repo))}},
		{inspection.NewBasicServiceCenter, []dig.ProvideOption{dig.As(new(inspection.ServiceCenter))}},
		{service.NewService, nil},
		{newRegistry, nil},
		{app.NewApp, nil},
	}
	for _, p := range providers {
		if err := c.Provide(p.constructor, p.opts...); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func newRegistry() (*items.Registry, error) {
	r := items.NewRegistry()
	for kind, f := range map[string]items.Factory{
		items.KindBike:    items.VehicleFactory(items.KindBike),
		items.KindScooter: items.ElectricVehicleFactory(items.KindScooter),
		items.KindEBike:   items.ElectricVehicleFactory(items.KindEBike),
		items.KindHelmet:  items.NewDevice,
		items.KindDock:    items.NewElectricDevice,
		items.KindCharger: items.NewElectricDevice,
	} {
		if err := r.Register(kind, f); err != nil {
			return nil, err
		}
	}
	return r, nil
}
