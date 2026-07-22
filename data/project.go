package data

import (
	"path/filepath"
)

func projectLocation() string {
	return filepath.Join(getDataDir(), "projects.json")
}

type Project struct {
	AppType  string `json:"apptype"`
	Location string `json:"location"`
}

func SaveProjects(apps []Project) error {
	return save(apps, projectLocation())
}

func LoadProjects() ([]Project, error) {
	return load[Project](projectLocation())
}
