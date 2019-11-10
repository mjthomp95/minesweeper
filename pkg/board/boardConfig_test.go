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

// TODO: add more tests
