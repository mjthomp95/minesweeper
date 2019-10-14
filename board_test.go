package main

import (
	"reflect"
	"testing"
)

func TestBlankBoard(t *testing.T) {
	var blankCell cell
	testCases := []struct {
		width  int
		height int
		result board
	}{
		{1, 1, board{{blankCell}}},
		{2, 2, board{{blankCell, blankCell}, {blankCell, blankCell}}},
		{1, 3, board{{blankCell}, {blankCell}, {blankCell}}},
	}

	cellCheck := blankBoard(1, 1)
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

//
// func TestNewBoard(t *testing.T) {
//
// }
