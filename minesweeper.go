package main

import (
	"fmt"
	"math/rand"
)

func newField(height, width, mines int) {

	// Make a slice with values ranging from 0 to (height X width) and shuffle it
	sequence := make([]int, height*width)
	for i := 0; i < (height * width); i++ {
		sequence[i] = i
	}
	rand.Shuffle(len(sequence), func(i, j int) {
		sequence[i], sequence[j] = sequence[j], sequence[i]
	})

	// Initialize 2D array representing field
	field := make([][]rune, height)
	for i := range field {
		field[i] = make([]rune, width)
		for j := range field[i] {
			field[i][j] = '0'
		}
	}

	for i := 0; i < mines; i++ {
		v := sequence[i]
		col := v % height // Modulus to get height
		row := v / height // Integer division to get row
		field[row][col] = '☀'
	}

	for _, row := range field {
		fmt.Printf("\n")
		for _, col := range row {
			fmt.Printf("%q", col)
		}
	}
}

func main() {
	newField(10, 10, 20)
}
