package game

import (
	"math/rand"
)

const Mine = -1

type tile struct {
	Val   int
	State int
}

const (
	Hidden = iota
	Revealed
	Flagged
)

type MineField struct {
	Tiles          [][]tile
	TilesRemaining int
	GameState      int
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

func NewMineField(height, width, mines int) MineField {
	field := newField(height, width, mines)
	tiles := make([][]tile, height)
	tilesRemaining := (height * width) - mines
	result := MineField{Tiles: tiles, TilesRemaining: tilesRemaining}
	for i := range result.Tiles {
		result.Tiles[i] = make([]tile, width)
	}
	for i := 0; i < height*width; i++ {
		col := i % width // Modulus to get height
		row := i / width // Integer division to get row

		result.Tiles[row][col].Val = field[row][col]
		result.Tiles[row][col].State = Hidden
	}

	return result
}

func (m MineField) FlagTile(x, y int) {
	switch m.Tiles[y][x].State {
	case Flagged:
		m.Tiles[y][x].State = Hidden
	case Hidden:
		m.Tiles[y][x].State = Flagged
	case Revealed:
		m.Tiles[y][x].State = Revealed
	}
}

// revealTile reveals a tile on field m and returns true if it's a mine.
// if the tile has no adjacent mines, surrounding tiles are also revealed.
func (m MineField) RevealTile(col, row int) int {

	if m.Tiles[row][col].State == Revealed {
		// tile is already revealed and there's nothing to do
		return 0
	}

	m.Tiles[row][col].State = Revealed
	if m.Tiles[row][col].Val == Mine {
		// player activated on a mine and died
		return -1
	}

	tilesRevealed := 1 // the tile revealed by the click
	if m.Tiles[row][col].Val == 0 {
		// if the tile is a "0", we need to reveal surrounding tiles, recursing if we encounter another "0"
		height := len(m.Tiles)
		width := len(m.Tiles[0])

		for x := -1; x <= 1; x++ {
			for y := -1; y <= 1; y++ {
				posx := col + x
				posy := row + y
				if posx >= 0 && posx <= width-1 {
					if posy >= 0 && posy <= height-1 {
						// Recursively reveal 0 tiles surrounding a 0 tile
						tilesRevealed += m.RevealTile(posx, posy)
					}
				}
			}
		}
	}
	return tilesRevealed
}
