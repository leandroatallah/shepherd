package sequences

import (
	"encoding/json"
	"os"

	"github.com/leandroatallah/firefly/internal/engine/contracts/sequences"
)

// Sequence is a list of commands to be executed in order, with additional properties.
type Sequence struct {
	commands            []sequences.Command
	BlockPlayerMovement bool
}

func (s *Sequence) Commands() []sequences.Command {
	return s.commands
}

// CommandData is a wrapper used for parsing commands from JSON.
// It holds the data for all possible command types.
type CommandData struct {
	Type string `json:"command"`

	// Fields for "dialogue"
	Lines       []string `json:"lines,omitempty"`
	Position    string   `json:"position,omitempty"`
	SpeechSpeed int      `json:"speech_speed,omitempty"`

	// Fields for "delay"
	Frames int `json:"frames,omitempty"`

	// Fields for "move_actor"
	TargetID string  `json:"target_id,omitempty"`
	EndX     float64 `json:"end_x,omitempty"`
	Speed    float64 `json:"speed,omitempty"`

	// Fields for "event"
	EventType string                 `json:"event_type,omitempty"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}

// SequenceData is a wrapper used for parsing a full sequence from JSON.
type SequenceData struct {
	Commands            []CommandData `json:"commands"`
	BlockPlayerMovement bool          `json:"block_player_movement,omitempty"`
}

// ToCommand converts the generic CommandData into a specific Command implementation.
func (cd *CommandData) ToCommand() sequences.Command {
	switch cd.Type {
	case "dialogue":
		speed := cd.SpeechSpeed
		if speed == 0 && cd.Speed > 0 {
			speed = int(cd.Speed)
		}
		return &DialogueCommand{Lines: cd.Lines, Position: cd.Position, Speed: speed}
	case "delay":
		return &DelayCommand{Frames: cd.Frames}
	case "move_actor":
		return &MoveActorCommand{
			TargetID: cd.TargetID,
			EndX:     cd.EndX,
			Speed:    cd.Speed,
		}
	case "event":
		return &EventCommand{
			EventType: cd.EventType,
			Payload:   cd.Payload,
		}
	}
	return nil
}

// NewSequenceFromJSON loads a sequence from a JSON file path.
func NewSequenceFromJSON(filePath string) (*Sequence, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return &Sequence{}, err
	}

	var sequenceData SequenceData
	if err := json.Unmarshal(data, &sequenceData); err != nil {
		return &Sequence{}, err
	}

	var commands []sequences.Command
	for _, cd := range sequenceData.Commands {
		cmd := cd.ToCommand()
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}

	return &Sequence{
		commands:            commands,
		BlockPlayerMovement: sequenceData.BlockPlayerMovement,
	}, nil
}
