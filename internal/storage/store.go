package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// user for loading and saving the data too file , or creating if not 
// exists
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
	// check for the data if it's valid or not
	if err := decoder.Decode(&data); err != nil {
		return Data{}, err
	}
	// if the data is valid but the best scores is nil we need to 
	// initialize it to an empty map
	if data.BestScores == nil {
		data.BestScores = map[string]BestScore{}
	}
	return data, nil
}

// saving data too file -> creating a temp file and then renaming and moving too
// the target path 
func Save(path string, data Data) error {
	// ensure the directory exists
	dir := filepath.Dir(path)

	// create the directory if it doesn't exist, with permissions 0755
	// 755 -> read and execute permissions for everyone, 
	// and write permission for the owner
	if err := os.MkdirAll(dir, 0o755); err != nil { return err }

	// creating a temp file, and remove it after the function returns
	file, err := os.CreateTemp(dir, "gotype-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	// encode the data to the temp file with indentation for readability
	encoder := json.NewEncoder(file)
	// SetIndent adds indentation to the JSON output for better readability.
	encoder.SetIndent("", "  ")


	/// some checks for the encoding and closing the file, 
	// if any error occurs we return it
	if err := encoder.Encode(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}

	// rename the temp file to the target path
	return os.Rename(file.Name(), path)
}
