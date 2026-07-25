package proj

import (
	"onyxide/store"
	"path/filepath"
)

func projectLocation() string {
	return filepath.Join(store.DataDir(), "projects.json")
}

func SaveProjects(apps []Project) error {
	return store.Save(apps, projectLocation())
}

func LoadProjects() ([]Project, error) {
	return store.Load[Project](projectLocation())
}

func ContainsLocation(projects []Project, location string) bool {
	for _, p := range projects {
		if p.Location == location {
			return true
		}
	}
	return false
}
