package proj

import (
	"fmt"
	"onyxide/features/app"
	"onyxide/ui"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
)

func startInteractive(cmd *cobra.Command, args []string) error {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()

	if err != nil {
		return err
	}

	if merr := m.(model).err; merr != nil {
		return merr
	}
	return nil
}

func initialModel() model {
	a, _ := LoadProjects()

	cols := []table.Column{
		{Title: "#", Width: 4},
		{Title: "APP_TYPE", Width: 20},
		{Title: "LOCATION", Width: 45},
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(toRows(a)),
		table.WithFocused(true),
	)

	ti := textinput.New()
	ti.CharLimit = 156
	ti.SetWidth(20)
	ti.SetVirtualCursor(false)

	return model{
		projects:     a,
		table:        t,
		input:        ti,
		editingIndex: -1,
		help:         help.New(),
		keys:         Keys,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		switch m.mode {
		case modeBrowse:
			return m.updateBrowser(msg)
		case modeAdd:
			return m.updatePicker(msg)
		case
			modeEdit:
			return m.updateInput(msg)
		}
	case tea.WindowSizeMsg:
		m.table.SetWidth(msg.Width)
		//m.table.SetHeight(msg.Height - 4) // subtract for help text
	}

	return m, nil
}

func (m model) updateBrowser(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	idx := m.table.Cursor()

	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.Add):
		apps, _ := app.LoadApps()
		if len(apps) == 0 {
			m.err = fmt.Errorf("no apps registered — run onyxide app first")
			return m, nil
		}
		names := make([]string, len(apps))
		for i, a := range apps {
			names[i] = a.Name
		}
		m.picker = ui.NewPicker("Select app", names)
		m.mode = modeAdd
		return m, nil

	case key.Matches(msg, m.keys.Edit):
		//if len(m.projects) > 0 {
		//	m.mode = modeEdit
		//	m.err = nil
		//	m.editingIndex = idx
		//m.input.SetValue(m.projects[idx])
		//cmd := m.input.Focus()
		//	return m, cmd
		//}

	case key.Matches(msg, m.keys.Save):
		m.err = SaveProjects(m.projects)
		return m, tea.Quit

	case key.Matches(msg, m.keys.Delete):
		if len(m.projects) == 0 {
			break
		}
		m.projects = append(m.projects[:idx], m.projects[idx+1:]...)
		if idx >= len(m.projects) && idx > 0 {
			idx--
			m.table.SetCursor(idx)
		}
		m.syncTable()
	case key.Matches(msg, m.keys.Open):
		if len(m.projects) == 0 {
			break
		}
		p := m.projects[idx]
		if err := openProject(p).Start(); err != nil {
			m.err = err
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) updateInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		//m.err = nil
		//name := m.input.Value()
		//if name != "" {
		//	if m.mode == modeEdit {
		//		m, m.err = m.editApp(name)
		//	} else {
		//		m, m.err = m.addApp(name)
		//	}
		//}
		//if m.err == nil {
		//	m.mode = modeBrowse
		//	m.editingIndex = -1
		//	m.input.SetValue("")
		//}
		m.syncTable()
	case "esc":
		m.mode = modeBrowse
		m.err = nil
		m.input.SetValue("")
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() tea.View {
	s := m.table.View()

	switch m.mode {
	case modeAdd:
		return tea.NewView(m.picker.View())

	case modeEdit:
		title := "Add App"
		if m.mode == modeEdit {
			title = "Edit App"
		}

		content := ui.DialogTitleStyle.Render(title) + "\n" + m.input.View()
		if m.err != nil {
			content += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(m.err.Error())
		}
		content += "\n\nenter to confirm • esc to cancel"

		s := ui.DialogStyle.Render(content)
		return tea.NewView(s)

	case modeBrowse:
		s += "\n" + m.help.View(m.keys) + "\n"
	}
	return tea.NewView(s)
}

func (m model) updatePicker(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.pendingApp = m.picker.SelectedItem()
		// TODO: switch to the location-input step
		return m, nil
	case "esc":
		m.mode = modeBrowse
		m.err = nil
		return m, nil
	}
	var cmd tea.Cmd
	m.picker, cmd = m.picker.Update(msg)
	return m, cmd
}

type mode int

const (
	modeBrowse mode = iota
	modeAdd
	modeEdit
)

type model struct {
	projects     []Project
	table        table.Model
	input        textinput.Model
	mode         mode
	editingIndex int
	err          error
	help         help.Model
	keys         KeyMap
	picker       ui.Picker
	pendingApp   string
}
