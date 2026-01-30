package scene

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/render/camera"
	"github.com/leandroatallah/firefly/internal/engine/render/tilemap"
)

type FlipType int

const (
	FlipTypeSmooth FlipType = iota
	FlipTypeInstant
)

type ScreenFlipper struct {
	cam     *camera.Controller
	player  body.Movable
	tilemap *tilemap.Tilemap

	screenWidth  float64
	screenHeight float64

	// Configuration
	FlipStrategy       func(dx, dy int) FlipType
	PlayerPushDistance float64 // pixels

	// State
	isFlipping    bool
	flipSourceX   float64
	flipSourceY   float64
	flipTargetX   float64
	flipTargetY   float64
	playerSourceX float64
	playerTargetX float64
	playerSourceY float64
	playerTargetY float64
	flipProgress  float64

	// Hooks
	OnFlipStart  func()
	OnFlipFinish func()
}

func NewScreenFlipper(cam *camera.Controller, player body.Movable, tm *tilemap.Tilemap) *ScreenFlipper {
	cfg := config.Get()
	return &ScreenFlipper{
		cam:                cam,
		player:             player,
		tilemap:            tm,
		screenWidth:        float64(cfg.ScreenWidth),
		screenHeight:       float64(cfg.ScreenHeight),
		PlayerPushDistance: 0.0,
	}
}

func (sf *ScreenFlipper) IsFlipping() bool {
	return sf.isFlipping
}

func (sf *ScreenFlipper) Update() {
	if sf.isFlipping {
		sf.updateFlip()
		return
	}

	sf.checkTrigger()
}

func (sf *ScreenFlipper) checkTrigger() {
	if sf.player == nil || sf.cam == nil || sf.tilemap == nil {
		return
	}

	// Current Camera View
	camX, camY := sf.cam.Kamera().Center()
	left := camX - sf.screenWidth/2
	right := camX + sf.screenWidth/2
	top := camY - sf.screenHeight/2
	bottom := camY + sf.screenHeight/2

	// Player Position (use center or min?)
	// Using Min for trigger to ensure they are fully out or touching edge?
	// User said: "when the player walks out of the screen"
	// Let's use the player's center for smoother detection or Min/Max for strict bounds.
	// Using Center is safer to avoid accidental triggers near edge.
	// Actually, usually it's when the player's leading edge crosses the boundary.
	px, py := sf.player.GetPositionMin()
	w, h := sf.player.GetShape().Width(), sf.player.GetShape().Height()
	pCenterX := float64(px) + float64(w)/2
	pCenterY := float64(py) + float64(h)/2

	// Thresholds (small buffer to avoid floating point issues)
	buffer := 2.0

	if pCenterX < left-buffer {
		sf.triggerFlip(-1, 0)
	} else if pCenterX > right+buffer {
		sf.triggerFlip(1, 0)
	} else if pCenterY < top-buffer {
		sf.triggerFlip(0, -1)
	} else if pCenterY > bottom+buffer {
		sf.triggerFlip(0, 1)
	}
}

func (sf *ScreenFlipper) triggerFlip(dx, dy int) {
	camX, camY := sf.cam.Kamera().Center()

	targetCamX := camX + float64(dx)*sf.screenWidth
	targetCamY := camY + float64(dy)*sf.screenHeight

	// Validate against Map Bounds
	// Assume layer 0 defines map size
	if len(sf.tilemap.Layers) == 0 {
		return
	}
	mapW := float64(sf.tilemap.Layers[0].Width * sf.tilemap.Tilewidth)
	mapH := float64(sf.tilemap.Layers[0].Height * sf.tilemap.Tileheight)

	// Check if target center is valid (should be within map bounds)
	// Actually, the camera view should be within map bounds.
	// Target View: [targetCamX - W/2, targetCamX + W/2]
	if targetCamX-sf.screenWidth/2 < 0 || targetCamX+sf.screenWidth/2 > mapW {
		log.Printf("Cannot flip: Target X out of bounds: %f", targetCamX)
		return
	}
	if targetCamY-sf.screenHeight/2 < 0 || targetCamY+sf.screenHeight/2 > mapH {
		log.Printf("Cannot flip: Target Y out of bounds: %f", targetCamY)
		return
	}

	// Calculate Player Target
	px, py := sf.player.GetPositionMin()
	playerSourceX := float64(px)
	playerSourceY := float64(py)
	playerTargetX := playerSourceX
	playerTargetY := playerSourceY

	// Push player into new screen
	if dx != 0 {
		playerTargetX += float64(dx) * sf.PlayerPushDistance
	}
	if dy != 0 {
		playerTargetY += float64(dy) * sf.PlayerPushDistance
	}

	// Determine Strategy
	flipType := FlipTypeSmooth
	if sf.FlipStrategy != nil {
		flipType = sf.FlipStrategy(dx, dy)
	}

	if flipType == FlipTypeInstant {
		if sf.OnFlipStart != nil {
			sf.OnFlipStart()
		}
		sf.cam.SetCenter(targetCamX, targetCamY)
		sf.player.SetPosition(int(playerTargetX), int(playerTargetY))
		if sf.OnFlipFinish != nil {
			sf.OnFlipFinish()
		}
		return
	}

	// Start Flip
	sf.isFlipping = true
	sf.flipSourceX = camX
	sf.flipSourceY = camY
	sf.flipTargetX = targetCamX
	sf.flipTargetY = targetCamY
	sf.flipProgress = 0

	// Player Transition
	sf.playerSourceX = playerSourceX
	sf.playerSourceY = playerSourceY
	sf.playerTargetX = playerTargetX
	sf.playerTargetY = playerTargetY

	if sf.OnFlipStart != nil {
		sf.OnFlipStart()
	}
}

func (sf *ScreenFlipper) updateFlip() {
	cfg := config.Get()
	speed := cfg.ScreenFlipSpeed
	if speed <= 0 {
		speed = 1.0 / 60.0 // Default fallback
	}

	sf.flipProgress += speed

	if sf.flipProgress >= 1.0 {
		sf.finishFlip()
		return
	}

	// Linear interpolation
	t := sf.flipProgress
	newCamX := sf.flipSourceX + (sf.flipTargetX-sf.flipSourceX)*t
	newCamY := sf.flipSourceY + (sf.flipTargetY-sf.flipSourceY)*t
	sf.cam.SetCenter(newCamX, newCamY)

	newPX := sf.playerSourceX + (sf.playerTargetX-sf.playerSourceX)*t
	newPY := sf.playerSourceY + (sf.playerTargetY-sf.playerSourceY)*t
	sf.player.SetPosition(int(newPX), int(newPY))
}

func (sf *ScreenFlipper) finishFlip() {
	sf.flipProgress = 1.0
	sf.isFlipping = false

	if sf.OnFlipFinish != nil {
		sf.OnFlipFinish()
	}

	// Force exact position
	sf.cam.SetCenter(sf.flipTargetX, sf.flipTargetY)
	sf.player.SetPosition(int(sf.playerTargetX), int(sf.playerTargetY))
}
