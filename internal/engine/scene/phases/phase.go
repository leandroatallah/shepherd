package phases

type GoalType string

type Phase struct {
	ID             int
	Name           string
	TilemapPath    string
	NextPhaseID    int
	SequencePath   string
	GoalType       GoalType
	ActorBehaviors map[string]ActorBehavior
}

type ActorBehavior struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}
