package data

import (
	"path/filepath"
)

type App struct {
	Name string `json:"name"`
}

func appLocation() string {
	return filepath.Join(getDataDir(), "apps.json")
}

func SaveApps(apps []App) error {
	return save(apps, appLocation())
}

func LoadApps() ([]App, error) {
	return load[App](appLocation())
}
