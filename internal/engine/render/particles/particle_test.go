package particles

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/render/camera"
)

func TestParticleUpdate(t *testing.T) {
	config := &Config{
		FrameCount: 2,
		FrameRate:  5,
	}
	p := &Particle{
		X:          10,
		Y:          20,
		VelX:       1,
		VelY:       2,
		Duration:   10,
		Scale:      1.0,
		ScaleSpeed: 0.1,
		Config:     config,
	}

	p.Update()

	if p.X != 11 || p.Y != 22 {
		t.Errorf("expected position (11, 22), got (%f, %f)", p.X, p.Y)
	}
	if p.Duration != 9 {
		t.Errorf("expected duration 9, got %d", p.Duration)
	}
	if p.Scale != 1.1 {
		t.Errorf("expected scale 1.1, got %f", p.Scale)
	}

	// FrameRate is 5. 1st update above made FrameTimer=1.
	// After 4 more updates, FrameTimer=5 -> Frame=1, FrameTimer=0.
	for i := 0; i < 4; i++ {
		p.Update()
	}
	if p.Frame != 1 {
		t.Errorf("expected frame 1 after 5 updates, got %d", p.Frame)
	}

	// 5 more updates to trigger reset to 0
	for i := 0; i < 5; i++ {
		p.Update()
	}
	if p.Frame != 0 {
		t.Errorf("expected frame reset to 0, got %d", p.Frame)
	}
}

func TestParticleIsExpired(t *testing.T) {
	p := &Particle{Duration: 1}
	if p.IsExpired() {
		t.Error("particle should not be expired yet")
	}
	p.Duration = 0
	if !p.IsExpired() {
		t.Error("particle should be expired")
	}
}

func TestSystem(t *testing.T) {
	sys := NewSystem()
	p1 := &Particle{Duration: 3, Config: &Config{}}
	p2 := &Particle{Duration: 2, Config: &Config{}}

	sys.Add(p1)
	sys.Add(p2)

	sys.Update()
	// p1: Dur=2, p2: Dur=1. Both still active.
	if len(sys.particles) != 2 {
		t.Errorf("expected 2 particles, got %d", len(sys.particles))
	}

	sys.Update()
	// p1: Dur=1, p2: Dur=0. p2 should be expired and removed.
	if len(sys.particles) != 1 {
		t.Errorf("expected 1 particle (p1), got %d", len(sys.particles))
	}

	sys.Update()
	// p1: Dur=0. Removed.
	if len(sys.particles) != 0 {
		t.Errorf("expected 0 particles, got %d", len(sys.particles))
	}
}

func TestParticleDrawSmoke(t *testing.T) {
	img := ebiten.NewImage(32, 32)
	img.Fill(color.White)
	config := &Config{
		Image:       img,
		FrameWidth:  16,
		FrameHeight: 16,
		FrameCount:  2,
	}
	p := &Particle{
		X:      10,
		Y:      10,
		Scale:  1.0,
		Config: config,
		Frame:  1,
	}

	screen := ebiten.NewImage(100, 100)
	cam := camera.NewController(0, 0)
	
	// Smoke test
	p.Draw(screen, cam)
	
	// Test with nil image
	p.Config.Image = nil
	p.Draw(screen, cam)
}
