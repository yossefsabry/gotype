package storage

import (
	"os"
	"path/filepath"
)

// using os for getting the user config path on operation system
func DefaultPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "gotype", "state.json"), nil
}
