package repo

import (
	"errors"
	"hw-1/items"
	"slices"
)

type SimpleRepo struct {
	items []items.Item
}

func NewSimpleRepo() *SimpleRepo {
	return &SimpleRepo{}
}

func (repo *SimpleRepo) AddItem(item items.Item) {
	repo.items = append(repo.items, item)
}

func (repo *SimpleRepo) GetItem(id items.ID) items.Item {
	for i := 0; i < len(repo.items); i++ {
		if repo.items[i].ID() == id {
			return repo.items[i]
		}
	}
	return nil
}

func (repo *SimpleRepo) All() []items.Item {
	return slices.Clone(repo.items)
}

func (repo *SimpleRepo) RemoveItem(id items.ID) error {
	for i := 0; i < len(repo.items); i++ {
		if repo.items[i].ID() == id {
			repo.items = slices.Delete(repo.items, i, i+1)
			return nil
		}
	}
	return errors.New("предмета с таким номером не существует")
}
