package app

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Quit           key.Binding
	Add            key.Binding
	Save           key.Binding
	Delete         key.Binding
	ToggleTerminal key.Binding
	Edit           key.Binding
	Back           key.Binding
}

var Keys = KeyMap{
	Quit:           key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Add:            key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
	Save:           key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "save")),
	Delete:         key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
	ToggleTerminal: key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "toggle")),
	Edit:           key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	Back:           key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Add, k.Save, k.Delete, k.ToggleTerminal, k.Edit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
