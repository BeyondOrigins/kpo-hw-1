package ui

import "fmt"

type UI interface {
	Run() error
	Parse()
	Println(string)
	PrintHeader(string)
}

type CLI struct {
}

func (c *CLI) Run() error {
	return nil
}

func (c *CLI) Println(info string) {
	fmt.Println(info)
}

func (c *CLI) PrintHeader(info string) {
	fmt.Println(info)
	fmt.Println("----------------------")
}
