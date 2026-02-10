package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Load(path string) (Data, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Data{BestScores: map[string]BestScore{}}, nil
		}
		return Data{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data Data
	if err := decoder.Decode(&data); err != nil {
		return Data{}, err
	}
	if data.BestScores == nil {
		data.BestScores = map[string]BestScore{}
	}
	return data, nil
}

func Save(path string, data Data) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	file, err := os.CreateTemp(dir, "gotype-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), path)
}
