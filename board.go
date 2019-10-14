package main

import (
	"fmt"
	"math/rand"
	"os"
	"text/tabwriter"
	"time"
)

type cell struct {
	value rune
	show  bool
	mark  bool
}

//Board is a 2d slice to keep all the cells
type Board [][]cell

//only for when out of bounds
type outOfBoundsError struct {
	height, width int
}

func (e *outOfBoundsError) Error() string {
	return fmt.Sprintf("Out of bounds. Height: %d, Width: %d", e.height, e.width)
}


//directions: up, down, left, right,
//up-left, up-right, down-left, down-right
var directions = [8][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1},
	{-1, 1}, {1, -1}, {1, 1}}

//Increases the count of the cells around a mine
func (b Board) increaseCount(row, col int) {
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
func (b Board) placeBombs() int {
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

func blankBoard(width, height int) Board {
	b := make(Board, height)
	for i := range b {
		b[i] = make([]cell, width)
	}
	return b
}

//NewBoard returns a Board of specified size and the number of mines
func NewBoard(width, height int) (Board, int) {

	b := blankBoard(width, height)

	numMines := b.placeBombs()
	return b, numMines
}

//Choose reveals an unrevealed cell
//assumes row and col are inbound
//shows the cell you chose
//returns whether you hit a mine or not and the number of cells actually chosen
func (b Board) Choose(row, col int) (bool, int) {
	//don't choose marked cells
	if b[row][col].mark || b[row][col].show {
		return false, 0
	}

	b[row][col].show = true
	numCells := 1
	if b[row][col].value == 'm' {
		return true, 1
	} else if b[row][col].value == '\x00' {
		_, cells := b.Expand(row, col)
		numCells += cells
	}

	return false, numCells

}

//Expand if given a shown cell it chooses all cells around it.
//returns whether you hit a mine or not and the number of cells actually chosen
func (b Board) Expand(row, col int) (bool, int) {
	if !b[row][col].show {
		return false, 0
	}

	height := len(b)
	width := len(b[0])
	numCells := 0

	for _, direction := range directions {
		rowCheck := row + direction[0]
		colCheck := col + direction[1]
		if inbound(rowCheck, colCheck, height, width) {
			end, i := b.Choose(rowCheck, colCheck)
			if end {
				return true, numCells
			}
			numCells += i
		}
	}

	return false, numCells

}

//Mark denotes a cell as having a mine
func (b Board) Mark(row, col int) (int, error) {
	if !b.inbound(row, col) {
		return 0,
	}
	if b[row][col].mark {
		b[row][col].mark = false
		return 1
	} else if b[row][col].show {
		return 0
	} else {
		b[row][col].mark = true
		return -1
	}
}

//center text
func center(s rune, w int) string {
	if s == 'm' {
		return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*c", (w+1)/2, s))
	}

	return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*v", (w+1)/2, s))
}

//PrintBoard prints out the board to the terminal
func (b Board) PrintBoard() {
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

//checks if a row or col is within the board
func (b Board) inbound(row, col int) bool {
	if row < 0 || row >= b.height {
		return false
	}

	if col < 0 || col >= b.width {
		return false
	}

	return true
}
