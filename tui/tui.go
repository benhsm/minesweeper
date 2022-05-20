package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

const (
	inGameMenu = iota
	inGame
)

type mainModel struct {
	sessionState  int
	gameComponent tea.Model
	menuComponent menuModel
}

func (m *gameModel) setSize(w, h int) {
	m.width = w
	m.height = h
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	switch m.sessionState {
	case inGame:
		m.gameComponent, cmd = m.gameComponent.Update(msg)
		return m, cmd
	default:
		m.menuComponent, cmd = m.menuComponent.update(msg)
		if !m.menuComponent.inMenu {
			m.sessionState = inGame
		}
		return m, cmd
	}
}

func (m mainModel) View() string {
	switch m.sessionState {
	case inGame:
		return m.gameComponent.View()
	default:
		return m.menuComponent.view()
	}
}

func NewModel() mainModel {
	return mainModel{
		sessionState:  inGameMenu,
		gameComponent: newGameModel(),
		menuComponent: menuModel{inMenu: true},
	}
}
