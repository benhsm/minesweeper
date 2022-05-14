package main

import (
	"math/rand"
)

// Constants representing characters used to render certain game elements
const (
	mineRune    = "☀"
	flagRune    = ""
	unknownRune = "?"
)

const mine = -1

const (
	playing = iota
	lost
	won
)

type point struct {
	x int
	y int
}

type tile struct {
	val   int
	state int
}

type mineField [][]tile

type model struct {
	field          mineField
	cursor         point
	gameState      int
	tilesRemaining int
	height         int
	width          int
}

func (m *model) setSize(w, h int) {
	m.width = w
	m.height = h
}

// newField takes dimensions and returns a 2D array
// representing a randomly generated minesweeper playing field.
// each cell is filled with either the value -1 representing a mine, or an
// integer reprsenting the number of adjacent mines.
func newField(height, width, mines int) [][]int {

	// Make a slice with values ranging from 0 to (height X width) and shuffle it
	sequence := make([]int, height*width)
	for i := 0; i < (height * width); i++ {
		sequence[i] = i
	}
	rand.Shuffle(len(sequence), func(i, j int) {
		sequence[i], sequence[j] = sequence[j], sequence[i]
	})

	// Initialize 2D array representing field
	field := make([][]int, height)
	for i := range field {
		field[i] = make([]int, width)
		for j := range field[i] {
			field[i][j] = 0
		}
	}

	// Use the random sequence to mine the empty field
	for i := 0; i < mines; i++ {
		v := sequence[i]
		col := v % width // Modulus to get column position
		row := v / width // Integer division to get row
		field[row][col] = -1
	}

	// Fill tiles with numbers representing adjacent mines
	for i := 0; i < height*width; i++ {
		col := i % width // Modulus to get column position
		row := i / width // Integer division to get row

		if field[row][col] == -1 {
			continue
		}

		val := 0
		for x := -1; x <= 1; x++ {
			for y := -1; y <= 1; y++ {
				posx := col + x
				posy := row + y
				if posx >= 0 && posx <= width-1 {
					if posy >= 0 && posy <= height-1 {
						if field[posy][posx] == -1 {
							val++
						}
					}
				}
			}
		}
		field[row][col] = val
	}

	return field
}

func newMineField(field [][]int) mineField {
	height := len(field)
	width := len(field[0])
	result := make(mineField, height)
	for i := range result {
		result[i] = make([]tile, width)
	}
	for i := 0; i < height*width; i++ {
		col := i % width // Modulus to get height
		row := i / width // Integer division to get row

		result[row][col].val = field[row][col]
		result[row][col].state = hidden
	}

	return result
}

func (m mineField) flagTile(x, y int) {
	switch m[y][x].state {
	case flagged:
		m[y][x].state = hidden
	case hidden:
		m[y][x].state = flagged
	case revealed:
		m[y][x].state = revealed
	}
}

// revealTile reveals a tile on field m and returns true if it's a mine.
// if the tile has no adjacent mines, surrounding tiles are also revealed.
func (m mineField) revealTile(col, row int) int {

	if m[row][col].state == revealed {
		// tile is already revealed and there's nothing to do
		return 0
	}

	m[row][col].state = revealed
	if m[row][col].val == mine {
		// player activated on a mine and died
		return -1
	}

	tilesRevealed := 1 // the tile revealed by the click
	if m[row][col].val == 0 {
		// if the tile is a "0", we need to reveal surrounding tiles, recursing if we encounter another "0"
		height := len(m)
		width := len(m[0])

		for x := -1; x <= 1; x++ {
			for y := -1; y <= 1; y++ {
				posx := col + x
				posy := row + y
				if posx >= 0 && posx <= width-1 {
					if posy >= 0 && posy <= height-1 {
						// Recursively reveal 0 tiles surrounding a 0 tile
						tilesRevealed += m.revealTile(posx, posy)
					}
				}
			}
		}
	}
	return tilesRevealed
}
