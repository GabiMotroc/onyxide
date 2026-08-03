package proj

import (
	"fmt"

	"charm.land/bubbles/v2/table"
)

func toRows(projects []Project) []table.Row {
	rows := make([]table.Row, len(projects))
	for i, p := range projects {
		rows[i] = table.Row{fmt.Sprintf("%d", i+1), p.AppType, p.Location}
	}
	return rows
}

//func (m model) addApp(name string) (model, error) {
//	if ContainsLocation(m.apps, name) {
//		return m, fmt.Errorf("app with name %s already exists", name)
//	}
//	m.apps = append(m.apps, App{Name: name})
//	m.syncTable()
//
//	return m, nil
//}

//func (m model) editApp(name string) (model, error) {
//	if ContainsAppName(m.apps, name) {
//		return m, fmt.Errorf("app with name %s already exists", name)
//	}
//	m.apps[m.editingIndex].Name = name
//
//	return m, nil
//}

func (m *model) syncTable() {
	m.table.SetRows(toRows(m.projects))
}
