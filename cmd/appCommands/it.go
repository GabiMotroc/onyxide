package appCommands

import (
	"fmt"
	"onyxide/data"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
)

var itCmd = &cobra.Command{
	Use:   "it",
	Short: "Manage apps",
	Long:  "Create, list, and manage applications.",
	RunE:  startInteractive,
}

func init() {
	AppCmd.AddCommand(itCmd)
}

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
	a, _ := data.LoadApps()

	ti := textinput.New()
	ti.CharLimit = 156
	ti.SetWidth(20)
	ti.SetVirtualCursor(false)
	return model{
		apps:         a,
		input:        ti,
		editingIndex: -1,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
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
	}

	return m, nil
}

func (m model) updateBrowser(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {

	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.apps)-1 {
			m.cursor++
		}

	case "a":
		m.mode = modeAdd
		m.err = nil
		m.input.SetValue("")
		cmd := m.input.Focus()
		return m, cmd

	case "e":
		if len(m.apps) > 0 {
			m.mode = modeEdit
			m.err = nil
			m.editingIndex = m.cursor
			m.input.SetValue(m.apps[m.cursor].Name)
			cmd := m.input.Focus()
			return m, cmd
		}

	case "s":
		m.err = data.SaveApps(m.apps)
		return m, tea.Quit

	case "d":
		if len(m.apps) == 0 {
			break
		}

		m.apps = append(m.apps[:m.cursor], m.apps[m.cursor+1:]...)
		if m.cursor >= len(m.apps) && m.cursor > 0 {
			m.cursor--
		}
	}

	return m, nil
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
	s := "Manage apps\n\n"

	for i, app := range m.apps {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%d) %s %s\n", i, cursor, app.Name)
	}

	switch m.mode {
	case modeAdd, modeEdit:
		if m.err != nil {
			s += "\n" + m.err.Error()
		}
		s += "\n" + m.input.View() + "\nenter to confirm, esc to cancel"
	case modeBrowse:
		s += "\nPress q to quit. Press a to add. Press s to save. Press d to delete\n"
	}
	return tea.NewView(s)
}

func (m model) addApp(name string) (model, error) {
	if data.ContainsAppName(m.apps, name) {
		return m, fmt.Errorf("app with name %s already exists", name)
	}
	m.apps = append(m.apps, data.App{Name: name})
	m.cursor = len(m.apps) - 1

	return m, nil
}

func (m model) editApp(name string) (model, error) {
	if data.ContainsAppName(m.apps, name) {
		return m, fmt.Errorf("app with name %s already exists", name)
	}
	m.apps[m.editingIndex].Name = name

	return m, nil
}

type mode int

const (
	modeBrowse mode = iota
	modeAdd
	modeEdit
)

type model struct {
	apps         []data.App
	cursor       int
	input        textinput.Model
	mode         mode
	editingIndex int
	err          error
}
