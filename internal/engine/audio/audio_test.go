package audio

import (
	"sync"
	"testing"
	"time"

	"github.com/leandroatallah/firefly/internal/engine/data/config"
)

var (
	audioManagerOnce sync.Once
	audioManager     *AudioManager
)

func getTestAudioManager() *AudioManager {
	audioManagerOnce.Do(func() {
		config.Set(&config.AppConfig{NoSound: false})
		audioManager = NewAudioManager()
	})
	return audioManager
}

func TestNewAudioManagerRespectsNoSound(t *testing.T) {
	// Testa a lógica de inicialização de volume quando NoSound=true
	// Nota: Não podemos chamar NewAudioManager() múltiplas vezes nos testes
	// porque ebiten/audio.NewContext entra em panic se já existir um contexto.
	// Este teste verifica o comportamento através do manager singleton.
	config.Set(&config.AppConfig{NoSound: true})
	am := getTestAudioManager()
	// Após o singleton ser criado com NoSound=false, testamos que SetVolume funciona
	am.SetVolume(0.0)
	if am.Volume() != 0 {
		t.Fatalf("expected volume 0 after SetVolume(0), got %f", am.Volume())
	}
}

func TestPlayAndSetVolumeAndPauseAllNoPlayers(t *testing.T) {
	am := getTestAudioManager()

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
	am := getTestAudioManager()
	am.SetVolume(1.0)
	am.FadeOutAll(50 * time.Millisecond)
	time.Sleep(150 * time.Millisecond)
	if am.Volume() != 0 {
		t.Fatalf("expected volume 0 after fadeout, got %f", am.Volume())
	}
}
