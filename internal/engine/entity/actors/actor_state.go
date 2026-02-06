package actors

type ActorState interface {
	State() ActorStateEnum
	OnStart(currentCount int)
	GetAnimationCount(currentCount int) int
	IsAnimationFinished() bool
}

type ActorStateEnum int

var (
	Idle    ActorStateEnum
	Walking ActorStateEnum
	Falling ActorStateEnum
	Landing ActorStateEnum
	Hurted  ActorStateEnum
)

func init() {
	Idle = RegisterState("idle", func(b BaseState) ActorState { return &IdleState{BaseState: b} })
	Walking = RegisterState("walk", func(b BaseState) ActorState { return &WalkState{BaseState: b} })
	Falling = RegisterState("fall", func(b BaseState) ActorState { return &FallState{BaseState: b} })
	Landing = RegisterState("land", func(b BaseState) ActorState { return &LandingState{BaseState: b} })
	Hurted = RegisterState("hurt", func(b BaseState) ActorState { return &HurtState{BaseState: b} })
}

type BaseState struct {
	actor      ActorEntity
	state      ActorStateEnum
	entryCount int
	tick       int
}

func NewBaseState(actor ActorEntity, state ActorStateEnum) BaseState {
	return BaseState{actor: actor, state: state}
}

func (s *BaseState) State() ActorStateEnum {
	return s.state
}
func (s *BaseState) GetActor() ActorEntity {
	return s.actor
}

func (s *BaseState) OnStart(currentCount int) {
	s.entryCount = currentCount
	s.tick = 0
}

func (s *BaseState) GetAnimationCount(currentCount int) int {
	return currentCount - s.entryCount
}

func (s *BaseState) IsAnimationFinished() bool {
	s.tick++

	character := s.GetActor().GetCharacter()
	if character == nil {
		return true
	}

	sprite := character.GetSpriteByState(s.State())
	if sprite == nil || sprite.Image == nil {
		return true
	}

	rect := character.Position()
	if rect.Dx() == 0 {
		return true
	}

	elementWidth := sprite.Image.Bounds().Dx()
	frameCount := elementWidth / rect.Dx()

	frameRate := character.FrameRate()
	if frameRate == 0 {
		frameRate = 1
	}

	// Calculate total duration in ticks
	duration := frameCount * frameRate

	return s.tick >= duration
}
