package components

type Resource struct {
	Kind   string
	Amount int
}

type Inventory struct{ Items map[string]int }

func (i *Inventory) Ensure() {
	if i.Items == nil {
		i.Items = make(map[string]int)
	}
}

type Faction struct{ Name string }

type Damage struct{ Amount int }
