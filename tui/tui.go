package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	inGameMenu = iota
	inGame
)

type mainModel struct {
	sessionState  int
	gameComponent gameModel
	menuComponent menuModel
	height        int
	width         int
}

func (m *mainModel) setSize(w, h int) {
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

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
	}
	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	switch m.sessionState {
	case inGame:
		m.gameComponent, cmd = m.gameComponent.update(msg)
		if !m.gameComponent.inGame {
			m.sessionState = inGameMenu
			m.menuComponent.inMenu = true
		}
		return m, cmd
	default:
		m.menuComponent, cmd = m.menuComponent.update(msg)
		if !m.menuComponent.inMenu {
			m.sessionState = inGame
			m.gameComponent.inGame = true
		}
		return m, cmd
	}
}

func (m mainModel) View() string {
	var s string
	switch m.sessionState {
	case inGame:
		s += m.gameComponent.view()
	default:
		s += m.menuComponent.view()
	}
	s = lipgloss.PlaceHorizontal(m.width, lipgloss.Center,
		lipgloss.PlaceVertical(m.height, lipgloss.Center, s))
	return s
}

func NewModel() mainModel {
	return mainModel{
		sessionState:  inGameMenu,
		gameComponent: newGameModel(),
		menuComponent: menuModel{inMenu: true},
	}
}
