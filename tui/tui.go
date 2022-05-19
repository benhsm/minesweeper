package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type mainModel struct {
	menuSelection int
	sessionState  int
	game          gameModel
}

func (m *gameModel) setSize(w, h int) {
	m.width = w
	m.height = h
}

func (m mainModel) Init() tea.Cmd {
	return nil
}
func (m gameModel) Init() tea.Cmd {
	return nil
}

func (m gameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	switch m.field.GameState {
	case lost, won:
		return updateGameOver(msg, m)
	default:
		return updateGameLoop(msg, m)
	}
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	return m.game.Update(msg)
}

func (m mainModel) View() string {
	return m.game.View()
}

func NewModel() mainModel {
	return mainModel{
		menuSelection: 0,
		sessionState:  0,
		game:          newGameModel(),
	}
}
