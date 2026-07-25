package app

import (
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
	a, _ := LoadApps()

	cols := []table.Column{
		{Title: "#", Width: 4},
		{Title: "NAME", Width: 20},
		{Title: "IS_TERMINAL", Width: 12},
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(appToRows(a)),
		table.WithFocused(true),
	)

	ti := textinput.New()
	ti.CharLimit = 156
	ti.SetWidth(20)
	ti.SetVirtualCursor(false)

	return model{
		apps:         a,
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
		case modeAdd, modeEdit:
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
		m.mode = modeAdd
		m.err = nil
		m.input.SetValue("")
		cmd := m.input.Focus()
		return m, cmd

	case key.Matches(msg, m.keys.Edit):
		if len(m.apps) > 0 {
			m.mode = modeEdit
			m.err = nil
			m.editingIndex = idx
			m.input.SetValue(m.apps[idx].Name)
			cmd := m.input.Focus()
			return m, cmd
		}

	case key.Matches(msg, m.keys.Save):
		m.err = SaveApps(m.apps)
		return m, tea.Quit

	case key.Matches(msg, m.keys.Delete):
		if len(m.apps) == 0 {
			break
		}
		m.apps, _ = RemoveByName(m.apps, m.apps[idx].Name)
		if idx >= len(m.apps) && idx > 0 {
			idx--
			m.table.SetCursor(idx)
		}
		m.syncTable()

	case key.Matches(msg, m.keys.ToggleTerminal):
		if len(m.apps) == 0 {
			break
		}
		m.apps[idx].IsTerminal = !m.apps[idx].IsTerminal
		m.syncTable()
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) updateInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.err = nil
		name := m.input.Value()
		if name != "" {
			if m.mode == modeEdit {
				m, m.err = m.editApp(name)
			} else {
				m, m.err = m.addApp(name)
			}
		}
		if m.err == nil {
			m.mode = modeBrowse
			m.editingIndex = -1
			m.input.SetValue("")
		}
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
	case modeAdd, modeEdit:
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

type mode int

const (
	modeBrowse mode = iota
	modeAdd
	modeEdit
)

type model struct {
	apps         []App
	table        table.Model
	input        textinput.Model
	mode         mode
	editingIndex int
	err          error
	help         help.Model
	keys         KeyMap
}
