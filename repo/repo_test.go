package repo

import (
	"testing"

	"hw-1/items"
)

func TestSimpleRepo(t *testing.T) {
	r := &SimpleRepo{}
	a, _ := items.NewDevice(items.Spec{ID: 1, Name: "A", Easiness: 5})
	b, _ := items.NewDevice(items.Spec{ID: 2, Name: "B", Easiness: 5})
	r.AddItem(a)
	r.AddItem(b)

	if got := r.GetItem(2); got != b {
		t.Errorf("GetItem(2) = %v", got)
	}
	if got := r.GetItem(3); got != nil {
		t.Errorf("GetItem(3) = %v, want nil", got)
	}

	all := r.All()
	all[0] = nil
	if r.GetItem(1) == nil {
		t.Error("изменение результата All не должно менять хранилище")
	}

	if err := r.RemoveItem(1); err != nil {
		t.Fatal(err)
	}
	if err := r.RemoveItem(1); err == nil {
		t.Error("повторное удаление должно вернуть ошибку")
	}
	if n := len(r.All()); n != 1 {
		t.Errorf("осталось %d, want 1", n)
	}
}
