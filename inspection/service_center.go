package inspection

import "hw-1/items"

type Verdict struct {
	Accepted bool
	Reason   string
}

type ServiceCenter interface {
	Inspect(it items.Item) Verdict
}

type BasicServiceCenter struct{}

func NewBasicServiceCenter() *BasicServiceCenter {
	return &BasicServiceCenter{}
}

func (*BasicServiceCenter) Inspect(it items.Item) Verdict {
	if e := it.Easiness(); e < 1 || e > 10 {
		return Verdict{Reason: "простота для новичка вне диапазона 1–10"}
	}
	if ec, ok := it.(items.EnergyConsumer); ok && ec.KWh() < 0 {
		return Verdict{Reason: "отрицательное энергопотребление"}
	}
	return Verdict{Accepted: true}
}
