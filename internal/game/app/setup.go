package gamesetup

import (
	"io/fs"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/audio"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/event"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
	"github.com/leandroatallah/firefly/internal/engine/scene"
	"github.com/leandroatallah/firefly/internal/engine/scene/phases"
	"github.com/leandroatallah/firefly/internal/engine/ui/speech"
	gamescene "github.com/leandroatallah/firefly/internal/game/scenes"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
	gamespeech "github.com/leandroatallah/firefly/internal/game/ui/speech"
)

func Setup(assets fs.FS) error {
	cfg := config.Get()
	// Basic Ebiten setup
	ebiten.SetWindowSize(cfg.ScreenWidth*3, cfg.ScreenHeight*3)
	ebiten.SetWindowTitle("Firefly")

	// Initialize all systems and managers
	audioManager := audio.NewAudioManager()
	sceneManager := scene.NewSceneManager()
	phaseManager := phases.NewManager()
	actorManager := actors.NewManager()

	// Initialize Dialogue Manager
	fontText, err := font.NewFontText(cfg.MainFontFace)
	if err != nil {
		return err
	}
	speechFont := speech.NewSpeechFont(fontText, 8, 14)
	speechBubble := gamespeech.NewSpeechBubble(speechFont)
	dialogueManager := speech.NewManager(speechBubble)

	// Load audio assets
	loadAudioAssetsFromFS(assets, audioManager)

	// Load phases
	phase0 := phases.Phase{
		ID:           1,
		Name:         "Phase 1",
		TilemapPath:  "assets/tilemap/shepherd-phase-0.tmj",
		NextPhaseID:  2,
		SequencePath: "assets/sequences/sample.json",
	}
	phase1 := phases.Phase{ID: 2, Name: "Phase 2", TilemapPath: "assets/tilemap/shepherd-phase-1.tmj", NextPhaseID: 1}
	phase2 := phases.Phase{ID: 2, Name: "Phase 2", TilemapPath: "assets/tilemap/shepherd-phase-2.tmj", NextPhaseID: 1}
	phase3 := phases.Phase{ID: 3, Name: "Phase 3", TilemapPath: "assets/tilemap/shepherd-phase-3.tmj", NextPhaseID: 1}
	phase4 := phases.Phase{ID: 4, Name: "Phase 4", TilemapPath: "assets/tilemap/shepherd-phase-4.tmj", NextPhaseID: 1}
	phase5 := phases.Phase{ID: 5, Name: "Phase 5", TilemapPath: "assets/tilemap/shepherd-phase-5.tmj", NextPhaseID: 1}
	phase6 := phases.Phase{ID: 6, Name: "Phase 6", TilemapPath: "assets/tilemap/shepherd-phase-6.tmj", NextPhaseID: 1}
	phaseManager.AddPhase(phase0)
	phaseManager.AddPhase(phase1)
	phaseManager.AddPhase(phase2)
	phaseManager.AddPhase(phase3)
	phaseManager.AddPhase(phase4)
	phaseManager.AddPhase(phase5)
	phaseManager.AddPhase(phase6)
	phaseManager.SetCurrentPhase(6)

	appContext := &app.AppContext{
		AudioManager:    audioManager,
		DialogueManager: dialogueManager,
		EventManager:    event.NewManager(),
		ActorManager:    actorManager,
		SceneManager:    sceneManager,
		PhaseManager:    phaseManager,
		ImageManager:    nil,
		DataManager:     nil,
		Assets:          assets,
		Config:          config.Get(),
		Space:           space.NewSpace(),
	}

	sceneFactory := scene.NewDefaultSceneFactory(gamescene.InitSceneMap(appContext))
	sceneFactory.SetAppContext(appContext)

	sceneManager.SetFactory(sceneFactory)
	sceneManager.SetAppContext(appContext)

	// Create and run the game
	game := app.NewGame(appContext)

	// Set initial game scene
	game.AppContext.SceneManager.NavigateTo(scenestypes.ScenePhases, nil, false)

	if err := ebiten.RunGame(game); err != nil {
		return err
	}

	return nil
}

// loadAudioAssetsFromFS is a helper function to load all audio files from an fs.FS.
func loadAudioAssetsFromFS(assets fs.FS, am *audio.AudioManager) {
	dir := "assets/audio"
	files, err := fs.ReadDir(assets, dir)
	if err != nil {
		log.Fatalf("error reading embedded audio dir: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		// Filter for supported audio types
		if !(strings.HasSuffix(fileName, ".ogg") || strings.HasSuffix(fileName, ".wav") || strings.HasSuffix(fileName, ".mp3")) {
			continue
		}

		fullPath := dir + "/" + fileName
		data, err := fs.ReadFile(assets, fullPath)
		if err != nil {
			log.Printf("failed to read embedded file %s: %v", fullPath, err)
			continue
		}

		// Use the existing Add method to process and store the player.
		am.Add(dir+"/"+fileName, data)
	}
}
