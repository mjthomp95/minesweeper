package board

import (
	"testing"
)

func TestCreateBoardHW(t *testing.T) {
	//tests the shouldn't create error
	maxHeight, maxWidth := 20, 15
	config := CreateHWConfig(maxHeight, maxWidth)
	testCases := []struct {
		height, width int
	}{
		{1, 1},
		{1, 3},
		{5, 5},
		{10, 10},
	}

	for _, test := range testCases {
		b, e := config.CreateBoard(test.width, test.height)
		if e != nil {
			t.Errorf("There should not be an error with given Height: %d and Width: %d and config Max Height: %d Max Width %d",
				test.height, test.width, maxHeight, maxWidth)
		} else if b.height != test.height || b.width != test.width {
			t.Errorf("Board doesn't have correct dimensions. Expected: {Height: %d, Width: %d} Got: {Height: %d, Width: %d}",
				test.height, test.width, b.height, b.width)
		}
	}

	//these should give errors
	errorCases := []struct {
		height, width int
	}{
		{0, 1},
		{1, 0},
		{-5, 1},
		{1, -5},
		{21, 5},
		{5, 16},
		{21, 16},
		{100, 100},
	}

	for _, test := range errorCases {
		b, e := config.CreateBoard(test.width, test.height)
		if e == nil {
			t.Errorf("This should throw an error. Given Height: %d Width: %d  Got: %v",
				test.height, test.width, b)
		}
	}
}

func TestCreateBoardFull(t *testing.T) {
	//test should succeed with square
	maxHeight, maxWidth, maxSize, square := 20, 20, 150, true
	config := CreateFullConfig(maxHeight, maxWidth, maxSize, square)
	testCases1 := []struct {
		height, width int
	}{
		{1, 1},
		{2, 2},
		{5, 5},
		{10, 10},
		{12, 12},
	}

	for _, test := range testCases1 {
		b, e := config.CreateBoard(test.width, test.height)
		if e != nil {
			t.Errorf("There should not be an error with given Height: %d and Width: %d and config Max Height: %d Max Width %d",
				test.height, test.width, maxHeight, maxWidth)
		} else if b.height != test.height || b.width != test.width {
			t.Errorf("Board doesn't have correct dimensions. Expected: {Height: %d, Width: %d} Got: {Height: %d, Width: %d}",
				test.height, test.width, b.height, b.width)
		}
	}

	//error tests with square
	errorCases := []struct {
		height, width int
	}{
		{0, 10},
		{-5, 5},
		{5, -5},
		{10, 0},
		{1, 2}, //not square
		{2, 1},
		{10, 15},
		{13, 13},
	}

	for _, test := range errorCases {
		b, e := config.CreateBoard(test.width, test.height)
		if e == nil {
			t.Errorf("This should throw an error. Given Height: %d Width: %d  Got: %v",
				test.height, test.width, b)
		}
	}

	//success tests 2 without square
	config.SetSquare(false)
	testCases2 := []struct {
		height, width int
	}{
		{1, 2},
		{2, 1},
		{5, 6},
		{6, 5},
		{10, 11},
		{11, 5},
	}

	for _, test := range testCases2 {
		b, e := config.CreateBoard(test.width, test.height)
		if e != nil {
			t.Errorf("There should not be an error with given Height: %d and Width: %d and config Max Height: %d Max Width %d",
				test.height, test.width, maxHeight, maxWidth)
		} else if b.height != test.height || b.width != test.width {
			t.Errorf("Board doesn't have correct dimensions. Expected: {Height: %d, Width: %d} Got: {Height: %d, Width: %d}",
				test.height, test.width, b.height, b.width)
		}
	}
}

func TestMakeCreateBoard(t *testing.T) {
	//test should succeed with square
	maxHeight, maxWidth, maxSize, square := 20, 20, 150, true
	createBoard := MakeCreateBoard(maxHeight, maxWidth, maxSize, square)
	testCases1 := []struct {
		height, width int
	}{
		{1, 1},
		{2, 2},
		{5, 5},
		{10, 10},
		{12, 12},
	}

	for _, test := range testCases1 {
		b, e := createBoard(test.width, test.height)
		if e != nil {
			t.Errorf("There should not be an error with given Height: %d and Width: %d and config Max Height: %d Max Width %d",
				test.height, test.width, maxHeight, maxWidth)
		} else if b.height != test.height || b.width != test.width {
			t.Errorf("Board doesn't have correct dimensions. Expected: {Height: %d, Width: %d} Got: {Height: %d, Width: %d}",
				test.height, test.width, b.height, b.width)
		}
	}

	//error tests with square
	errorCases := []struct {
		height, width int
	}{
		{0, 10},
		{-5, 5},
		{5, -5},
		{10, 0},
		{1, 2}, //not square
		{2, 1},
		{10, 15},
		{13, 13},
	}

	for _, test := range errorCases {
		b, e := createBoard(test.width, test.height)
		if e == nil {
			t.Errorf("This should throw an error. Given Height: %d Width: %d  Got: %v",
				test.height, test.width, b)
		}
	}

	maxHeight, maxWidth, maxSize, square = 20, 20, 150, false
	createBoard = MakeCreateBoard(maxHeight, maxWidth, maxSize, square)
	//success tests 2 without square
	testCases2 := []struct {
		height, width int
	}{
		{1, 2},
		{2, 1},
		{5, 6},
		{6, 5},
		{10, 11},
		{11, 5},
	}

	for _, test := range testCases2 {
		b, e := createBoard(test.width, test.height)
		if e != nil {
			t.Errorf("There should not be an error with given Height: %d and Width: %d and config Max Height: %d Max Width %d",
				test.height, test.width, maxHeight, maxWidth)
		} else if b.height != test.height || b.width != test.width {
			t.Errorf("Board doesn't have correct dimensions. Expected: {Height: %d, Width: %d} Got: {Height: %d, Width: %d}",
				test.height, test.width, b.height, b.width)
		}
	}
}
