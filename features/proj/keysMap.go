package proj

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Quit      key.Binding
	Add       key.Binding
	Save      key.Binding
	Delete    key.Binding
	Edit      key.Binding
	Open      key.Binding
	Back      key.Binding
	Next      key.Binding
	Confirm   key.Binding
	Increment key.Binding
	Decrement key.Binding
}

var Keys = KeyMap{
	Quit:      key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
	Add:       key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
	Save:      key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "save")),
	Delete:    key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
	Edit:      key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	Open:      key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open")),
	Back:      key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Next:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next")),
	Confirm:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
	Increment: key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "increment")),
	Decrement: key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "decrement")),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Add, k.Save, k.Delete, k.Edit, k.Open, k.Increment, k.Decrement}
}
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
