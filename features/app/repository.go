package app

import (
	"onyxide/store"
	"path/filepath"
	"strings"
)

func repoPath() string {
	return filepath.Join(store.DataDir(), "apps.json")
}

func SaveApps(apps []App) error {
	return store.Save(apps, repoPath())
}

func LoadApps() ([]App, error) {
	return store.Load[App](repoPath())
}

func ContainsAppName(apps []App, name string) bool {
	for _, a := range apps {
		if strings.EqualFold(a.Name, name) {
			return true
		}
	}
	return false
}
