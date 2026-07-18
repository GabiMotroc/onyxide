package data

import (
	"encoding/json"
	"os"
)

type App struct {
	Name string `json:"name"`
}

var appLocation = "persistent/apps.json"

func SaveApps(apps []App) error {
	bytes, err := json.MarshalIndent(apps, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(appLocation, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadApps() ([]App, error) {
	bytes, err := os.ReadFile(appLocation)
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
