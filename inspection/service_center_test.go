package inspection

import (
	"testing"

	"hw-1/items"
)

type thing struct{ easiness int }

func (thing) ID() items.ID    { return 1 }
func (thing) Name() string    { return "Шлем" }
func (t thing) Easiness() int { return t.easiness }

type consumer struct {
	thing
	kwh int
}

func (c consumer) KWh() int { return c.kwh }

func TestBasicServiceCenter(t *testing.T) {
	tests := []struct {
		name string
		item items.Item
		want bool
	}{
		{"нормальная вещь", thing{6}, true},
		{"простота 1", thing{1}, true},
		{"простота 10", thing{10}, true},
		{"простота 0", thing{0}, false},
		{"простота 11", thing{11}, false},
		{"нормальный потребитель", consumer{thing{5}, 2}, true},
		{"нулевой расход", consumer{thing{5}, 0}, true},
		{"отрицательный расход", consumer{thing{5}, -1}, false},
	}
	sc := NewBasicServiceCenter()
	for _, tt := range tests {
		v := sc.Inspect(tt.item)
		if v.Accepted != tt.want {
			t.Errorf("%s: Accepted = %v, want %v", tt.name, v.Accepted, tt.want)
		}
		if !v.Accepted && v.Reason == "" {
			t.Errorf("%s: при отказе нужна причина", tt.name)
		}
	}
}
