package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	mBoard "github.com/mjthomp95/minesweeper/pkg/board"
)

// TODO: Change name
type newBoard struct {
	Cells [][]rune `json:"cells"`
	Mines int      `json:"mines"`
}

// TODO: change all responses from 'mark', 'choose', 'expand', 'new', and 'output' to json.
// IDEA: use a map of ids to different boards to have multiple instances of a game.
func main() {
	//one instance of game
	var board *mBoard.Board
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		if board != nil {
			fmt.Fprintf(w, "Game already in session")
			return
		}

		if r.Method != "POST" {
			fmt.Fprintf(w, "Only POST is supported")
			return
		}

		if err := r.ParseForm(); err != nil {
			fmt.Fprint(w, "Trouble parsing your POST form")
			return
		}

		heightString := r.FormValue("height")
		widthString := r.FormValue("width")

		if heightString != "" && widthString != "" {
			height, err := strconv.Atoi(heightString)
			if err != nil {
				fmt.Fprint(w, "height not a number")
				return
			}

			width, err := strconv.Atoi(widthString)
			if err != nil {
				fmt.Fprint(w, "width not a number")
				return
			}

			tmpBoard, err := mBoard.NewBoard(height, width)
			if err != nil {
				fmt.Fprint(w, err.Error())
				return
			}
			board = tmpBoard
			log.Println(board.PrintBoard())
			log.Println(board.OutputBoard())
			resp := newBoard{Cells: board.OutputBoard(), Mines: board.GetMines()}
			data, err := json.Marshal(resp) // TODO: Change this up
			if err != nil {
				fmt.Fprint(w, "couldn't marshal board")
				return
			}
			log.Println(data)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		} else {
			fmt.Fprint(w, "Trouble getting your height and width values")
		}
	})

	http.HandleFunc("/choose", func(w http.ResponseWriter, r *http.Request) { // IDEA: Abstract closure also
		gameOver, win := chooseHandler(w, r, board)
		if gameOver {
			log.Println(board.PrintBoard())
			board = nil
			if win {
				fmt.Fprint(w, "Win")
			} else {
				fmt.Fprint(w, "Lose")
			}
		}
	})
	http.HandleFunc("/expand", func(w http.ResponseWriter, r *http.Request) {
		gameOver, win := expandHandler(w, r, board)
		if gameOver {
			log.Println("\n", board.PrintBoard())
			board = nil
			if win {
				fmt.Fprint(w, "Win")
			} else {
				fmt.Fprint(w, "Lose")
			}
		}
	})
	http.HandleFunc("/mark", func(w http.ResponseWriter, r *http.Request) {
		gameOver, win := markHandler(w, r, board)
		if gameOver {
			log.Println("\n", board.PrintBoard())
			board = nil
			if win {
				fmt.Fprint(w, "Win")
			} else {
				fmt.Fprint(w, "Lose")
			}
		}
	})
	http.HandleFunc("/output", func(arg1 http.ResponseWriter, arg2 *http.Request) {
		return
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	//give web app
	return
}

// IDEA: very similar handlers can I abstract and make a common handler that takes a board function to run and handle?
func chooseHandler(w http.ResponseWriter, r *http.Request, board *mBoard.Board) (gameOver, win bool) {
	if r.Method != "POST" {
		fmt.Fprint(w, "Only POST is supported")
		return
	}
	if board != nil {
		if err := r.ParseForm(); err != nil {
			fmt.Fprint(w, "Trouble parsing your POST form")
			return
		}
		rowString := r.FormValue("row")
		colString := r.FormValue("col")
		if rowString != "" && colString != "" {
			row, err := strconv.Atoi(rowString)
			if err != nil {
				//bad value
				fmt.Fprint(w, "row not a number")
				return
			}
			col, err := strconv.Atoi(colString)
			if err != nil {
				fmt.Fprint(w, "col not a number")
				return
			}
			_, err = board.Choose(row, col) // TODO: Do something with output
			switch err.(type) {
			case nil:
			case mBoard.GameOver:
				gameOver = true
				win = err.(mBoard.GameOver).Win
				return
			default:
				fmt.Fprint(w, err.Error())
			}
			log.Println("\n", board.PrintBoard())
			resp := newBoard{board.OutputBoard(), board.GetMines()}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		fmt.Fprint(w, "Trouble getting row and col values")
		return
	}

	fmt.Fprint(w, "No game running")
	return
}

func expandHandler(w http.ResponseWriter, r *http.Request, board *mBoard.Board) (gameOver, win bool) {
	log.Println("Expand Start")
	if r.Method != "POST" {
		fmt.Fprint(w, "Only POST is supported")
		return
	}
	if board != nil {
		if err := r.ParseForm(); err != nil {
			fmt.Fprint(w, "Trouble parsing your POST form")
			return
		}
		rowString := r.FormValue("row")
		colString := r.FormValue("col")
		if rowString != "" && colString != "" {
			row, err := strconv.Atoi(rowString)
			if err != nil {
				//bad value
				fmt.Fprint(w, "row not a number")
				return
			}
			col, err := strconv.Atoi(colString)
			if err != nil {
				fmt.Fprint(w, "col not a number")
				return
			}
			_, err = board.Expand(row, col) // TODO: Do something with output
			switch err.(type) {
			case nil:
			case mBoard.GameOver:
				gameOver = true
				win = err.(mBoard.GameOver).Win
				return
			default:
				fmt.Fprint(w, err.Error())
			}
			log.Println("\n", board.PrintBoard())
			log.Println("Expand")
			resp := newBoard{board.OutputBoard(), board.GetMines()}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		fmt.Fprint(w, "Trouble getting row and col values")
		return
	}

	fmt.Fprint(w, "No game running")
	return
}

func markHandler(w http.ResponseWriter, r *http.Request, board *mBoard.Board) (gameOver, win bool) {
	if r.Method != "POST" {
		fmt.Fprint(w, "Only POST is supported")
		return
	}
	if board != nil {
		if err := r.ParseForm(); err != nil {
			fmt.Fprint(w, "Trouble parsing your POST form")
			return
		}
		rowString := r.FormValue("row")
		colString := r.FormValue("col")
		if rowString != "" && colString != "" {
			row, err := strconv.Atoi(rowString)
			if err != nil {
				//bad value
				fmt.Fprint(w, "row not a number")
				return
			}
			col, err := strconv.Atoi(colString)
			if err != nil {
				fmt.Fprint(w, "col not a number")
				return
			}
			_, err = board.Mark(row, col) // TODO: Do something with output
			switch err.(type) {
			case nil:
			case mBoard.GameOver:
				gameOver = true
				win = err.(mBoard.GameOver).Win
				return
			default:
				fmt.Fprint(w, err.Error())
			}
			log.Println("\n", board.PrintBoard())
			resp := newBoard{board.OutputBoard(), board.GetMines()}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		fmt.Fprint(w, "Trouble getting row and col values")
		return
	}

	fmt.Fprint(w, "No game running")
	return
}
