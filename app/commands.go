package app

import (
	"errors"
	"fmt"
	"strings"
	"text/tabwriter"

	"hw-1/items"
	"hw-1/ui"
)

const addUsage = `add <вид> <номер> "<название>" <простота> [кВт·ч/сут]`

func (app *App) add(cmd ui.Command) error {
	if n := len(cmd.Args); n < 4 || n > 5 {
		return fmt.Errorf("%w: %s", ui.ErrArgsCount, addUsage)
	}
	id, err := cmd.IntArg(1)
	if err != nil {
		return err
	}
	easiness, err := cmd.IntArg(3)
	if err != nil {
		return err
	}
	spec := items.Spec{ID: items.ID(id), Name: cmd.Args[2], Easiness: easiness}
	if len(cmd.Args) > 4 {
		if spec.KWh, err = cmd.IntArg(4); err != nil {
			return err
		}
	}

	it, err := app.registry.New(cmd.Args[0], spec)
	if errors.Is(err, items.ErrUnknownKind) {
		return fmt.Errorf("%w; доступные виды: %s", err, strings.Join(app.registry.Kinds(), ", "))
	}
	if err != nil {
		return err
	}
	verdict, err := app.srv.Accept(it)
	if err != nil {
		return err
	}
	if !verdict.Accepted {
		app.cli.Printf("Техосмотр не пройден: %s (№%d): %s\n", it.Name(), it.ID(), verdict.Reason)
		return nil
	}
	app.cli.Printf("Принято в парк: %s (№%d)\n", it.Name(), it.ID())
	return nil
}

func (app *App) remove(cmd ui.Command) error {
	if err := cmd.RequireArgs(1); err != nil {
		return err
	}
	id, err := cmd.IntArg(0)
	if err != nil {
		return err
	}
	if err := app.srv.Remove(items.ID(id)); err != nil {
		return err
	}
	app.cli.Printf("Списано: №%d\n", id)
	return nil
}

func (app *App) kinds(cmd ui.Command) error {
	if err := cmd.RequireArgs(0); err != nil {
		return err
	}
	app.cli.Printf("Виды: %s\n", strings.Join(app.registry.Kinds(), ", "))
	return nil
}

func (app *App) inventory(cmd ui.Command) error {
	if err := cmd.RequireArgs(0); err != nil {
		return err
	}
	list := app.srv.Inventory()
	if len(list) == 0 {
		app.cli.Println("На балансе ничего нет.")
		return nil
	}
	app.cli.PrintHeader("Инвентарь")
	app.printItems(list)
	return nil
}

func (app *App) report(cmd ui.Command) error {
	if err := cmd.RequireArgs(0); err != nil {
		return err
	}
	list := app.srv.Inventory()
	transport := 0
	for _, it := range list {
		if _, ok := it.(items.Transport); ok {
			transport++
		}
	}
	app.cli.PrintHeader("Отчёт по парку")
	app.cli.Printf("Всего на балансе: %d\n", len(list))
	app.cli.Printf("  транспорт:      %d\n", transport)
	app.cli.Printf("  вещи:           %d\n", len(list)-transport)
	app.cli.Printf("Для новичков:     %d\n", len(app.srv.BeginnerFriendly()))
	app.cli.Printf("Суммарный расход: %d кВт·ч/сут\n", app.srv.TotalDailyKWh())
	return nil
}

func (app *App) energy(cmd ui.Command) error {
	if err := cmd.RequireArgs(0); err != nil {
		return err
	}
	app.cli.Printf("Суммарный расход: %d кВт·ч/сут\n", app.srv.TotalDailyKWh())
	return nil
}

func (app *App) beginners(cmd ui.Command) error {
	if err := cmd.RequireArgs(0); err != nil {
		return err
	}
	list := app.srv.BeginnerFriendly()
	if len(list) == 0 {
		app.cli.Println("Нет устройств, подходящих для новичков.")
		return nil
	}
	app.cli.PrintHeader("Устройства для новичков")
	app.printItems(list)
	return nil
}

var demoScript = []string{
	"kinds",
	`add scooter 1 "Ninebot Max G30" 7 1`,
	`add ebike 2 "Volteco Bolt" 4 2`,
	`add bike 3 "Stels Navigator" 9`,
	`add helmet 10 "Шлем M" 10`,
	`add dock 20 "Док-станция у Покровки" 3 1`,
	`add charger 21 "Зарядный шкаф" 7 3`,
	`add scooter 4 "Самокат без тормозов" 12 1`,
	`add bike 3 "Дубликат номера" 5`,
	`add helmet 11 "Шлем с фонарём" 9 1`,
	`add rocket 5 "Ракета" 5`,
	"inventory",
	"beginners",
	"energy",
	"report",
	"remove 3",
	"report",
}

func (app *App) demo(cmd ui.Command) error {
	if err := cmd.RequireArgs(0); err != nil {
		return err
	}
	for _, line := range demoScript {
		app.cli.Printf("\n> %s\n", line)
		if _, err := app.cli.Execute(line); err != nil {
			app.cli.Printf("Ошибка: %v\n", err)
		}
	}
	return nil
}

func (app *App) printItems(list []items.Item) {
	var b strings.Builder
	w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "№\tНазвание\tВид\tПростота\tкВт·ч/сут")
	for _, it := range list {
		kind, kwh := "вещь", "—"
		if t, ok := it.(items.Transport); ok {
			kind = t.Kind()
		}
		if c, ok := it.(items.EnergyConsumer); ok {
			kwh = fmt.Sprint(c.KWh())
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%s\n", it.ID(), it.Name(), kind, it.Easiness(), kwh)
	}
	w.Flush()
	app.cli.Printf("%s", b.String())
}
