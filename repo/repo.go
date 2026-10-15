package repo

import (
	"hw-1/items"
	"slices"
)

type Repo struct {
	items []items.Item
}

func (repo *Repo) AddItem(item items.Item) {
	repo.items = append(repo.items, item)
}

func (repo *Repo) GetItem(id items.ID) items.Item {
	for i := 0; i < len(repo.items); i++ {
		if repo.items[i].ID() == id {
			return repo.items[i]
		}
	}
	return nil
}

func (repo *Repo) RemoveItem(id items.ID) {
	for i := 0; i < len(repo.items); i++ {
		if repo.items[i].ID() == id {
			repo.items = slices.Delete(repo.items, i, i+1)
			return
		}
	}
}
