package sequences

import (
	"fmt"
	"math"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/event"
	"github.com/leandroatallah/firefly/internal/engine/ui/speech"
)

// EventCommand publishes an event to the global event manager.
type EventCommand struct {
	EventType string
	Payload   map[string]interface{}

	eventManager *event.Manager
}

func (c *EventCommand) Init(appContext any) {
	c.eventManager = appContext.(*app.AppContext).EventManager
	if c.eventManager != nil {
		evt := event.GenericEvent{
			EventType: c.EventType,
			Payload:   c.Payload,
		}
		c.eventManager.Publish(evt)
	}
}

func (c *EventCommand) Update() bool {
	return true
}

// DialogueCommand displays one or more lines of text and waits for player input.
type DialogueCommand struct {
	Lines           []string
	Position        string
	Speed           int
	dialogueManager *speech.Manager
}

func (c *DialogueCommand) Init(appContext any) {
	c.dialogueManager = appContext.(*app.AppContext).DialogueManager
	c.dialogueManager.ShowMessages(c.Lines, c.Position, c.Speed)
}

func (c *DialogueCommand) Update() bool {
	// The command is done when the dialogue manager is no longer speaking.
	return !c.dialogueManager.IsSpeaking()
}

// DelayCommand waits for a specified number of frames.
type DelayCommand struct {
	Frames int
	timer  int
}

func (c *DelayCommand) Init(appContext any) {
	c.timer = 0
}

func (c *DelayCommand) Update() bool {
	c.timer++
	return c.timer >= c.Frames
}

// MoveActorCommand moves a target actor to a specified X position.
type MoveActorCommand struct {
	TargetID string
	EndX     float64
	Speed    float64

	targetActor actors.ActorEntity
	isDone      bool

	lastX         float64
	stuckFrames   int
	initWaitCount int
}

func (c *MoveActorCommand) Init(appContext any) {
	actor, found := appContext.(*app.AppContext).ActorManager.Find(c.TargetID)
	if !found {
		fmt.Printf("MoveActorCommand: Actor with ID '%s' not found.\n", c.TargetID)
		c.isDone = true
		return
	}
	c.targetActor = actor

	if model := actor.MovementModel(); model != nil {
		model.SetIsScripted(true)
	}

	c.lastX = float64(actor.Position().Min.X)
	c.stuckFrames = 0
	c.initWaitCount = 0
}

func (c *MoveActorCommand) Update() bool {
	if c.isDone || c.targetActor == nil {
		return true
	}

	currentX := float64(c.targetActor.Position().Min.X)

	// Stuck detection: if we haven't moved significantly for a while, finish the command.
	// We wait a few frames (initWaitCount) to allow physics to kick in.
	if c.initWaitCount < 10 {
		c.initWaitCount++
	} else {
		if math.Abs(currentX-c.lastX) < 0.1 {
			c.stuckFrames++
		} else {
			c.stuckFrames = 0
		}
	}
	c.lastX = currentX

	const stuckThreshold = 60 // 1 second at 60fps
	if c.stuckFrames >= stuckThreshold {
		c.isDone = true
	}

	speed := c.Speed
	if speed == 0 {
		speed = float64(c.targetActor.Speed())
	}

	distance := c.EndX - currentX

	const arrivalThreshold = 20.0
	const brakingDistance = 10.0 // This value may need tuning depending on friction and speed

	if c.isDone || math.Abs(distance) < arrivalThreshold {
		c.isDone = true
		// Restore player control before finishing the command.
		if model := c.targetActor.MovementModel(); model != nil {
			model.SetIsScripted(false)
		}
		return true
	}

	// When we are close, stop applying force and let friction do the work.
	if math.Abs(distance) < brakingDistance {
		// No-op, just wait for the actor to glide to a stop.
	} else {
		// Apply force to move towards the target.
		if distance > 0 {
			c.targetActor.OnMoveRight(int(speed))
		} else {
			c.targetActor.OnMoveLeft(int(speed))
		}
	}

	return false
}
