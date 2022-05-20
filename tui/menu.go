package tui

import tea "github.com/charmbracelet/bubbletea"

type menuModel struct {
	menuSelection int
	menuChoices   int
	difficulty    int
	inMenu        bool
}

func (m menuModel) view() string {
	return "This is a menu! Press 'enter' to go start a new game."
}

func (m menuModel) update(msg tea.Msg) (menuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.inMenu = false
		}
	}
	return m, nil
}
