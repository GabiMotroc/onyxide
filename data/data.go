package data

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func getDataDir() string {
	base := os.Getenv("APPDATA")
	if base == "" {
		base = filepath.Join(os.Getenv("HOME"), ".local", "share")
	}
	return filepath.Join(base, "onyxide", "persistent")
}

func createFile(location string) error {
	return writeToFile([]byte("[]"), location)
}

func writeToFile(bytes []byte, location string) error {
	err := os.MkdirAll(filepath.Dir(location), 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(location, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func load[T any](location string) ([]T, error) {
	bytes, err := os.ReadFile(location)
	if os.IsNotExist(err) {
		err := createFile(location)
		if err != nil {
			return nil, err
		}
		return []T{}, nil
	}
	if err != nil {
		return nil, err
	}

	var items []T
	err = json.Unmarshal(bytes, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func save[T any](items []T, location string) error {
	bytes, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	return writeToFile(bytes, location)
}
