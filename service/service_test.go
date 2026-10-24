package service

import (
	"errors"
	"testing"

	"hw-1/inspection"
	"hw-1/items"
	"hw-1/repo"
)

type thing struct {
	id       items.ID
	name     string
	easiness int
}

func (t thing) ID() items.ID  { return t.id }
func (t thing) Name() string  { return t.name }
func (t thing) Easiness() int { return t.easiness }

type transport struct {
	thing
}

func (transport) Kind() string { return "test" }

type consumer struct {
	thing
	kwh int
}

func (c consumer) KWh() int { return c.kwh }

type stubCenter struct {
	verdict   inspection.Verdict
	inspected []items.Item
}

func (s *stubCenter) Inspect(it items.Item) inspection.Verdict {
	s.inspected = append(s.inspected, it)
	return s.verdict
}

func newService(v inspection.Verdict) (*Service, *stubCenter) {
	sc := &stubCenter{verdict: v}
	return NewService(&repo.SimpleRepo{}, sc), sc
}

var accept = inspection.Verdict{Accepted: true}

func TestAccept_RejectedByInspection(t *testing.T) {
	s, sc := newService(inspection.Verdict{Reason: "тормоза"})

	v, err := s.Accept(thing{1, "Самокат", 7})
	if err != nil {
		t.Fatal(err)
	}
	if v.Accepted || v.Reason != "тормоза" {
		t.Errorf("вердикт = %+v, ожидался отказ техосмотра", v)
	}
	if len(sc.inspected) != 1 {
		t.Errorf("техосмотр вызван %d раз, ожидался 1", len(sc.inspected))
	}
	if n := len(s.Inventory()); n != 0 {
		t.Errorf("на балансе %d единиц, отклонённое устройство не должно попасть в парк", n)
	}
}

func TestAccept_AcceptedByInspection(t *testing.T) {
	s, _ := newService(accept)

	if v, err := s.Accept(thing{1, "Самокат", 7}); err != nil || !v.Accepted {
		t.Fatalf("Accept = %+v, %v", v, err)
	}
	if n := len(s.Inventory()); n != 1 {
		t.Errorf("на балансе %d единиц, ожидалась 1", n)
	}
}

func TestAccept_InvalidItemNotInspected(t *testing.T) {
	s, sc := newService(accept)
	if _, err := s.Accept(thing{7, "Шлем", 5}); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		item items.Item
		want error
	}{
		{"nil", nil, ErrNilItem},
		{"без номера", thing{0, "Шлем", 5}, ErrNoID},
		{"дубликат номера", thing{7, "Другой шлем", 5}, ErrDuplicate},
	}
	for _, tt := range tests {
		if _, err := s.Accept(tt.item); !errors.Is(err, tt.want) {
			t.Errorf("%s: ошибка %v, ожидалась %v", tt.name, err, tt.want)
		}
	}
	if len(sc.inspected) != 1 {
		t.Errorf("некорректные устройства не должны доходить до техосмотра")
	}
}

func TestTotalDailyKWh(t *testing.T) {
	s, _ := newService(accept)
	for _, it := range []items.Item{
		consumer{thing{1, "Электросамокат", 7}, 3},
		consumer{thing{2, "Зарядный шкаф", 4}, 5},
		transport{thing{3, "Велосипед", 8}},
		thing{4, "Шлем", 9},
	} {
		if _, err := s.Accept(it); err != nil {
			t.Fatal(err)
		}
	}

	if got := s.TotalDailyKWh(); got != 8 {
		t.Errorf("TotalDailyKWh = %d, want 8", got)
	}
}

func TestBeginnerFriendly(t *testing.T) {
	s, _ := newService(accept)
	for _, it := range []items.Item{
		transport{thing{1, "Простой", MinBeginnerEasiness}},
		transport{thing{2, "Сложный", MinBeginnerEasiness - 1}},
		thing{3, "Шлем", MinBeginnerEasiness + 4},
		thing{4, "Док-станция", 1},
	} {
		if _, err := s.Accept(it); err != nil {
			t.Fatal(err)
		}
	}

	got := s.BeginnerFriendly()
	if len(got) != 2 || got[0].ID() != 1 || got[1].ID() != 3 {
		t.Errorf("BeginnerFriendly = %v, ожидались №1 и №3", got)
	}
}

func TestRemove(t *testing.T) {
	s, _ := newService(accept)
	if _, err := s.Accept(thing{1, "Шлем", 5}); err != nil {
		t.Fatal(err)
	}

	if err := s.Remove(1); err != nil {
		t.Fatal(err)
	}
	if err := s.Remove(1); err == nil {
		t.Error("повторное удаление должно вернуть ошибку")
	}
	if n := len(s.Inventory()); n != 0 {
		t.Errorf("на балансе %d единиц после удаления", n)
	}
}
