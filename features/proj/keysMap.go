package proj

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Quit   key.Binding
	Add    key.Binding
	Save   key.Binding
	Delete key.Binding
	Edit   key.Binding
	Open   key.Binding
}

var Keys = KeyMap{
	Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Add:    key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
	Save:   key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "save")),
	Delete: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
	Edit:   key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	Open:   key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open")),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Add, k.Save, k.Delete, k.Edit, k.Open}
}
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
