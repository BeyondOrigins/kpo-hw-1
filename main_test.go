package main

import (
	"strings"
	"testing"

	"hw-1/app"
	"hw-1/inspection"
	"hw-1/items"
)

type rejectAll struct{}

func (rejectAll) Inspect(items.Item) inspection.Verdict {
	return inspection.Verdict{Reason: "stub: отказ"}
}

func run(t *testing.T, input string, replace func() inspection.ServiceCenter) string {
	t.Helper()
	var out strings.Builder
	c, err := buildContainer(strings.NewReader(input), &out)
	if err != nil {
		t.Fatal(err)
	}
	if replace != nil {
		if err := c.Decorate(func(inspection.ServiceCenter) inspection.ServiceCenter { return replace() }); err != nil {
			t.Fatal(err)
		}
	}
	if err := c.Invoke(func(a *app.App) error { return a.Run() }); err != nil {
		t.Fatal(err)
	}
	return out.String()
}

func TestContainer_RealServiceCenter(t *testing.T) {
	out := run(t, `add scooter 1 "Ninebot" 7 1`+"\nreport\n", nil)
	for _, want := range []string{"Принято в парк: Ninebot (№1)", "Всего на балансе: 1"} {
		if !strings.Contains(out, want) {
			t.Errorf("в выводе нет %q:\n%s", want, out)
		}
	}
}

func TestContainer_ServiceCenterReplaced(t *testing.T) {
	out := run(t, `add scooter 1 "Ninebot" 7 1`+"\nreport\n",
		func() inspection.ServiceCenter { return rejectAll{} })
	for _, want := range []string{"Техосмотр не пройден: Ninebot (№1): stub: отказ", "Всего на балансе: 0"} {
		if !strings.Contains(out, want) {
			t.Errorf("в выводе нет %q:\n%s", want, out)
		}
	}
}
