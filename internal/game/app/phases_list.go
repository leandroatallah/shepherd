package gamesetup

import (
	"github.com/leandroatallah/firefly/internal/engine/scene/phases"
	gamescenephases "github.com/leandroatallah/firefly/internal/game/scenes/phases"
)

func GetPhases() []phases.Phase {
	return []phases.Phase{
		{
			ID:           1,
			Name:         "Phase 1",
			TilemapPath:  "assets/tilemap/shepherd-phase-0.tmj",
			NextPhaseID:  2,
			SequencePath: "assets/sequences/sample.json",
			GoalType:     gamescenephases.RescueSheepType,
		},
		{
			ID:          2,
			Name:        "Phase 2",
			TilemapPath: "assets/tilemap/shepherd-phase-1.tmj",
			NextPhaseID: 1,
			GoalType:    gamescenephases.RescueSheepType,
		},
		{
			ID:          3,
			Name:        "Phase 3",
			TilemapPath: "assets/tilemap/shepherd-phase-2.tmj",
			NextPhaseID: 1,
			GoalType:    gamescenephases.NoGoalType,
			// Find a NPC by ID "DOG_NPC" and set its behavior to "follow_player"
			ActorBehaviors: map[string]phases.ActorBehavior{
				"DOG_NPC": {
					Type: "follow_player",
				},
			},
		},
		{
			ID:          4,
			Name:        "Phase 4",
			TilemapPath: "assets/tilemap/shepherd-phase-3.tmj",
			NextPhaseID: 1,
			GoalType:    gamescenephases.RescueSheepType,
		},
		{
			ID:          5,
			Name:        "Phase 5",
			TilemapPath: "assets/tilemap/shepherd-phase-4.tmj",
			NextPhaseID: 1,
			GoalType:    gamescenephases.RescueSheepType,
		},
		{
			ID:          6,
			Name:        "Phase 6",
			TilemapPath: "assets/tilemap/shepherd-phase-5.tmj",
			NextPhaseID: 1,
			GoalType:    gamescenephases.RescueSheepType,
		},
		{
			ID:          7,
			Name:        "Phase 7",
			TilemapPath: "assets/tilemap/shepherd-phase-6.tmj",
			NextPhaseID: 1,
			GoalType:    gamescenephases.RescueSheepType,
		},
	}
}
