package audio

import (
	"testing"
	"time"

	"github.com/leandroatallah/firefly/internal/engine/data/config"
)

func TestNewAudioManagerRespectsNoSound(t *testing.T) {
	config.Set(&config.AppConfig{NoSound: true})
	am := NewAudioManager()
	if am.Volume() != 0 {
		t.Fatalf("expected volume 0 when NoSound=true, got %f", am.Volume())
	}
}

func TestPlayAndSetVolumeAndPauseAllNoPlayers(t *testing.T) {
	config.Set(&config.AppConfig{NoSound: true})
	am := NewAudioManager()

	if am.PlayMusic("missing") != nil {
		t.Fatalf("expected nil player for missing key")
	}

	am.SetVolume(0.5)
	if am.Volume() != 0.5 {
		t.Fatalf("volume not set")
	}

	am.PauseAll() // no panic
}

func TestFadeOutAllReachesZero(t *testing.T) {
	config.Set(&config.AppConfig{NoSound: false})
	am := NewAudioManager()
	am.SetVolume(1.0)
	am.FadeOutAll(50 * time.Millisecond)
	time.Sleep(150 * time.Millisecond)
	if am.Volume() != 0 {
		t.Fatalf("expected volume 0 after fadeout, got %f", am.Volume())
	}
}
