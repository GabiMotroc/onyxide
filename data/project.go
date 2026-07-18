package data

import (
	"encoding/json"
	"os"
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
	bytes, err := json.MarshalIndent(apps, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(projectLocation()), 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(projectLocation(), bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadProjects() ([]Project, error) {
	bytes, err := os.ReadFile(projectLocation())
	if os.IsNotExist(err) {
		return []Project{}, nil
	}
	if err != nil {
		return nil, err
	}

	var proj []Project
	err = json.Unmarshal(bytes, &proj)
	if err != nil {
		return nil, err
	}

	return proj, nil
}
