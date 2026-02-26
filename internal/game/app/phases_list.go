package gamesetup

import (
	"github.com/leandroatallah/firefly/internal/engine/scene/phases"
	gamescenephases "github.com/leandroatallah/firefly/internal/game/scenes/phases"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
)

func GetPhases() []phases.Phase {
	return []phases.Phase{
		{
			ID:           1,
			Name:         "Story Intro - Part 1",
			NextPhaseID:  2,
			SequencePath: "assets/sequences/story-1.json",
			GoalType:     gamescenephases.SequenceGoalType,
			SceneType:    scenestypes.SceneStory,
		},
		{
			ID:                  2,
			Name:                "Story Intro - Part 2",
			TilemapPath:         "assets/tilemap/shepherd-phase-story-intro-part-2.tmj",
			NextPhaseID:         3,
			GoalType:            gamescenephases.SequenceGoalType,
			SceneType:           scenestypes.ScenePhases,
			SequencePath:        "assets/sequences/story-2.json",
			BlockPlayerMovement: true,
		},
		{
			ID:           3,
			Name:         "Story Intro - Part 3",
			NextPhaseID:  4,
			SequencePath: "assets/sequences/story-3.json",
			GoalType:     gamescenephases.SequenceGoalType,
			SceneType:    scenestypes.SceneStory,
		},
		{
			ID:                  4,
			Name:                "Story Intro - Part 4",
			TilemapPath:         "assets/tilemap/shepherd-phase-4.tmj",
			NextPhaseID:         5,
			GoalType:            gamescenephases.SequenceGoalType,
			SceneType:           scenestypes.ScenePhases,
			SequencePath:        "assets/sequences/story-4.json",
			BlockPlayerMovement: true,
		},
		{
			ID:           5,
			Name:         "Story Intro - Part 5",
			NextPhaseID:  6,
			SequencePath: "assets/sequences/story-5.json",
			GoalType:     gamescenephases.SequenceGoalType,
			SceneType:    scenestypes.SceneStory,
		},
		{
			ID:          6,
			Name:        "Story Intro - Part 6",
			Title:       "Phase title",
			NextPhaseID: 7,
			GoalType:    gamescenephases.NoGoalType,
			SceneType:   scenestypes.ScenePhaseTitle,
		},
		{
			ID:           7,
			Name:         "Area 1 - Stage 1",
			TilemapPath:  "assets/tilemap/shepherd-phase-1-1.tmj",
			NextPhaseID:  8,
			SequencePath: "assets/sequences/phase-1-1.json",
			GoalType:     gamescenephases.ReactEndpointType,
			SceneType:    scenestypes.ScenePhases,
		},
		{
			ID:          8,
			Name:        "Area 1 - Stage 2",
			TilemapPath: "assets/tilemap/shepherd-phase-1-2.tmj",
			NextPhaseID: 9,
			GoalType:    gamescenephases.ReactEndpointType,
			SceneType:   scenestypes.ScenePhases,
		},
		{
			ID:          9,
			Name:        "Area 1 - Stage 2",
			TilemapPath: "assets/tilemap/shepherd-phase-1-3.tmj",
			NextPhaseID: 1,
			GoalType:    gamescenephases.ReactEndpointType,
			SceneType:   scenestypes.ScenePhases,
		},
	}
}
