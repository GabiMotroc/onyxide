package ui

import "strings"

import tea "charm.land/bubbletea/v2"

type Picker struct {
	Title    string
	Items    []string
	Selected int
}

func NewPicker(title string, items []string) Picker {
	return Picker{Title: title, Items: items}
}

func (p Picker) SelectedItem() string {
	if len(p.Items) == 0 {
		return ""
	}
	return p.Items[p.Selected]
}

func (p Picker) Update(msg tea.Msg) (Picker, tea.Cmd) {
	k, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return p, nil
	}
	switch k.String() {
	case "up":
		if p.Selected > 0 {
			p.Selected--
		}
	case "down":
		if p.Selected < len(p.Items)-1 {
			p.Selected++
		}
	case "home":
		p.Selected = 0
	case "end":
		p.Selected = len(p.Items) - 1
	}
	return p, nil
}

func (p Picker) View() string {
	var s strings.Builder
	s.WriteString(DialogTitleStyle.Render(p.Title) + "\n")
	for i, item := range p.Items {
		marker := "  "
		if i == p.Selected {
			marker = "> "
		}
		s.WriteString(marker + item + "\n")
	}
	return DialogStyle.Render(s.String())
}
