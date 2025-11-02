package store

import (
	"encoding/json"
	"os"
)

const storeFile = "notes.json"

func Save(data interface{}) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(storeFile, d, 0644)
}

func Load() ([]byte, error) {
	data, err := os.ReadFile(storeFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []byte{}, nil
		}
		return nil, err
	}
	return data, nil
}
