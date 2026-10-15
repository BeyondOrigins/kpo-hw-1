package items

type itemBase struct {
	name     string
	id       ID
	easiness int
}

func (obj *itemBase) ID() ID {
	return obj.id
}

func (obj *itemBase) Name() string {
	return obj.name
}

func (obj *itemBase) Easiness() int {
	return obj.easiness
}
