package service

import (
	"hw-1/items"
	"hw-1/repo"
	servicecenter "hw-1/service_center"
	"hw-1/ui"
)

type Service struct {
	ui                ui.UI
	repo              repo.Repo
	sc                servicecenter.ServiceCenter
	energyConsumption int
}

func NewService(ui ui.UI, repo repo.Repo, sc servicecenter.ServiceCenter) *Service {
	return &Service{
		ui:   ui,
		repo: repo,
		sc:   sc,
	}
}

func (s *Service) GetEnergyConsumption() int {
	return s.energyConsumption
}

func (s *Service) AddItem(item items.Identifiable) {
	// обновить подсчёт потребления энергии в день
	if с, ok := item.(items.EnergyConsumer); ok {
		s.energyConsumption += 24 * с.KWh()
	}
	s.repo.AddItem(item)
}

func (s *Service) RemoveItem(id items.ID) {
	item := s.repo.GetItem(id)
	if item == nil {
		return
	}
	if c, ok := item.(items.EnergyConsumer); ok {
		s.energyConsumption -= 24 * c.KWh()
	}
}

func (s *Service) PrintReport() {
	s.ui.Println("")
}
