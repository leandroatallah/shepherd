package actors

func init() {
	Idle = RegisterState("idle", func(b BaseState) ActorState { return &IdleState{BaseState: b} })
	Walking = RegisterState("walk", func(b BaseState) ActorState { return &WalkState{BaseState: b} })
	Falling = RegisterState("fall", func(b BaseState) ActorState { return &FallState{BaseState: b} })
	Hurted = RegisterState("hurt", func(b BaseState) ActorState { return &HurtState{BaseState: b} })
}

// Idle
type IdleState struct {
	BaseState
}

func (s *IdleState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
}

// Walking
type WalkState struct {
	BaseState
}

func (s *WalkState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
}

// Falling
type FallState struct {
	BaseState
}

func (s *FallState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
}

// Hurt
type HurtState struct {
	BaseState
	count         int
	durationLimit int
}

func (s *HurtState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
	s.durationLimit = 30 // 0.5 sec, duration of the hurt animation
}

func (s *HurtState) IsAnimationFinished() bool {
	s.count++
	return s.count > s.durationLimit
}
