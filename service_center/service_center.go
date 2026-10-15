package servicecenter

import (
	"bufio"
	"fmt"
	"hw-1/items"
	"os"
)

type Verdict struct {
	Accepted bool
	Reason   string
}

type ServiceCenter interface {
	Inspect(item *items.Identifiable) Verdict
}

type CLIServiceCenter struct {
}

func NewVerdict(acc bool, reason string) *Verdict {
	return &Verdict{Accepted: acc, Reason: reason}
}

func (sc *CLIServiceCenter) Inspect(item items.Identifiable) Verdict {
	fmt.Printf("Принять %s", item.Name())
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {

	}

	return Verdict{}
}
