package navigation

type Phase interface {
	// TODO: Add at least one generic method
}

type PhaseManager interface {
	AddPhase(p Phase)
	AdvanceToNextPhase() error
	GetCurrentPhase() (Phase, error)
	GetPhase(id int) (Phase, error)
	SetCurrentPhase(id int) error
}
