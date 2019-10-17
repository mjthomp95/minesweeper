package main

import (
	"errors"
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
type Board struct {
	cells         [][]cell
	height, width int
}

//OutOfBoundsError tells if a certain row and column are not within a board.
type OutOfBoundsError struct {
	row, col, height, width int
}

func (e *OutOfBoundsError) Error() string {
	return fmt.Sprintf("out of bounds; Row: %d and Column: %d not within"+
		"Board with Height: %d, Width: %d", e.row, e.col, e.height, e.width)
}

//directions: up, down, left, right,
//up-left, up-right, down-left, down-right
var directions = [8][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1},
	{-1, 1}, {1, -1}, {1, 1}}

//Increases the count of the cells around a mine
func (b Board) increaseCount(row, col int) {

	for _, direction := range directions {
		rowCheck := row + direction[0]
		colCheck := col + direction[1]
		if b.Inbound(rowCheck, colCheck) {
			val := b.cells[rowCheck][colCheck].value
			if val != 'm' {
				b.cells[rowCheck][colCheck].value = rune(int(val) + 1)
			}
		}
	}
}

//Randomly places mines across the input board.
func (b Board) placeBombs() int {
	//maxMines is 10% of the number of cells
	maxMines := int(float32(b.height) * float32(b.width) * .1)
	numMines := 0
	rand.Seed(time.Now().UnixNano())
	for numMines < maxMines {
		//not likely to take more than a couple passes
		for i := range b.cells {
			for j, c := range b.cells[i] {
				if numMines >= maxMines {
					return numMines
				}
				if c.value == 'm' {
					continue
				}
				if rand.Intn(100) < 10 {
					b.cells[i][j].value = 'm'
					b.increaseCount(i, j)
					numMines++
				}
			}
		}
	}
	return numMines
}

//blankBoard helper function for NewBoard. Doesn't check validity of row and height
//creates a Board with given row and height, but a blank cells array.
func blankBoard(width, height int) Board {
	vals := make([][]cell, height)
	for i := range vals {
		vals[i] = make([]cell, width)
	}
	return Board{vals, height, width}
}

//NewBoard returns a Board of specified size, the number of mines or an error if
// the height or width isn't a positive non-zero value
func NewBoard(width, height int) (*Board, int, error) {

	if width <= 0 || height <= 0 {
		return nil, 0, errors.New("need positive, non-zero values for width and height")
	}
	b := blankBoard(width, height)

	numMines := b.placeBombs()
	return &b, numMines, nil
}

//Choose reveals an unrevealed cell
//assumes row and col are inbound
//shows the cell you chose
//returns whether you hit a mine or not and the number of cells actually chosen
// TODO: Possibly add wrapper functions for all actions such that the wrapped function doesn't check for errors but wrapper function does.
//prevents double checking for errors since some actions are called within other actions.
func (b Board) Choose(row, col int) (bool, int, error) {
	//don't choose marked cells
	if !b.Inbound(row, col) {
		return false, 0, &OutOfBoundsError{row, col, b.height, b.width}
	}

	if b.cells[row][col].mark || b.cells[row][col].show {
		return false, 0, nil
	}

	b.cells[row][col].show = true
	numCells := 1
	if b.cells[row][col].value == 'm' {
		return true, 1, nil
	} else if b.cells[row][col].value == '\x00' {
		_, cells, _ := b.Expand(row, col) //inbounds already checked
		numCells += cells
	}

	return false, numCells, nil

}

//Expand if given a shown cell it chooses all cells around it.
//returns whether you hit a mine or not and the number of cells actually chosen
func (b Board) Expand(row, col int) (bool, int, error) {
	if !b.Inbound(row, col) {
		return false, 0, &OutOfBoundsError{row, col, b.height, b.width}
	}

	if !b.cells[row][col].show {
		return false, 0, nil
	}

	numCells := 0

	for _, direction := range directions {
		rowCheck := row + direction[0]
		colCheck := col + direction[1]
		end, i, err := b.Choose(rowCheck, colCheck) //inbounds gets checked
		if err != nil {
			continue
		}
		if end {
			return true, numCells, nil
		}
		numCells += i
	}

	return false, numCells, nil

}

//Mark denotes a cell as having a mine
func (b Board) Mark(row, col int) (int, error) {
	if !b.Inbound(row, col) {
		return 0, &OutOfBoundsError{row, col, b.height, b.width}
	}

	if b.cells[row][col].mark {
		b.cells[row][col].mark = false
		return 1, nil
	} else if b.cells[row][col].show {
		return 0, nil
	} else {
		b.cells[row][col].mark = true
		return -1, nil
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
//// TODO: Change this to returning a string
func (b Board) PrintBoard() {
	clear()
	w := tabwriter.NewWriter(os.Stdout, 5, 1, 0, ' ', 0)
	fmt.Fprint(w, "\t") //spacing
	if b.height <= 0 {
		return // TODO: quick fix make more permanent fix
	}
	for i := range b.cells[0] {
		fmt.Fprint(w, center(rune(i), 5), "\t")
	}
	fmt.Fprint(w, "col\n")
	for i := range b.cells {
		fmt.Fprint(w, center(rune(i), 5), "\t")
		for j := range b.cells[0] {
			if b.cells[i][j].show {
				fmt.Fprint(w, "[", center(b.cells[i][j].value, 3), "]\t")
			} else {
				if b.cells[i][j].mark {
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

//Inbound checks if a row or col is within the board
func (b Board) Inbound(row, col int) bool {
	if row < 0 || row >= b.height {
		return false
	}

	if col < 0 || col >= b.width {
		return false
	}

	return true
}
