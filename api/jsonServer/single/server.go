package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	mBoard "github.com/mjthomp95/minesweeper/pkg/board"
)

// IDEA: Create two versions of api where returning full board or just changes
type newBoard struct {
	Cells [][]rune `json:"cells"`
	Mines int      `json:"mines"`
}

type result struct {
	Cells    []mBoard.Output `json:"cells"`
	Mines    int             `json:"mines"`
	NumCells int             `json:"numcells"`
	Err      string          `json:"error"` //Send "Win" or "Lose" in Err when game over
}

// IDEA: use a map of ids to different boards to have multiple instances of a game.
// IDEA: when making map, have lock for map and board. create closure to access board.
func main() {
	//one instance of game
	var board *mBoard.Board
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		log.Println()
		if tmp := newHandler(w, r, board); tmp != nil {
			board = tmp
		}
	})
	http.HandleFunc("/end", func(w http.ResponseWriter, r *http.Request) {
		board = nil
		w.Header().Set("Access-Control-Allow-Origin", "*")
	})

	choose := moveWrapper(func(row, col int, b **mBoard.Board) (output []mBoard.Output, err error) {
		return (*b).Choose(row, col)
	}, &board)

	expand := moveWrapper(func(row, col int, b **mBoard.Board) (output []mBoard.Output, err error) {
		return (*b).Expand(row, col)
	}, &board)

	mark := moveWrapper(func(row, col int, b **mBoard.Board) (output []mBoard.Output, err error) {
		return (*b).Mark(row, col)
	}, &board)

	moveHandler := func(move func(*http.Request, *result) error) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var res result
			if board == nil {
				res.Err = "No game running"
			} else {
				move(r, &res)
				if board != nil {
					log.Println("\n", board.PrintBoard())
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			json.NewEncoder(w).Encode(res)
		}
	}

	http.HandleFunc("/choose", moveHandler(choose))

	http.HandleFunc("/expand", moveHandler(expand))

	http.HandleFunc("/mark", moveHandler(mark))

	http.HandleFunc("/output", func(w http.ResponseWriter, r *http.Request) {
		return
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	//give web app
	return
}

// TODO: Change all to json
func newHandler(w http.ResponseWriter, r *http.Request, board *mBoard.Board) *mBoard.Board {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if board != nil {
		fmt.Fprintf(w, "Game already in session")
		return nil
	}

	if r.Method != "POST" {
		fmt.Fprintf(w, "Only POST is supported")
		return nil
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprint(w, "Trouble parsing your POST form")
		return nil
	}

	heightString := r.FormValue("height")
	widthString := r.FormValue("width")

	if heightString != "" && widthString != "" {
		height, err := strconv.Atoi(heightString)
		if err != nil {
			fmt.Fprint(w, "height not a number")
			return nil
		}

		width, err := strconv.Atoi(widthString)
		if err != nil {
			fmt.Fprint(w, "width not a number")
			return nil
		}

		tmpBoard, err := mBoard.NewBoard(height, width)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return nil
		}
		resp := newBoard{Cells: tmpBoard.OutputBoard(), Mines: tmpBoard.GetMines()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return tmpBoard
	}
	fmt.Fprint(w, "Trouble getting your height and width values")
	return nil
}

func rowColHelper(r *http.Request) (row, col int, err error) {
	if r.Method != "POST" {
		err = errors.New("Only POST is supported")
		return
	}
	if err = r.ParseForm(); err != nil {
		return
	}
	rowString := r.FormValue("row")
	colString := r.FormValue("col")
	if rowString != "" && colString != "" {
		row, err = strconv.Atoi(rowString)
		col, err = strconv.Atoi(colString)
	}
	return
}

func moveWrapper(move func(int, int, **mBoard.Board) ([]mBoard.Output, error), board **mBoard.Board) func(*http.Request, *result) error {
	// TODO: change to take board and a function that takes board to make a move if this doesn't work
	return func(r *http.Request, res *result) (err error) {
		row, col, rcErr := rowColHelper(r)
		if rcErr != nil {
			err = rcErr
			return
		}
		res.Cells, err = move(row, col, board)
		log.Println((*board).PrintBoard()) // TEMP: used for when no frontend
		res.Mines = (*board).GetMines()
		res.NumCells = (*board).GetNumCells()
		switch err.(type) {
		case nil:
		case mBoard.GameOver:
			(*board) = nil
			if err.(mBoard.GameOver).Win {
				res.Err = "Win"
			} else {
				res.Err = "Lose"
			}
		}
		return
	}
}
