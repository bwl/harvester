package components

type Quest struct {
	ID           string
	PlanetID     int
	Description  string
	Objectives   []QuestObjective
	Rewards      []Reward
	EscapeReward bool
}

type QuestObjective struct {
	Type      string
	Target    string
	Count     int
	Location  string
	Completed bool
}

type Reward struct {
	Type  string
	Value interface{}
}
