package proj

import (
	"fmt"

	"charm.land/bubbles/v2/table"
)

func toRows(projects []Project) []table.Row {
	rows := make([]table.Row, len(projects))
	for i, p := range projects {
		rows[i] = table.Row{fmt.Sprintf("%d", i+1), fmt.Sprintf("%d", p.Score), p.AppType, p.Location}
	}
	return rows
}

func (m model) syncTable() model {
	m.table.SetRows(toRows(m.projects))
	m.dirty = true
	return m
}

func appendProject(projects []Project, appType string, loc string) ([]Project, error) {

	if ContainsLocation(projects, loc) {
		return projects, fmt.Errorf("project at %q already exists", loc)
	}

	projects = append(projects, Project{
		AppType:  appType,
		Location: loc,
	})
	return projects, nil
}
