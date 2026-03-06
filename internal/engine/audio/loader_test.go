package audio

import (
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestLoadAudioAssetsFromFS(t *testing.T) {
	am := getTestAudioManager()
	wavData := createMinimalWAV()

	mockFS := fstest.MapFS{
		"assets/audio/test.wav": &fstest.MapFile{Data: wavData},
		"assets/audio/test.ogg": &fstest.MapFile{Data: []byte("invalid ogg")}, // Should log error but not crash
		"assets/audio/test.txt": &fstest.MapFile{Data: []byte("text file")},    // Should be ignored
		"assets/audio/subdir/ignored.wav": &fstest.MapFile{Data: wavData},    // Should be ignored (IsDir)
	}

	LoadAudioAssetsFromFS(mockFS, am)

	if _, ok := am.audioPlayers["assets/audio/test.wav"]; !ok {
		t.Error("expected assets/audio/test.wav to be loaded")
	}

	if _, ok := am.audioPlayers["assets/audio/test.txt"]; ok {
		t.Error("expected assets/audio/test.txt to be ignored")
	}
}

func TestLoadAudioAssetsFromFS_EmptyDir(t *testing.T) {
	am := getTestAudioManager()
	mockFS := fstest.MapFS{
		"assets/audio/": &fstest.MapFile{Mode: fs.ModeDir},
	}

	LoadAudioAssetsFromFS(mockFS, am) // should not panic/exit
}
