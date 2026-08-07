package app

import (
	"fmt"
	"onyxide/ui"

	"charm.land/bubbles/v2/table"
)

func appToRows(apps []App) []table.Row {
	rows := make([]table.Row, len(apps))
	for i, app := range apps {
		rows[i] = table.Row{fmt.Sprintf("%d", i+1), app.Name, ui.BoolIcon(app.IsTerminal)}
	}
	return rows
}

func (m model) addApp(name string) (model, error) {
	if ContainsAppName(m.apps, name) {
		return m, fmt.Errorf("app with name %s already exists", name)
	}
	m.apps = append(m.apps, App{Name: name})
	m.syncTable()

	return m, nil
}

func (m model) editApp(name string) (model, error) {
	if ContainsAppName(m.apps, name) {
		return m, fmt.Errorf("app with name %s already exists", name)
	}
	m.apps[m.editingIndex].Name = name

	return m, nil
}

func (m *model) syncTable() {
	m.table.SetRows(appToRows(m.apps))
	m.dirty = true
}
