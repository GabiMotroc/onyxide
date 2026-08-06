package proj

import (
	"fmt"
	"onyxide/features/app"
	"onyxide/ui"
	"os"

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

	styles := table.DefaultStyles()
	styles.Selected = ui.SelectedStyle
	t := table.New(
		table.WithColumns(cols),
		table.WithRows(toRows(a)),
		table.WithStyles(styles),
		table.WithFocused(true),
	)

	ti := textinput.New()
	ti.Prompt = ""
	ti.CharLimit = 1024
	ti.SetWidth(60)
	ti.SetVirtualCursor(false)

	apps, _ := app.LoadApps()
	appNames := make([]string, len(apps))
	for i, a := range apps {
		appNames[i] = a.Name
	}

	return model{
		projects:     a,
		table:        t,
		input:        ti,
		editingIndex: -1,
		help:         help.New(),
		keys:         Keys,
		picker:       ui.NewPicker(appNames),
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
			return m.updateForm(msg)
		}
	case tea.WindowSizeMsg:
		m.table.SetWidth(msg.Width)
	}

	return m, nil
}

func (m model) updateForm(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Next):
		if m.input.Focused() {
			m.input.Blur()
			return m, nil
		}
		cmd := m.input.Focus()
		return m, cmd

	case key.Matches(msg, m.keys.Back):
		return m.exitForm(), nil

	case key.Matches(msg, m.keys.Confirm):
		if !m.input.Focused() {
			cmd := m.input.Focus()
			return m, cmd
		}
		m.err = nil
		appType := m.picker.SelectedItem()
		loc := m.input.Value()
		if !pathExists(loc) {
			m.err = fmt.Errorf("location %q does not exist", loc)
			return m, nil
		}
		if m.mode == modeEdit {
			m.projects[m.editingIndex].AppType = appType
			m.projects[m.editingIndex].Location = loc
			return m.exitForm(), nil
		}

		m.projects, m.err = appendProject(m.projects, appType, loc)
		if m.err != nil {
			return m, nil
		}

		return m.exitForm(), nil
	}

	if !m.input.Focused() {
		var cmd tea.Cmd
		m.picker, cmd = m.picker.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) exitForm() model {
	m.mode = modeBrowse
	m.editingIndex = -1
	m.input.SetValue("")
	m.input.Blur()
	return m.syncTable()
}

func (m model) updateBrowser(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	idx := m.table.Cursor()

	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.Back):
		return m, tea.Quit
	case key.Matches(msg, m.keys.Add):
		apps := m.picker.Items
		if len(apps) == 0 {
			m.err = fmt.Errorf("no apps registered — run onyxide app first")
			return m, nil
		}
		m.mode = modeAdd
		return m, nil

	case key.Matches(msg, m.keys.Edit):
		if len(m.projects) > 0 {
			m.mode = modeEdit
			m.err = nil
			m.editingIndex = idx
			m.input.SetValue(m.projects[idx].Location)
			selected := 0
			for i, a := range m.picker.Items {
				if m.projects[idx].AppType == a {
					selected = i
					break
				}
			}
			m.picker.Selected = selected
			m.input.Blur()
			return m, nil
		}

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
		return m.syncTable(), nil
	case key.Matches(msg, m.keys.Open):
		if len(m.projects) == 0 {
			break
		}
		p := m.projects[idx]
		c, terminal := openProject(p)
		if terminal {
			return m, tea.ExecProcess(c, nil)
		}
		if err := c.Start(); err != nil {
			m.err = err
		}

	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	s := m.table.View()

	switch m.mode {
	case modeAdd, modeEdit:
		title := "Add Project"
		if m.mode == modeEdit {
			title = "Edit Project"
		}
		isLocFocused := m.input.Focused()
		appTitle, locTitle := ui.SelectedStyle, ui.SelectedStyle
		if !isLocFocused {
			locTitle = ui.UnselectedStyle
		} else {
			appTitle = ui.UnselectedStyle
		}

		content := ui.DialogTitleStyle.Render(title) + "\n" +
			appTitle.Render("App type:") + "\n" + m.picker.List() + "\n" +
			locTitle.Render("Location:") + "\n" + m.input.View()

		if m.err != nil {
			content += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(m.err.Error())
		}
		content += "\n\n" + m.help.ShortHelpView([]key.Binding{
			m.keys.Next,
			m.keys.Confirm,
			m.keys.Back,
		})
		return tea.NewView(ui.DialogStyle.Render(content))

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
	projects     []Project
	table        table.Model
	input        textinput.Model
	mode         mode
	editingIndex int
	err          error
	help         help.Model
	keys         KeyMap
	picker       ui.Picker
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
