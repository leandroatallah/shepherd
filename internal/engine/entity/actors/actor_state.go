package actors

type ActorState interface {
	State() ActorStateEnum
	OnStart(currentCount int)
	GetAnimationCount(currentCount int) int
}

type ActorStateEnum int

var (
	Idle    ActorStateEnum
	Walking ActorStateEnum
	Falling ActorStateEnum
	Hurted  ActorStateEnum
)

func init() {
	Idle = RegisterState("idle", func(b BaseState) ActorState { return &IdleState{BaseState: b} })
	Walking = RegisterState("walk", func(b BaseState) ActorState { return &WalkState{BaseState: b} })
	Falling = RegisterState("fall", func(b BaseState) ActorState { return &FallState{BaseState: b} })
	Hurted = RegisterState("hurt", func(b BaseState) ActorState { return &HurtState{BaseState: b} })
}

type BaseState struct {
	actor      ActorEntity
	state      ActorStateEnum
	entryCount int
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
}

func (s *BaseState) GetAnimationCount(currentCount int) int {
	return currentCount - s.entryCount
}
