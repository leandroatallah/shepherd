package gamescenephases

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/enemies"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/npcs"
	"github.com/leandroatallah/firefly/internal/engine/entity/items"
	"github.com/leandroatallah/firefly/internal/engine/event"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/leandroatallah/firefly/internal/engine/scene"
	"github.com/leandroatallah/firefly/internal/engine/scene/transition"
	gameenemies "github.com/leandroatallah/firefly/internal/game/entity/actors/enemies"
	gamenpcs "github.com/leandroatallah/firefly/internal/game/entity/actors/npcs"
	gameitems "github.com/leandroatallah/firefly/internal/game/entity/items"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
	"github.com/leandroatallah/firefly/internal/game/events"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
)

const (
	bgSound = "assets/audio/Sketchbook.ogg"
)

type PhasesScene struct {
	scene.TilemapScene
	count       int
	player      gameentitytypes.PlatformerActorEntity
	mainText    *font.FontText
	bodyCounter *BodyCounter

	// Complete phase
	isConcludingPhase   bool
	phaseCompleted      bool
	phaseCompletedDelay int

	// Reboot
	isRebooting bool
	rebootDelay int

	// UI effects
	ShowDrawScreenFlash int

	screenFlipper *scene.ScreenFlipper
}

func NewPhasesScene(context *app.AppContext) *PhasesScene {
	mainText, err := font.NewFontText(config.Get().MainFontFace)
	if err != nil {
		log.Fatal(err)
	}
	tilemapScene := scene.NewTilemapScene(context)
	scene := PhasesScene{
		TilemapScene: *tilemapScene,
		mainText:     mainText,
		bodyCounter:  &BodyCounter{},
	}
	scene.SetAppContext(context)

	// Subscribe events
	context.EventManager.Subscribe(events.CharacterDiedEventType, func(e event.Event) {
		scene.Reboot()
	})

	return &scene
}

func (s *PhasesScene) OnStart() {
	s.TilemapScene.OnStart()
	s.count = 0

	// Create player and register to space and context
	p, err := createPlayer(s.AppContext(), gameentitytypes.ShepherdPlayerType)
	if err != nil {
		log.Fatal(err)
	}
	s.player = p
	s.AppContext().ActorManager.Register(s.player)
	s.PhysicsSpace().AddBody(s.player)

	s.initTilemap()

	// After init bodies, set body counter
	s.bodyCounter.setBodyCounter(s.PhysicsSpace())

	s.PhysicsSpace().Bodies()

	// Init camera target
	s.SetCameraConfig(scene.CameraConfig{Mode: scene.CameraModeFixed})
	s.Camera().SetFollowTarget(s.player)
	s.Camera().SetCenter(float64(config.Get().ScreenWidth)/2, float64(config.Get().ScreenHeight)/2)

	// Init collisions bodies and touch trigger for endpoints
	s.Tilemap().CreateCollisionBodies(s.PhysicsSpace(), func(id string) body.Touchable {
		return bodyphysics.NewTouchTrigger(func() {
			s.endpointTrigger(id)
		}, s.player)
	})

	// Init screen flipper
	s.screenFlipper = scene.NewScreenFlipper(s.Camera(), s.player, s.Tilemap())
	s.screenFlipper.OnFlipStart = func() {
		s.player.SetImmobile(true)
	}
	s.screenFlipper.OnFlipFinish = func() {
		s.player.SetImmobile(false)
	}
}

func (s *PhasesScene) Update() error {
	if s.screenFlipper != nil {
		s.screenFlipper.Update()
		if s.screenFlipper.IsFlipping() {
			return nil
		}
	}

	if config.Get().CamDebug {
		s.CamDebug()
	}

	if s.checkReboot() {
		return nil
	}

	if s.checkPhaseCompleted() {
		s.completePhase()
	}

	s.TilemapScene.Update() // Update the camera if in follow mode

	s.playBackgroundMusic()

	s.count++

	// Execute bodies updates
	space := s.PhysicsSpace()
	for _, i := range space.Bodies() {
		switch b := i.(type) {
		// ActorEntity case should came first. It can be confused with body.Obstacle
		case gameentitytypes.PlatformerActorEntity:
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
		case gameentitytypes.PlatformerActorEntity:
			opts := sb.ImageOptions()
			sb.UpdateImageOptions()
			if config.Get().CollisionBox {
				s.Camera().Draw(sb.ImageCollisionBox(), opts, screen)
			} else {
				s.Camera().Draw(sb.Image(), opts, screen)
			}
		case items.Item:
			if sb.IsRemoved() {
				continue
			}
			opts := sb.ImageOptions()
			sb.UpdateImageOptions()
			if config.Get().CollisionBox {
				s.Camera().Draw(sb.ImageCollisionBox(), opts, screen)
			} else {
				s.Camera().Draw(sb.Image(), opts, screen)
			}
		case body.Obstacle:
			if config.Get().CollisionBox {
				opts := sb.ImageOptions()
				sb.UpdateImageOptions()
				s.Camera().Draw(sb.ImageCollisionBox(), opts, screen)
			}
		}
	}

	if s.ShowDrawScreenFlash > 0 {
		DrawScreenFlash(screen)
		s.ShowDrawScreenFlash--
	}
}

func (s *PhasesScene) Reboot() {
	s.ShowDrawScreenFlash = 4 // frames
	s.isRebooting = true
	s.rebootDelay = 60 // frames
}

func (s *PhasesScene) OnFinish() {
	s.TilemapScene.OnFinish()
	s.AppContext().ActorManager.Unregister(s.player)
}

func (s *PhasesScene) endpointTrigger(eventID string) {
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
			am.SetVolume(0.25)
		}
	}

}

func (s *PhasesScene) initTilemap() {
	// Set items map to factory creation process
	itemsMap := map[int]items.ItemType{
		0: gameitems.CollectibleCoinType,
	}

	// Set items position from tilemap
	f := items.NewItemFactory(gameitems.InitItemMap(s.AppContext()))
	s.InitItems(itemsMap, f)

	// Set enemies position from tilemap
	enemyFactory := enemies.NewEnemyFactory(gameenemies.InitEnemyMap(s.AppContext()))
	scene.InitEnemies(&s.TilemapScene, enemyFactory)

	// Set NPCs position from tilemap
	npcFactory := npcs.NewNpcFactory(gamenpcs.InitNpcMap(s.AppContext()))
	scene.InitNPCs(&s.TilemapScene, npcFactory)

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

// TODO: Define phase objective. For now, it checks if all sheeps were rescued.
func (s *PhasesScene) checkPhaseCompleted() bool {
	if s.isConcludingPhase {
		return true
	}

	if s.bodyCounter.sheep > s.bodyCounter.sheepRescued {
		return false
	}

	s.player.SetImmobile(true)
	s.Audiomanager().FadeOut(bgSound, time.Second)
	s.isConcludingPhase = true
	s.phaseCompletedDelay = 120 // frames
	return true
}

func (s *PhasesScene) completePhase() {
	if s.phaseCompleted {
		return
	}

	if s.phaseCompletedDelay == 0 {
		s.AppContext().PhaseManager.AdvanceToNextPhase()
		s.AppContext().SceneManager.NavigateTo(
			scenestypes.ScenePhases,
			transition.NewFader(),
			true,
		)
		s.phaseCompleted = true
		return
	}

	s.phaseCompletedDelay--
}
