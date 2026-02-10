package storage

import (
	"os"
	"path/filepath"
)

func DefaultPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "gotype", "state.json"), nil
}
