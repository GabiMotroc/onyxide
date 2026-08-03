package proj

import (
	"onyxide/store"
	"path/filepath"
	"strings"
)

func projectLocation() string {
	return filepath.Join(store.DataDir(), "projects.json")
}

func SaveProjects(projects []Project) error {
	return store.Save(projects, projectLocation())
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

func FirstMatchingProject(location string, projects []Project) (bool, Project) {
	for _, project := range projects {
		if isLocationMatching(project, location) {
			return true, project
		}
	}
	return false, Project{}
}

func AllMatchingProjects(location string, projects []Project) []Project {
	matches := make([]Project, 0, len(projects))
	for _, p := range projects {
		if isLocationMatching(p, location) {
			matches = append(matches, p)
		}
	}
	return matches
}

func RemoveMatching(location string, projects []Project) []Project {
	filtered := make([]Project, 0, len(projects))
	for _, p := range projects {
		if !isLocationMatching(p, location) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func isLocationMatching(project Project, location string) bool {
	return strings.Contains(
		strings.ToLower(project.Location),
		strings.ToLower(location),
	)
}
