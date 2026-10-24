package ui

import (
	"errors"
	"strings"
	"testing"
)

func TestCLI_Run(t *testing.T) {
	in := strings.NewReader("ping a \"b c\"\n\nfail\nunknown\nhelp\nexit\nping never\n")
	var out strings.Builder
	cli := NewCLI(in, &out)

	var got []string
	mustRegister(t, cli, CommandSpec{Name: "ping", Usage: "ping <args...>", Help: "эхо",
		Run: func(cmd Command) error {
			got = append(got, strings.Join(cmd.Args, "|"))
			return nil
		}})
	mustRegister(t, cli, CommandSpec{Name: "fail", Help: "всегда ошибка",
		Run: func(Command) error { return errors.New("boom") }})

	if err := cli.Run(); err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 || got[0] != "a|b c" {
		t.Errorf("ping вызван с %q; команды после exit выполняться не должны", got)
	}
	for _, want := range []string{"Ошибка: fail: boom", `неизвестная команда "unknown"`, "ping <args...>", "эхо"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("в выводе нет %q:\n%s", want, out.String())
		}
	}
}

func TestCLI_RunStopsAtEOF(t *testing.T) {
	var out strings.Builder
	if err := NewCLI(strings.NewReader("help"), &out).Run(); err != nil {
		t.Fatal(err)
	}
}

func TestCLI_Register(t *testing.T) {
	cli := NewCLI(strings.NewReader(""), &strings.Builder{})
	noop := func(Command) error { return nil }

	mustRegister(t, cli, CommandSpec{Name: "a", Run: noop})
	for _, spec := range []CommandSpec{
		{Name: "a", Run: noop},
		{Name: "", Run: noop},
		{Name: "b"},
		{Name: "help", Run: noop},
	} {
		if err := cli.Register(spec); err == nil {
			t.Errorf("Register(%q) должен вернуть ошибку", spec.Name)
		}
	}
}

func mustRegister(t *testing.T, cli *CLI, spec CommandSpec) {
	t.Helper()
	if err := cli.Register(spec); err != nil {
		t.Fatal(err)
	}
}
