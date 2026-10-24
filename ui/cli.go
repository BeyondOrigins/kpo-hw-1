package ui

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"text/tabwriter"
)

type Handler func(cmd Command) error

type CommandSpec struct {
	Name  string
	Usage string
	Help  string
	Run   Handler
}

const (
	cmdHelp = "help"
	cmdExit = "exit"
)

type CLI struct {
	in       *bufio.Scanner
	out      io.Writer
	commands map[string]CommandSpec
	order    []string
}

func NewCLI(in io.Reader, out io.Writer) *CLI {
	return &CLI{
		in:       bufio.NewScanner(in),
		out:      out,
		commands: make(map[string]CommandSpec),
	}
}

func (c *CLI) Register(spec CommandSpec) error {
	switch {
	case spec.Name == "":
		return errors.New("register: пустое имя команды")
	case spec.Run == nil:
		return fmt.Errorf("register %q: нет обработчика", spec.Name)
	case spec.Name == cmdHelp || spec.Name == cmdExit:
		return fmt.Errorf("register %q: имя зарезервировано", spec.Name)
	}
	if _, exists := c.commands[spec.Name]; exists {
		return fmt.Errorf("register %q: команда уже есть", spec.Name)
	}
	c.commands[spec.Name] = spec
	c.order = append(c.order, spec.Name)
	return nil
}

func (c *CLI) Run() error {
	c.Println("ВышКат — учёт парка. help — список команд, exit — выход.")
	for {
		fmt.Fprint(c.out, "> ")
		if !c.in.Scan() {
			c.Println("")
			return c.in.Err()
		}
		quit, err := c.Execute(c.in.Text())
		if err != nil {
			c.Printf("Ошибка: %v\n", err)
		}
		if quit {
			return nil
		}
	}
}

func (c *CLI) Execute(line string) (quit bool, err error) {
	cmd, err := ParseLine(line)
	if errors.Is(err, ErrEmptyLine) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	switch cmd.Name {
	case cmdExit:
		return true, nil
	case cmdHelp:
		c.printHelp()
		return false, nil
	}

	spec, ok := c.commands[cmd.Name]
	if !ok {
		return false, fmt.Errorf("неизвестная команда %q, введите help", cmd.Name)
	}
	if err := spec.Run(cmd); err != nil {
		return false, fmt.Errorf("%s: %w", cmd.Name, err)
	}
	return false, nil
}

func (c *CLI) printHelp() {
	c.PrintHeader("Команды")
	w := tabwriter.NewWriter(c.out, 0, 0, 2, ' ', 0)
	for _, name := range c.order {
		spec := c.commands[name]
		usage := spec.Usage
		if usage == "" {
			usage = spec.Name
		}
		fmt.Fprintf(w, "  %s\t%s\n", usage, spec.Help)
	}
	fmt.Fprintf(w, "  %s\t%s\n", cmdHelp, "список команд")
	fmt.Fprintf(w, "  %s\t%s\n", cmdExit, "выход")
	w.Flush()
}

func (c *CLI) Printf(format string, args ...any) {
	fmt.Fprintf(c.out, format, args...)
}

func (c *CLI) Println(s string) {
	fmt.Fprintln(c.out, s)
}

func (c *CLI) PrintHeader(title string) {
	fmt.Fprintf(c.out, "\n== %s ==\n", title)
}
