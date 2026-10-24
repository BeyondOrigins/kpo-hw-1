package app

import (
	"hw-1/items"
	"hw-1/service"
	"hw-1/ui"
)

type App struct {
	srv      *service.Service
	cli      *ui.CLI
	registry *items.Registry
}

func NewApp(srv *service.Service, cli *ui.CLI,
	registry *items.Registry) (*App, error) {
	app := &App{cli: cli, srv: srv, registry: registry}
	for _, spec := range app.commands() {
		if err := cli.Register(spec); err != nil {
			return nil, err
		}
	}
	return app, nil
}

func (app *App) Run() error {
	return app.cli.Run()
}

func (app *App) commands() []ui.CommandSpec {
	return []ui.CommandSpec{
		{Name: "add", Usage: addUsage, Help: "принять устройство в парк через техосмотр", Run: app.add},
		{Name: "remove", Usage: "remove <номер>", Help: "списать с баланса", Run: app.remove},
		{Name: "kinds", Usage: "kinds", Help: "виды, которые можно добавить", Run: app.kinds},
		{Name: "inventory", Usage: "inventory", Help: "всё на балансе с инвентарными номерами", Run: app.inventory},
		{Name: "report", Usage: "report", Help: "число единиц и суммарный расход", Run: app.report},
		{Name: "energy", Usage: "energy", Help: "суммарный суточный расход", Run: app.energy},
		{Name: "beginners", Usage: "beginners", Help: "устройства, подходящие для новичков", Run: app.beginners},
		{Name: "demo", Usage: "demo", Help: "прогнать все сценарии (добавляет устройства в парк)", Run: app.demo},
	}
}
