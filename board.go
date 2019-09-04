package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"text/tabwriter"
	"time"
)

type cell struct {
	value rune
	show  bool
	mark  bool
}

//Increases the count of the cells around a mine
func (b board) increaseCount(row, col int) {
	height := len(b)
	width := len(b[0])

	for _, direction := range directions {
		rowCheck := row + direction[0]
		colCheck := col + direction[1]
		if inbound(rowCheck, colCheck, height, width) {
			val := b[rowCheck][colCheck].value
			if val != 'm' {
				b[rowCheck][colCheck].value = rune(int(val) + 1)
			}
		}
	}
}

//Randomly places mines across the input board.
func (b board) placeBombs() int {
	//maxMines is 10% of the number of cells
	maxMines := int(float32(len(b)) * float32(len(b[0])) * .1)
	numMines := 0
	rand.Seed(time.Now().UnixNano())
	for numMines < maxMines {
		//not likely to take more than a couple passes
		for i := range b {
			for j, c := range b[i] {
				if numMines >= maxMines {
					return numMines
				}
				if c.value == 'm' {
					continue
				}
				if rand.Intn(100) < 10 {
					b[i][j].value = 'm'
					b.increaseCount(i, j)
					numMines++
				}
			}
		}
	}
	return numMines
}

//newBoard returns a Board of specified size and the number of mines
func newBoard(width, height int) (board, int) {

	b := make(board, height)
	for i := range b {
		b[i] = make([]cell, width)
	}

	numMines := b.placeBombs()
	return b, numMines
}

//center text
func center(s rune, w int) string {
	if s == 'm' {
		return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*c", (w+1)/2, s))
	}

	return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*v", (w+1)/2, s))
}

//printBoard prints out the board
func (b board) printBoard() {
	clear()
	w := tabwriter.NewWriter(os.Stdout, 5, 1, 0, ' ', 0)
	fmt.Fprint(w, "\t") //spacing
	for i := range b[0] {
		fmt.Fprint(w, center(rune(i), 5), "\t")
	}
	fmt.Fprint(w, "col\n")
	for i := range b {
		fmt.Fprint(w, center(rune(i), 5), "\t")
		for j := range b[0] {
			if b[i][j].show {
				fmt.Fprint(w, "[", center(b[i][j].value, 3), "]\t")
			} else {
				if b[i][j].mark {
					fmt.Fprint(w, "[ x ]\t")
				} else {
					fmt.Fprint(w, "[   ]\t")
				}
			}
		}
		fmt.Fprintf(w, "\n")
	}
	fmt.Fprintf(w, " row \n")
	w.Flush()
}

//clears the terminal if it is a linux system
func clear() {
	if runtime.GOOS == "linux" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

//checks if a row or col is within the board
func inbound(row, col, height, width int) bool {
	if row < 0 || row >= height {
		return false
	}

	if col < 0 || col >= width {
		return false
	}

	return true
}
