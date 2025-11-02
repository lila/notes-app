package cache

import (
	"encoding/json"
	"os"
)

const cacheFile = "concept_cache.json"

var conceptCache map[string]string

func init() {
	conceptCache = make(map[string]string)
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return // Cache file doesn't exist yet, that's fine.
		}
		// For other errors, we can log them but continue.
		// A broken cache shouldn't crash the app.
		return
	}
	_ = json.Unmarshal(data, &conceptCache)
}

func Get(key string) (string, bool) {
	name, found := conceptCache[key]
	return name, found
}

func Set(key, name string) {
	conceptCache[key] = name
	data, err := json.MarshalIndent(conceptCache, "", "  ")
	if err != nil {
		return // Don't crash on cache write failure
	}
	_ = os.WriteFile(cacheFile, data, 0644)
}
