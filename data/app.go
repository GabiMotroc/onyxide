package data

import (
	"path/filepath"
	"strings"
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

func ContainsAppName(apps []App, name string) bool {
	for _, a := range apps {
		if strings.EqualFold(a.Name, name) {
			return true
		}
	}
	return false
}
