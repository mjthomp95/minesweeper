package board

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"text/tabwriter"
	"time"
)

//Move is a function that does an action on the board
type Move func(int, int) ([]Output, error)

type cell struct {
	value rune
	show  bool
	mark  bool
}

//OutputCell is the type of each element of the output board given.
type OutputCell = rune

//Output is the type to tell that a cell has changed (shown or marked/unmarked)
type Output struct {
	Row, Col int
	Value    rune
}

//Board is a 2d slice to keep all the cells
type Board struct {
	cells         [][]cell
	outputCells   [][]OutputCell
	height, width int
	mines         int
	unshownCells  int
	gameOver      bool
	win           bool
}

//OutOfBoundsError tells if a certain row and column are not within a board.
type OutOfBoundsError struct {
	row, col, height, width int
}

func (e OutOfBoundsError) Error() string {
	return fmt.Sprintf("out of bounds; Row: %d and Column: %d not within"+
		"Board with Height: %d, Width: %d", e.row, e.col, e.height, e.width)
}

//GameOver tells you if you try to make a move and the game is over
type GameOver struct {
	Win bool
}

func (e GameOver) Error() string {
	if e.Win {
		return "You Won!"
	}
	return "You Lost"
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
func (b *Board) placeBombs() int {
	//maxMines is 10% of the number of cells
	maxMines := int(float32(b.height) * float32(b.width) * .1)
	numMines := 0
	rand.Seed(time.Now().UnixNano())
	for numMines < maxMines {
		//not likely to take more than a couple passes
		for i := range b.cells {
			for j, c := range b.cells[i] {
				if numMines >= maxMines {
					b.mines = numMines
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
	b.mines = numMines
	return numMines
}

//blankBoard helper function for NewBoard. Doesn't check validity of row and height
//creates a Board with given row and height, but a blank cells array.
func blankBoard(width, height int) *Board {
	vals := make([][]cell, height)
	outputVals := make([][]OutputCell, height)
	for i := range vals {
		vals[i] = make([]cell, width)
		outputVals[i] = make([]OutputCell, width)
	}
	return &Board{vals, outputVals, height, width, 0, width * height, false, false} //int overflow?
}

//NewBoard returns a Board of specified size, the number of mines or an error if
// the height or width isn't a positive non-zero value
func NewBoard(width, height int) (*Board, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("need positive, non-zero values for width and height")
	}
	b := blankBoard(width, height)
	b.placeBombs()
	return b, nil
}

//GetMines returns the number of mines subracted from the number of marked cells.
func (b Board) GetMines() int {
	return b.mines
}

//GetNumCells returns the number of cells that aren't shown and unmarked.
func (b Board) GetNumCells() int {
	return b.unshownCells
}

//Win returns whether you have won or not.
func (b Board) Win() bool {
	return b.win
}

//IsGameOver returns whether the game is done or not.
func (b Board) IsGameOver() bool {
	return b.gameOver
}

//checkGameOver checks to see if the game is over.
//should only check for unshownCells
func (b Board) checkWin() bool {
	if b.unshownCells == 0 && b.mines == 0 {
		return true
	}

	return false
}

//Choose reveals an unrevealed cell
//assumes row and col are inbound
//shows the cell you chose
//returns whether you hit a mine or not and the number of cells actually chosen
// IDEA: Possibly add wrapper functions for all actions such that the wrapped function doesn't check for errors but wrapper function does.
//prevents double checking for errors since some actions are called within other actions.
func (b *Board) Choose(row, col int) (output []Output, err error) {
	if b.gameOver {
		err = GameOver{b.win}
		return
	}
	if !b.Inbound(row, col) {
		err = OutOfBoundsError{row, col, b.height, b.width}
		return
	}

	//don't choose marked cells
	if b.cells[row][col].mark || b.cells[row][col].show {
		return
	}

	b.cells[row][col].show = true
	b.outputCells[row][col] = b.cells[row][col].value + 48 //want unicode for the integer
	output = append(output, Output{row, col, b.outputCells[row][col]})
	if b.cells[row][col].value == 'm' {
		b.gameOver = true
		err = GameOver{false}
		return
	} else if b.cells[row][col].value == '\x00' {
		expOut, expErr := b.Expand(row, col) //inbounds already checked
		output = append(output, expOut...)
		switch expErr.(type) {
		case GameOver:
			err = expErr
		}
	}

	b.unshownCells--
	if b.checkWin() {
		b.win = true
		b.gameOver = true
		err = GameOver{true}
	}
	return
}

//Expand if given a shown cell it chooses all cells around it.
//returns the list of cells you hit and what their value was.
func (b *Board) Expand(row, col int) (output []Output, err error) {
	if b.gameOver {
		err = GameOver{b.win}
		return
	}
	if !b.Inbound(row, col) {
		err = OutOfBoundsError{row, col, b.height, b.width}
		return
	}

	if !b.cells[row][col].show {
		return
	}

	for _, direction := range directions {
		rowCheck := row + direction[0]
		colCheck := col + direction[1]
		chOut, chErr := b.Choose(rowCheck, colCheck) //inbounds gets checked
		output = append(output, chOut...)
		switch chErr.(type) {
		case GameOver:
			err = chErr
			return
		}
	}

	return

}

//Mark denotes a cell as having a mine
func (b *Board) Mark(row, col int) ([]Output, error) {
	if b.gameOver {
		return nil, GameOver{b.win}
	}
	if !b.Inbound(row, col) {
		return nil, OutOfBoundsError{row, col, b.height, b.width}
	}

	output := make([]Output, 1, 1)
	if b.cells[row][col].mark {
		b.cells[row][col].mark = false
		b.mines++
		b.unshownCells++
		b.outputCells[row][col] = 0
		output[0] = Output{row, col, 0}
		return output, nil
	} else if b.cells[row][col].show {
		output[0] = Output{row, col, b.outputCells[row][col]}
		return output, nil
	} else {
		b.cells[row][col].mark = true
		b.mines--
		b.unshownCells--
		b.outputCells[row][col] = 'm'
		output[0] = Output{row, col, 'm'}
		if b.checkWin() {
			b.win = true
			b.gameOver = true
			return output, GameOver{true}
		}
		return output, nil
	}
}

//center text
func center(s rune, w int) string {
	if s == 'm' {
		return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*c", (w+1)/2, s))
	}

	return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*v", (w+1)/2, s))
}

//PrintBoard creates a string representation of the board to be printed out
// TODO: Change to use outputCells
func (b Board) PrintBoard() string {
	if b.height == 0 || b.width == 0 {
		return ""
	}
	var strBuilder strings.Builder
	w := tabwriter.NewWriter(&strBuilder, 7, 1, 0, ' ', 0)
	fmt.Fprint(w, "\t") //spacing
	for i := range b.cells[0] {
		fmt.Fprint(w, "  ", center(rune(i), 5), "\t")
	}
	fmt.Fprint(w, "col\n\n")
	for i := range b.cells {
		fmt.Fprint(w, "  ", center(rune(i), 5), "\t")
		for j := range b.cells[0] {
			if b.cells[i][j].show {
				fmt.Fprint(w, "  [", center(b.cells[i][j].value, 3), "]\t")
			} else {
				if b.cells[i][j].mark {
					fmt.Fprint(w, "  [ x ]\t")
				} else {
					fmt.Fprint(w, "  [   ]\t")
				}
			}
		}
		fmt.Fprintf(w, "\n\n")
	}
	fmt.Fprintf(w, " row \n")
	w.Flush()
	return strBuilder.String()
}

// IDEA: Add reset function to reuse the same board

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

//OutputBoard returns a representation of the whole board only showing cells that
//are marked or shown. Any other cell is blank.
func (b Board) OutputBoard() [][]OutputCell {
	return b.outputCells
}
