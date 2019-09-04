package main

import (
	"fmt"
)

//board is a 2d slice to keep all the cells
type board [][]cell

//directions: up, down, left, right,
//up-left, up-right, down-left, down-right
var directions = [8][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1},
	{-1, 1}, {1, -1}, {1, 1}}

func main() {
	fmt.Println("Welcome to Minesweeper!")
	run := true
	for run {
		fmt.Print("Start (Type s), Help (Type h), Quit (type q): ")
		var choice string
		if _, err := fmt.Scan(&choice); err != nil {
			fmt.Println(err.Error())
			continue
		}
		switch choice {
		case "s":
		case "h":
			help()
			continue
		case "q":
			return
		default:
			continue
		}
		fmt.Println("What size board would you like?")

		var width int
		getSetupInput(&width, "Width (between 10 and 20): ")

		var height int
		getSetupInput(&height, "Height (between 10 and 20): ")

		b, numMines := newBoard(width, height)

		b.printBoard()
		numCells := height*width - numMines
		fmt.Println("Number of Mines:", numMines)
		fmt.Println("Number of Unshown Cells:", numCells)
		win := b.runGame(height, width, height*width-numMines, numMines)

		if win {
			fmt.Println("You Win!!!")
		} else {
			fmt.Println("You Lose")
		}
		fmt.Println("Play again?")

		for choice != "Y" && choice != "N" {
			fmt.Print("Y/N: ")
			fmt.Scan(&choice)
		}
		if choice == "N" {
			run = false
		}
	}

}

func help() {
	fmt.Println(`This is Minesweeper!
First, you choose a board size.
Both Height and Width of the board should be between 10 and 20 (inclusive)
Each move, you have the choice of either Choosing, Expanding or Marking a cell.
After making a choice, you are then asked the row and column of the cell that
you want to perform that action on. Choosing means you want to show a cell.
It does not work on marked or already shown cells. If a zero cell is chosen,
it gets expanded. If a mine cell is chosen, you lose the game. Expanding means
you pick a cell and every cell around it is chosen. All the surrounding cells
go through the choosing process. Marking a cell means the cell is considered a
mine and is given an 'x'. Each mine cell is marked with an 'm'. Each non-mine
cell has a number telling you the number of mines in the 8 surrounding cells.
The game ends when you either hit a mine and lose or when you mark all mines
and show all non-mine cells to win.`)
}

//getSetupInput makes a loop to get valid input.
func getSetupInput(num *int, inputString string) {
	for true {
		fmt.Print(inputString)
		if _, err := fmt.Scan(num); err != nil {
			fmt.Println("Problem Reading in Value. Input an integer between 10 and 20 (inclusive)")
			continue
		}
		if *num > 20 || *num < 10 {
			fmt.Println("Value needs to be between 10 and 20 Cells")
		} else {
			return
		}
	}

}

//runGame starts the minesweeper game
func (b board) runGame(height, width, numCells, numMines int) bool {
	end := false
	for !end {
		fmt.Println("What do you want to do?")
		fmt.Println("Choose (type c), Expand (type e), Mark (type m): ")
		var choice string
		var row int
		var col int
		fmt.Scan(&choice)
		switch choice {
		case "c", "e", "m":
			fmt.Println("Which Cell?")
			fmt.Print("Row: ")
			if _, err := fmt.Scan(&row); err != nil {
				fmt.Println(err.Error())
				continue
			}
			if !inbound(row, 0, height, width) {
				fmt.Println("Row out of Bounds")
				continue
			}
			fmt.Print("Col: ")
			if _, err := fmt.Scan(&col); err != nil {
				fmt.Println(err.Error())
				continue
			}
			if !inbound(row, col, height, width) {
				fmt.Println("Column out of Bounds")
				continue
			}
		default:
			b.printBoard()
			fmt.Println("Not 'c', 'e' or 'm'")
			continue
		}

		var cells int
		switch choice {
		case "c":
			end, cells = b.choose(row, col)
			numCells -= cells
		case "e":
			end, cells = b.expand(row, col)
			numCells -= cells
		case "m":
			marked := b.mark(row, col)
			numMines += marked
		default:
			continue
		}

		b.printBoard()
		fmt.Println("Number of Mines:", numMines)
		fmt.Println("Number of Unshown Cells:", numCells)
		if end {
			return false
		}

		if numCells == 0 && numMines == 0 {
			end = true
		}

	}
	return true
}

//chooses an unrevealed cell to reveal
//assumes row and col are inbound
//shows the cell you chose
//returns whether you hit a mine or not and the number of cells actually chosen
func (b board) choose(row, col int) (bool, int) {
	//don't choose marked cells
	if b[row][col].mark || b[row][col].show {
		return false, 0
	}

	b[row][col].show = true
	numCells := 1
	if b[row][col].value == 'm' {
		return true, 1
	} else if b[row][col].value == '\x00' {
		_, cells := b.expand(row, col)
		numCells += cells
	}

	return false, numCells

}

//if given a shown cell it chooses all cells around it.
//returns whether you hit a mine or not and the number of cells actually chosen
func (b board) expand(row, col int) (bool, int) {
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
			end, i := b.choose(rowCheck, colCheck)
			if end {
				return true, numCells
			}
			numCells += i
		}
	}

	return false, numCells

}

//marks a cell as having a mine
func (b board) mark(row, col int) int {
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
