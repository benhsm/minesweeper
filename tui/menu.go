package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
	menuSelection int
	menuChoices   []string
	difficulty    int
	inMenu        bool
}

func newMenuModel() menuModel {
	return menuModel{
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
			c = menuSelectStyle.Render(c)
		}
		s += fmt.Sprintf("%s\n", c)
	}
	return s
}

func (m menuModel) update(msg tea.Msg) (menuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k":
			if m.menuSelection > 0 {
				m.menuSelection--
			} else {
				m.menuSelection = len(m.menuChoices) - 1
			}
		case "j":
			if m.menuSelection < len(m.menuChoices)-1 {
				m.menuSelection++
			} else {
				m.menuSelection = 0
			}

		case "h":
			if m.menuSelection == difficultyOption {
				if m.difficulty > beginner {
					m.difficulty--
				}
			} else {
				m.difficulty = expert
			}
		case "l":
			if m.menuSelection == difficultyOption {
				if m.difficulty < expert {
					m.difficulty++
				} else {
					m.difficulty = beginner
				}
			}
		case "enter":
			if m.menuSelection == playButton {
				m.inMenu = false
			}
		}
	}
	return m, nil
}
