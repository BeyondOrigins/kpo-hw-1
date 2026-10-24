package service

import (
	"errors"
	"fmt"

	"hw-1/inspection"
	"hw-1/items"
)

const MinBeginnerEasiness = 6

var (
	ErrNilItem   = errors.New("устройство не передано")
	ErrNoID      = errors.New("инвентарный номер должен быть положительным")
	ErrDuplicate = errors.New("инвентарный номер уже занят")
)

type Repo interface {
	AddItem(items.Item)
	GetItem(items.ID) items.Item
	RemoveItem(items.ID) error
	All() []items.Item
}

type Service struct {
	repo Repo
	sc   inspection.ServiceCenter
}

func NewService(repo Repo, sc inspection.ServiceCenter) *Service {
	return &Service{
		repo: repo,
		sc:   sc,
	}
}

func (s *Service) Accept(it items.Item) (inspection.Verdict, error) {
	if it == nil {
		return inspection.Verdict{}, ErrNilItem
	}
	if it.ID() <= 0 {
		return inspection.Verdict{}, ErrNoID
	}
	if s.repo.GetItem(it.ID()) != nil {
		return inspection.Verdict{}, fmt.Errorf("%w: №%d", ErrDuplicate, it.ID())
	}

	verdict := s.sc.Inspect(it)
	if verdict.Accepted {
		s.repo.AddItem(it)
	}
	return verdict, nil
}

func (s *Service) Remove(id items.ID) error {
	return s.repo.RemoveItem(id)
}

func (s *Service) Inventory() []items.Item {
	return s.repo.All()
}

func (s *Service) TotalDailyKWh() int {
	total := 0
	for _, it := range s.repo.All() {
		if c, ok := it.(items.EnergyConsumer); ok {
			total += c.KWh()
		}
	}
	return total
}

func (s *Service) BeginnerFriendly() []items.Item {
	var res []items.Item
	for _, it := range s.repo.All() {
		if it.Easiness() >= MinBeginnerEasiness {
			res = append(res, it)
		}
	}
	return res
}
