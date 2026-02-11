package gamescenephases

// RescueSheepGoal: Complete when all sheep are rescued
type RescueSheepGoal struct {
	scene *PhasesScene
}

func (g *RescueSheepGoal) IsCompleted() bool {
	return g.scene.bodyCounter.sheep > 0 && g.scene.bodyCounter.sheep == g.scene.bodyCounter.sheepRescued
}

func (g *RescueSheepGoal) OnCompletion() {
	g.scene.defaultCompletion()
}

// ReachEndpointGoal: Complete when player reaches endpoint
type ReachEndpointGoal struct {
	scene *PhasesScene
}

func (g *ReachEndpointGoal) IsCompleted() bool {
	return g.scene.reachedEndpoint
}

func (g *ReachEndpointGoal) OnCompletion() {
	g.scene.defaultCompletion()
}
