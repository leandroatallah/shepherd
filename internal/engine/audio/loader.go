package audio

import (
	"io/fs"
	"log"
	"strings"
)

// LoadAudioAssetsFromFS is a helper function to load all audio files from an fs.FS.
func LoadAudioAssetsFromFS(assets fs.FS, am *AudioManager) {
	dir := "assets/audio"
	files, err := fs.ReadDir(assets, dir)
	if err != nil {
		log.Fatalf("error reading embedded audio dir: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		// Filter for supported audio types
		if !(strings.HasSuffix(fileName, ".ogg") || strings.HasSuffix(fileName, ".wav") || strings.HasSuffix(fileName, ".mp3")) {
			continue
		}

		fullPath := dir + "/" + fileName
		data, err := fs.ReadFile(assets, fullPath)
		if err != nil {
			log.Printf("failed to read embedded file %s: %v", fullPath, err)
			continue
		}

		// Use the existing Add method to process and store the player.
		am.Add(dir+"/"+fileName, data)
	}
}
