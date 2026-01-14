package gamesetup

import (
	"flag"

	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
)

const (
	// Celeste is 320 x 180
	ScreenWidth   = 320
	ScreenHeight  = 224
	DefaultVolume = 0.5
	MainFontFace  = "assets/fonts/pressstart2p.ttf"
)

func NewConfig() *config.AppConfig {
	defaultPhysics := config.PhysicsConfig{
		HorizontalInertia:     2.0,
		AirFrictionMultiplier: 0.5,
		AirControlMultiplier:  0.25,
		CoyoteTimeFrames:      6,
		JumpBufferFrames:      6,
		JumpForce:             4,
		JumpCutMultiplier:     0.5,
		UpwardGravity:         4,
		DownwardGravity:       4,
		MaxFallSpeed:          fp16.To16(3),
	}

	cfg := &config.AppConfig{
		ScreenWidth:  ScreenWidth,
		ScreenHeight: ScreenHeight,
		Physics:      defaultPhysics,

		DefaultVolume: DefaultVolume,

		MainFontFace: MainFontFace,
	}

	flag.BoolVar(&cfg.CamDebug, "cam-debug", false, "Enable camera debug")
	flag.BoolVar(&cfg.CollisionBox, "collision-box", false, "Enable collision box debug")
	flag.BoolVar(&cfg.NoSound, "no-sound", false, "Disable game sound")

	return cfg
}
