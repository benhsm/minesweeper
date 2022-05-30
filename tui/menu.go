package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type menuKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Help   key.Binding
	Quit   key.Binding
	Select key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k menuKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k menuKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit, k.Select},      // second column
	}
}

var menuKeys = menuKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k", "w"),
		key.WithHelp("↑/k/w", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j", "s"),
		key.WithHelp("↓/j/s", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h", "a"),
		key.WithHelp("←/h/a", "change setting"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l", "d"),
		key.WithHelp("→/l/d", "change setting"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "select"),
	),
}

var menuSelectStyle = lipgloss.NewStyle().Bold(true)

const (
	beginner = iota
	intermediate
	expert
)

const (
	difficultyOption = iota
	playButton
)

type menuModel struct {
	keys          menuKeyMap
	help          help.Model
	menuSelection int
	menuChoices   []string
	difficulty    int
	inMenu        bool
}

func newMenuModel() menuModel {
	return menuModel{
		keys:        menuKeys,
		help:        help.New(),
		menuChoices: []string{"Difficulty: ", "Play"},
		difficulty:  beginner,
		inMenu:      true,
	}
}

func (m menuModel) view() string {
	var s string
	for i, c := range m.menuChoices {
		if i == difficultyOption {
			switch m.difficulty {
			case beginner:
				c += "Beginner"
			case intermediate:
				c += "Intermediate"
			case expert:
				c += "Expert"
			}
		}
		if i == m.menuSelection {
			c = "<< " + menuSelectStyle.Render(c) + " >> "
		}
		s += fmt.Sprintf("\n%s\n", c)
	}
	s += "\n" + m.help.View(m.keys)
	return s
}

func (m menuModel) update(msg tea.Msg) (menuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.menuSelection > 0 {
				m.menuSelection--
			} else {
				m.menuSelection = len(m.menuChoices) - 1
			}
		case key.Matches(msg, m.keys.Down):
			if m.menuSelection < len(m.menuChoices)-1 {
				m.menuSelection++
			} else {
				m.menuSelection = 0
			}

		case key.Matches(msg, m.keys.Left):
			if m.menuSelection == difficultyOption {
				if m.difficulty > beginner {
					m.difficulty--
				} else {
					m.difficulty = expert
				}
			}
		case key.Matches(msg, m.keys.Right):
			if m.menuSelection == difficultyOption {
				if m.difficulty < expert {
					m.difficulty++
				} else {
					m.difficulty = beginner
				}
			}
		case key.Matches(msg, m.keys.Select):
			if m.menuSelection == playButton {
				m.inMenu = false
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}
	return m, nil
}
