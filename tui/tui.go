package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const banner = `▙▗▌▗                             
▌▘▌▄ ▛▀▖▞▀▖▞▀▘▌  ▌▞▀▖▞▀▖▛▀▖▞▀▖▙▀▖
▌ ▌▐ ▌ ▌▛▀ ▝▀▖▐▐▐ ▛▀ ▛▀ ▙▄▘▛▀ ▌  
▘ ▘▀▘▘ ▘▝▀▘▀▀  ▘▘ ▝▀▘▝▀▘▌  ▝▀▘▘  
`

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
		m.menuComponent.help.Width = msg.Width
		m.gameComponent.help.Width = m.width
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
			switch m.menuComponent.difficulty {
			case beginner:
				m.gameComponent = newGameModel(9, 9, 10)
			case intermediate:
				m.gameComponent = newGameModel(16, 16, 40)
			case expert:
				m.gameComponent = newGameModel(30, 16, 99)
			default:
				m.gameComponent = newGameModel(9, 9, 10)
			}
			m.sessionState = inGame
			m.gameComponent.help.Width = m.width
			m.gameComponent.inGame = true
			cmd = m.gameComponent.stopwatch.Start()
		}
		return m, cmd
	}
}

func (m mainModel) View() string {
	var s string
	s = banner
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
		menuComponent: newMenuModel(),
	}
}
