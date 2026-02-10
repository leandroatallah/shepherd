package gamescenephases

type PhaseGoal interface {
	IsCompleted(s *PhasesScene) bool
	OnCompletion(s *PhasesScene)
}

// RescueSheepGoal: Complete when all sheep are rescued
type RescueSheepGoal struct{}

func (g *RescueSheepGoal) IsCompleted(s *PhasesScene) bool {
	return s.bodyCounter.sheep > 0 && s.bodyCounter.sheep == s.bodyCounter.sheepRescued
}

func (g *RescueSheepGoal) OnCompletion(s *PhasesScene) {
	s.defaultCompletion()
}

// ReachEndpointGoal: Complete when player reaches endpoint
type ReachEndpointGoal struct{}

func (g *ReachEndpointGoal) IsCompleted(s *PhasesScene) bool {
	return s.reachedEndpoint
}

func (g *ReachEndpointGoal) OnCompletion(s *PhasesScene) {
	s.defaultCompletion()
}

// SequenceGoal: Complete when sequence finishes
type SequenceGoal struct{}

func (g *SequenceGoal) IsCompleted(s *PhasesScene) bool {
	return s.sequencePlayer != nil && !s.sequencePlayer.IsPlaying()
}

func (g *SequenceGoal) OnCompletion(s *PhasesScene) {
	s.defaultCompletion()
}

// NoGoal: Never completes
type NoGoal struct{}

func (g *NoGoal) IsCompleted(s *PhasesScene) bool {
	return false
}

func (g *NoGoal) OnCompletion(s *PhasesScene) {}
