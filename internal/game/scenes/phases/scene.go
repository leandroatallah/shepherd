package gamescenephases

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	sequencestypes "github.com/leandroatallah/firefly/internal/engine/contracts/sequences"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/enemies"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/npcs"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/platformer"
	"github.com/leandroatallah/firefly/internal/engine/entity/items"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/leandroatallah/firefly/internal/engine/render/particles/vfx"
	"github.com/leandroatallah/firefly/internal/engine/render/screenutil"
	"github.com/leandroatallah/firefly/internal/engine/scene"
	"github.com/leandroatallah/firefly/internal/engine/scene/pause"
	"github.com/leandroatallah/firefly/internal/engine/scene/phases"
	"github.com/leandroatallah/firefly/internal/engine/scene/transition"
	"github.com/leandroatallah/firefly/internal/engine/sequences"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
	gameenemies "github.com/leandroatallah/firefly/internal/game/entity/actors/enemies"
	gamenpcs "github.com/leandroatallah/firefly/internal/game/entity/actors/npcs"
	gameitems "github.com/leandroatallah/firefly/internal/game/entity/items"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
)

const (
	bgSound = "assets/audio/Goblins_Den_Regular.ogg"
)

type PhasesScene struct {
	*scene.TilemapScene

	count       int
	player      platformer.PlatformerActorEntity
	mainText    *font.FontText
	bodyCounter *BodyCounter
	allowPause  bool

	// Complete phase
	isConcludingPhase   bool
	phaseCompleted      bool
	phaseCompletedDelay int
	reachedEndpoint     bool
	hasEndpoints        bool
	goal                phases.Goal

	// Reboot
	isRebooting bool
	rebootDelay int

	// UI effects
	ShowDrawScreenFlash int

	screenFlipper *scene.ScreenFlipper
	// sequencePlayer *sequences.SequencePlayer
	sequencePlayer sequencestypes.Player
	pauseScreen    *pause.PauseScreen
	vfxManager     *vfx.Manager
}

func NewPhasesScene(ctx *app.AppContext) *PhasesScene {
	mainText, err := font.NewFontText(config.Get().MainFontFace)
	if err != nil {
		log.Fatal(err)
	}
	tilemapScene := scene.NewTilemapScene(ctx)
	scene := &PhasesScene{
		TilemapScene: tilemapScene,
		mainText:     mainText,
		bodyCounter:  &BodyCounter{},
	}
	scene.SetAppContext(ctx)

	subscribeEvents(ctx, scene)

	return scene
}

func (s *PhasesScene) OnStart() {
	s.TilemapScene.OnStart()
	s.count = 0
	s.vfxManager = vfx.NewManager("assets/particles/vfx.json")

	// Create player and register to space and context
	p, err := createPlayer(s.AppContext(), gameentitytypes.ShepherdPlayerType)
	if err != nil {
		log.Fatal(err)
	}
	s.player = p
	s.AppContext().ActorManager.Register(s.player)
	s.PhysicsSpace().AddBody(s.player)

	// Optionally block input for the current player of this phase
	if phase, err := s.AppContext().PhaseManager.GetCurrentPhase(); err == nil && phase.BlockPlayerMovement {
		if p, ok := s.AppContext().ActorManager.GetPlayer(); ok {
			p.BlockMovement()
		}
	}

	s.initTilemap()

	// After init bodies, set body counter
	s.bodyCounter.setBodyCounter(s.PhysicsSpace())

	s.PhysicsSpace().Bodies()

	s.SetCameraConfig(scene.CameraConfig{Mode: scene.CameraModeFollow})
	s.Camera().SetFollowTarget(s.player)

	// Init collisions bodies and touch trigger for endpoints
	s.Tilemap().CreateCollisionBodies(s.PhysicsSpace(), func(id string) body.Touchable {
		return bodyphysics.NewTouchTrigger(func() {
			s.endpointTrigger(id)
		}, s.player)
	})

	s.screenFlipper = scene.NewScreenFlipper(s.Camera(), s.player, s.Tilemap(), s.AppContext())
	tileWidth := s.Tilemap().Tilewidth
	s.screenFlipper.PlayerPushDistance = float64(tileWidth / 2)
	s.screenFlipper.FlipStrategy = func(dx, dy int) scene.FlipType {
		if dy != 0 {
			return scene.FlipTypeInstant
		}
		return scene.FlipTypeSmooth
	}
	s.screenFlipper.OnFlipStart = func() {
		s.player.SetImmobile(true)
	}
	s.screenFlipper.OnFlipFinish = func() {
		s.player.SetImmobile(false)
	}
	s.screenFlipper.SnapToCurrentRoom()

	s.pauseScreen = pause.NewPauseScreen(ebiten.KeyEnter, 250*time.Millisecond)

	phase, err := s.AppContext().PhaseManager.GetCurrentPhase()
	if err == nil && phase.SequencePath != "" {
		s.sequencePlayer = sequences.NewSequencePlayer(s.AppContext())
		s.allowPause = phase.GoalType != SequenceGoalType
		seq, err := sequences.NewSequenceFromJSON(phase.SequencePath)
		if err != nil {
			log.Printf("Failed to load sequence: %v", err)
		} else {
			s.sequencePlayer.Play(seq)
		}
	}

	s.initGoal()
}

func (s *PhasesScene) initGoal() {
	phase, _ := s.AppContext().PhaseManager.GetCurrentPhase()
	switch phase.GoalType {
	case RescueSheepType:
		s.goal = &RescueSheepGoal{scene: s}
	case ReactEndpointType:
		s.goal = &ReachEndpointGoal{scene: s}
	case SequenceGoalType:
		s.goal = &phases.SequenceGoal{
			Player:         s.sequencePlayer,
			OnCompleteFunc: s.defaultCompletion,
		}
	case NoGoalType:
		s.goal = &phases.NoGoal{}
	default:
		s.goal = &phases.NoGoal{}
	}
}

func (s *PhasesScene) freezeAllActors() {
	if s.AppContext().ActorManager != nil {
		s.AppContext().ActorManager.ForEach(func(actor actors.ActorEntity) {
			actor.SetImmobile(true)
			actor.SetFreeze(true)
		})
	}
}

func (s *PhasesScene) defaultCompletion() {
	s.isConcludingPhase = true
	s.phaseCompletedDelay = timing.FromDuration(2 * time.Second)
}

func (s *PhasesScene) Update() error {
	if s.pauseScreen != nil && s.allowPause {
		s.pauseScreen.Update()
		if s.pauseScreen.IsPaused() {
			return nil
		}
	}

	if s.sequencePlayer != nil {
		s.sequencePlayer.Update()
	}

	if s.vfxManager != nil {
		s.vfxManager.Update()
	}

	if s.screenFlipper != nil {
		s.screenFlipper.Update()
		if s.screenFlipper.IsFlipping() {
			return nil
		}
	}

	if config.Get().CamDebug {
		s.Camera().CamDebug()
	}

	if s.checkReboot() {
		return nil
	}

	if s.checkPhaseCompleted() {
		s.completePhase()
	}

	// This calls TilemapScene.Update -> BaseScene.Update (handling Schedule)
	if err := s.TilemapScene.Update(); err != nil {
		return err
	}

	s.playBackgroundMusic()

	s.count++

	// Execute bodies updates
	space := s.PhysicsSpace()
	for _, i := range space.Bodies() {
		switch b := i.(type) {
		// ActorEntity case should came first. It can be confused with body.Obstacle
		case platformer.PlatformerActorEntity:
			if err := b.Update(space); err != nil {
				return err
			}
		case items.Item:
			// Remove items marked as removed
			if b.IsRemoved() {
				s.PhysicsSpace().RemoveBody(i)
				continue
			}
			if err := b.Update(space); err != nil {
				return err
			}
		case body.Obstacle:
			continue
		}
	}

	// Remove bodies queued for removal
	space.ProcessRemovals()

	return nil
}

func (s *PhasesScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0xff}) // force black

	// Get tilemap image and draw based on camera
	tilemap, err := s.Tilemap().Image(screen)
	if err != nil {
		log.Fatal(err)
	}
	s.Camera().Draw(tilemap, s.Tilemap().ImageOptions(), screen)

	// Draw bodies based on camera
	space := s.PhysicsSpace()
	for _, b := range space.Bodies() {
		switch sb := b.(type) {
		case platformer.PlatformerActorEntity:
			opts := sb.ImageOptions()
			sb.UpdateImageOptions()
			s.Camera().Draw(sb.Image(), opts, screen)
			if config.Get().CollisionBox {
				s.Camera().DrawCollisionBox(screen, sb)
			}
		case items.Item:
			if sb.IsRemoved() {
				continue
			}
			opts := sb.ImageOptions()
			sb.UpdateImageOptions()
			s.Camera().Draw(sb.Image(), opts, screen)
			if config.Get().CollisionBox {
				s.Camera().DrawCollisionBox(screen, sb)
			}
		case body.Obstacle:
			if config.Get().CollisionBox {
				s.Camera().DrawCollisionBox(screen, sb)
			}
		}
	}

	if s.ShowDrawScreenFlash > 0 {
		screenutil.DrawScreenFlash(screen)
		s.ShowDrawScreenFlash--
	}

	if s.vfxManager != nil {
		s.vfxManager.Draw(screen, s.Camera())
	}

	if s.pauseScreen.IsPaused() {
		s.drawPause(screen)
	}
}

func (s *PhasesScene) Reboot() {
	s.ShowDrawScreenFlash = timing.FromDuration(67 * time.Millisecond) // 4 frames
	s.isRebooting = true
	s.rebootDelay = timing.FromDuration(1 * time.Second) // 60 frames
}

func (s *PhasesScene) OnFinish() {
	s.TilemapScene.OnFinish()
	// Ensure we remove any movement block applied at phase start (for whichever actor is the current player)
	if p, ok := s.AppContext().ActorManager.GetPlayer(); ok {
		p.UnblockMovement()
	}
	s.AppContext().ActorManager.Unregister(s.player)
}

func (s *PhasesScene) endpointTrigger(eventID string) {
	if eventID == "SPIKE" {
		s.player.OnDie()
		return
	}

	s.reachedEndpoint = true

	sheepCarrier, ok := s.player.(gameentitytypes.SheepCarrier)
	if !ok {
		return
	}

	if sheepCarrier.IsCarryingSheep() {
		sheepCarrier.DropSheep()
		s.bodyCounter.sheepRescued++
	}
}

func (s *PhasesScene) playBackgroundMusic() {
	if s.AppContext().Config.NoSound {
		return
	}

	if s.count == 60 {
		if am := s.AppContext().AudioManager; !am.IsPlaying(bgSound) {
			am.PlayMusic(bgSound)
			am.SetVolume(1)
		}
	}

}

func (s *PhasesScene) initTilemap() {
	// Set items position from tilemap
	f := items.NewItemFactory(gameitems.InitItemMap(s.AppContext()))
	scene.InitItems(s.TilemapScene, f)

	// Set enemies position from tilemap
	enemyFactory := enemies.NewEnemyFactory(gameenemies.InitEnemyMap(s.AppContext()))
	scene.InitEnemies(s.TilemapScene, enemyFactory)

	// Set NPCs position from tilemap
	npcFactory := npcs.NewNpcFactory(gamenpcs.InitNpcMap(s.AppContext()))
	scene.InitNPCs(s.TilemapScene, npcFactory)

	s.SetPlayerStartPosition(s.player)
}

func (s *PhasesScene) checkReboot() bool {
	if !s.isRebooting {
		return false
	}

	if s.rebootDelay == 0 {
		s.AppContext().SceneManager.NavigateTo(
			scenestypes.ScenePhaseReboot,
			transition.NewFader(),
			true,
		)
	}

	s.rebootDelay--
	return false
}

func (s *PhasesScene) checkPhaseCompleted() bool {
	if s.isConcludingPhase {
		return true
	}

	if s.goal != nil && s.goal.IsCompleted() {
		s.goal.OnCompletion()
		return true
	}

	return false
}

func (s *PhasesScene) completePhase() {
	if s.phaseCompleted {
		return
	}

	if s.phaseCompletedDelay == 0 {
		s.AppContext().CompleteCurrentPhase(transition.NewFader(), true)
		s.phaseCompleted = true
		return
	}

	s.phaseCompletedDelay--
}

func (s *PhasesScene) drawPause(screen *ebiten.Image) {
	if !s.allowPause || s.pauseScreen == nil || !s.pauseScreen.IsPaused() {
		return
	}

	cfg := config.Get()
	for x := 0; x < cfg.ScreenWidth; x++ {
		for y := 0; y < cfg.ScreenWidth; y++ {
			if x%2 == 0 && y%2 == 0 {
				vector.DrawFilledRect(screen, float32(x), float32(y), 1, 1, color.Black, false)
			}
		}
	}

	speed := 10
	initialW, initialH := cfg.ScreenWidth/4, cfg.ScreenHeight/4
	w := max(min(initialW+s.pauseScreen.Count()*speed, cfg.ScreenWidth/2), 1)
	h := max(min(initialH+s.pauseScreen.Count()*speed, cfg.ScreenHeight/2), 1)
	container := ebiten.NewImage(w, h)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(cfg.ScreenWidth)/2, float64(cfg.ScreenHeight)/2)
	op.GeoM.Translate(-float64(w/2), -float64(h/2))
	container.Fill(color.Black)
	screen.DrawImage(container, op)
}
