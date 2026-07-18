package data

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type App struct {
	Name string `json:"name"`
}

func appLocation() string {
	return filepath.Join(getDataDir(), "apps.json")
}

func SaveApps(apps []App) error {
	bytes, err := json.MarshalIndent(apps, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(appLocation()), 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(appLocation(), bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadApps() ([]App, error) {
	bytes, err := os.ReadFile(appLocation())
	if err != nil {
		return nil, err
	}

	var apps []App
	err = json.Unmarshal(bytes, &apps)
	if err != nil {
		return nil, err
	}

	return apps, nil
}
