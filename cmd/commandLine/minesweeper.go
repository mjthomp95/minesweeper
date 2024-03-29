package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/mjthomp95/minesweeper/pkg/board"
)

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

		b, err := board.NewBoard(width, height)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		clear()
		fmt.Println(b.PrintBoard())
		fmt.Println("Number of Mines:", b.GetMines())
		fmt.Println("Number of Unshown Cells:", b.GetNumCells())
		win := runGame(*b, height, width)
		clear()
		fmt.Println(b.PrintBoard())
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

//clears the terminal if it is a linux system
func clear() {
	if runtime.GOOS == "linux" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

//runGame starts the minesweeper game
func runGame(b board.Board, height, width int) bool {
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
			if !b.Inbound(row, 0) {
				fmt.Println("Row out of Bounds")
				continue
			}
			fmt.Print("Col: ")
			if _, err := fmt.Scan(&col); err != nil {
				fmt.Println(err.Error())
				continue
			}
			if !b.Inbound(row, col) {
				fmt.Println("Column out of Bounds")
				continue
			}
		default:
			clear()
			fmt.Println(b.PrintBoard())
			fmt.Println("Not 'c', 'e' or 'm'")
			continue
		}

		switch choice {
		case "c":
			_, err := b.Choose(row, col)
			switch err.(type) {
			case nil:
			case board.GameOver:
				return err.(board.GameOver).Win
			default:
				fmt.Println(err.Error())
				continue
			}
		case "e":
			_, err := b.Expand(row, col) //don't need the returned list
			switch err.(type) {
			case nil:
			case board.GameOver:
				return err.(board.GameOver).Win
			default:
				fmt.Println(err.Error())
				continue
			}
		case "m":
			_, err := b.Mark(row, col)
			switch err.(type) {
			case nil:
			case board.GameOver:
				return err.(board.GameOver).Win
			default:
				fmt.Println(err.Error())
				continue
			}
		default:
			continue
		}

		clear()
		fmt.Println(b.PrintBoard())
		fmt.Println("Number of Mines:", b.GetMines())
		fmt.Println("Number of Cells left:", b.GetNumCells())
	}
	return true
}
