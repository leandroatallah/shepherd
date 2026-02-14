package phases

import "github.com/leandroatallah/firefly/internal/engine/contracts/navigation"

type GoalType string

type Phase struct {
	ID             int
	Name           string
	TilemapPath    string
	NextPhaseID    int
	SequencePath   string
	GoalType       GoalType
	ActorBehaviors map[string][]ActorBehavior
	SceneType      navigation.SceneType
}

type ActorBehavior struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}
