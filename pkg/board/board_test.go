package board

import (
	"reflect"
	"testing"
)

func TestBlankBoard(t *testing.T) {
	var blankCell cell
	testCases := []struct {
		height int
		width  int
		result Board
	}{
		{1, 1, Board{[][]cell{{blankCell}}, 1, 1}},
		{2, 2, Board{[][]cell{{blankCell, blankCell}, {blankCell, blankCell}}, 2, 2}},
		{1, 3, Board{[][]cell{{blankCell, blankCell, blankCell}}, 1, 3}},
	}

	cellCheck := blankBoard(1, 1).cells
	if !reflect.DeepEqual(cellCheck[0][0], blankCell) {
		t.Fatal("Should be Equal")
	}

	for _, test := range testCases {
		b := blankBoard(test.width, test.height)
		if !reflect.DeepEqual(b, test.result) {
			t.Error("Should be Equal", b, test.result)
		}
	}
}

func countBombs(cells [][]cell) (count int) {
	for _, cellRow := range cells {
		for _, cell := range cellRow {
			if cell.value == 'm' {
				count++
			}
		}
	}
	return
}

func TestPlaceBombs(t *testing.T) {
	testCases := []struct {
		height int
		width  int
		result int
	}{
		{10, 1, 1},
		{20, 1, 2},
		{10, 2, 2},
		{1, 10, 1},
		{1, 20, 2},
		{5, 5, 2},
		{2, 10, 2},
		{10, 10, 10},
	}

	for _, test := range testCases {
		b := blankBoard(test.width, test.height)
		numMines := b.placeBombs()
		if numMines != countBombs(b.cells) {
			t.Error("numMines not equal to the counted bombs")
		} else if numMines != test.result {
			t.Errorf("Not the correct number of bombs placed. Got: %d Expected: %d",
				numMines, test.result)
		}
	}
}

func TestNewBoard(t *testing.T) {
	testCases := []struct {
		height int
		width  int
		result int //number of mines
	}{
		{10, 1, 1},
		{20, 1, 2},
		{10, 2, 2},
		{1, 10, 1},
		{1, 20, 2},
		{5, 5, 2},
		{2, 10, 2},
		{10, 10, 10},
	}

	for _, test := range testCases {
		b, numMines, err := NewBoard(test.width, test.height)
		if err == nil {
			if counted := countBombs(b.cells); numMines != counted {
				t.Errorf("Returned number of mines doesn't equal counted."+
					"NumMines: %d Counted: %d", numMines, counted)
			} else if numMines != test.result {
				t.Errorf("Number of mines doesn't equal expected. Expected: %d Got: %d",
					test.result, numMines)
			}
		} else {
			t.Error("Should not be an error", err.Error())
		}
	}

	errorCases := []struct {
		height int
		width  int
	}{
		{0, 0},
		{-1, 2},
		{-10, -5},
		{5, -5},
	}

	for _, error := range errorCases {
		b, numMines, err := NewBoard(error.width, error.height)
		if err == nil {
			t.Errorf("Error should not be nil. Board: %v NumMines: %d", b, numMines)
		} else if numMines != 0 {
			t.Errorf("NumMines should be zero. Board: %v", *b)
		} else if b != nil {
			t.Errorf("Board should be nil. Board: %v", *b)
		}
	}
}

// TODO: More Tests
