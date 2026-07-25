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
