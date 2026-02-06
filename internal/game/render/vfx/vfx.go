package vfx

import (
	"encoding/json"
	"image/color"
	"log"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/firefly/internal/engine/data/schemas"
	"github.com/leandroatallah/firefly/internal/engine/render/camera"
	engineparticles "github.com/leandroatallah/firefly/internal/engine/render/particles"
)

type VFXConfig struct {
	Type string `json:"type"`
	schemas.ParticleData
}

// Manager handles all visual effects for the game.
type Manager struct {
	system  *engineparticles.System
	configs map[string]*engineparticles.Config
}

func NewManager() *Manager {
	configs := make(map[string]*engineparticles.Config)

	// Load vfx.json
	path := "assets/particles/vfx.json"
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("failed to load vfx config: %v", err)
	} else {
		var vfxList []VFXConfig
		if err := json.Unmarshal(data, &vfxList); err != nil {
			log.Printf("failed to parse vfx config: %v", err)
		}

		for _, vfx := range vfxList {
			img, _, err := ebitenutil.NewImageFromFile(vfx.Image)
			if err != nil {
				log.Printf("failed to load particle image %s: %v", vfx.Image, err)
				// Fallback to white pixel
				img = ebiten.NewImage(1, 1)
				img.Fill(color.White)
			}

			frameCount := 1
			if vfx.FrameWidth > 0 {
				frameCount = img.Bounds().Dx() / vfx.FrameWidth
			}

			config := &engineparticles.Config{
				Image:       img,
				FrameWidth:  vfx.FrameWidth,
				FrameHeight: vfx.FrameHeight,
				FrameCount:  frameCount,
				FrameRate:   vfx.FrameRate,
			}
			configs[vfx.Type] = config
		}
	}

	return &Manager{
		system:  engineparticles.NewSystem(),
		configs: configs,
	}
}

// spawnPuff creates a puff of particles of the specified type at the given location.
func (m *Manager) spawnPuff(typeKey string, x, y float64, count int, randRange float64) {
	config, ok := m.configs[typeKey]
	if !ok {
		return
	}

	for i := 0; i < count; i++ {
		p := &engineparticles.Particle{
			X:           x,
			Y:           y,
			VelX:        (rand.Float64() - 0.5) * randRange,
			VelY:        (rand.Float64() - 0.5) * randRange,
			Duration:    config.FrameCount * config.FrameRate,
			MaxDuration: config.FrameCount * config.FrameRate,
			Scale:       1.0,
			ScaleSpeed:  0,
			Config:      config,
		}
		m.system.Add(p)
	}
}

// SpawnJumpPuff creates a jump dust effect at the specified location.
// The randRange parameter controls the randomness of the particle velocities.
func (m *Manager) SpawnJumpPuff(x, y float64, count int) {
	m.spawnPuff("jump", x, y, count, 0.1)
}

// SpawnLandingPuff creates a landing dust effect at the specified location.
// The randRange parameter controls the randomness of the particle velocities.
func (m *Manager) SpawnLandingPuff(x, y float64, count int) {
	m.spawnPuff("landing", x, y, count, 0.1)
}

func (m *Manager) Update() {
	m.system.Update()
}

func (m *Manager) Draw(screen *ebiten.Image, cam *camera.Controller) {
	m.system.Draw(screen, cam)
}
