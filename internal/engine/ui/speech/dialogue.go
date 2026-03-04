package speech

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Manager handles the display of dialogue and speech bubbles.
type Manager struct {
	speeches        map[string]Speech
	activeSpeech    string
	isSpeaking      bool
	currentText     string
	lines           []string
	currentLine     int
	waitingForInput bool
}

const (
	BubbleSpeechID = "bubble"
	StorySpeechID  = "story"
)

// NewManager creates a new dialogue manager.
func NewManager(s ...Speech) *Manager {
	m := &Manager{
		speeches: make(map[string]Speech),
	}
	for _, s := range s {
		m.AddSpeech(s)
	}
	return m
}

func (m *Manager) AddSpeech(s Speech) {
	m.speeches[s.ID()] = s
}

func (m *Manager) SetSpeech(id string) {
	if _, ok := m.speeches[id]; ok {
		m.activeSpeech = id
	}
}

func (m *Manager) getActiveSpeech() Speech {
	return m.speeches[m.activeSpeech]
}

func (m *Manager) SetActiveSpeech(id string) {
	m.SetSpeech(id)
}

// ShowMessages displays a list of messages.
func (m *Manager) ShowMessages(lines []string, position string, speed int) {
	if len(lines) == 0 {
		return
	}
	s := m.getActiveSpeech()
	m.lines = lines
	m.currentLine = 0
	m.isSpeaking = true
	m.waitingForInput = false
	s.ResetText()
	s.SetPosition(position)
	if speed > 0 {
		s.SetSpeed(speed)
	} else {
		// Default speed if not specified
		s.SetSpeed(4)
	}
	s.Show()
}

// IsSpeaking returns true if the dialogue manager is currently displaying a message.
func (m *Manager) IsSpeaking() bool {
	return m.isSpeaking
}

// Update updates the dialogue state. It handles input for proceeding.
func (m *Manager) Update() error {
	if !m.isSpeaking {
		return nil
	}

	s := m.getActiveSpeech()
	if err := s.Update(); err != nil {
		return err
	}

	if s.IsSpellingComplete() && !m.waitingForInput {
		m.waitingForInput = true
	}

	if m.waitingForInput {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			m.currentLine++
			if m.currentLine >= len(m.lines) {
				s.Hide()
				m.isSpeaking = false
			} else {
				s.ResetText()
				m.waitingForInput = false
			}
		}
	}
	return nil
}

// Draw draws the speech bubble if it's active.
func (m *Manager) Draw(screen *ebiten.Image) {
	if !m.isSpeaking {
		return
	}

	if m.currentLine < len(m.lines) {
		m.getActiveSpeech().Draw(screen, m.lines[m.currentLine])
	}
}
