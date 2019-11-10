package board

import (
	"errors"
	"fmt"
)

// TODO: add tests

//Config is a wrapper for creating a new board to make sure it is constrained
//to the configuration parameters.
type Config struct {
	maxHeight, maxWidth, maxSize int
	square                       bool
}

//CreateBlankConfig gives you a Config that doesn't have any max values or
//flags set with default NewBoard
func CreateBlankConfig() Config {
	return Config{}
}

//CreateFullConfig gives you a config with all max values and flags set by user
func CreateFullConfig(height, width, size int, square bool) Config {
	return Config{height, width, size, square}
}

//CreateHWConfig gives you a config with only height and width max values set
func CreateHWConfig(height, width int) Config {
	return Config{maxHeight: height, maxWidth: width}
}

//SetMaxWidth sets the max value for width in Config
func (c *Config) SetMaxWidth(width int) {
	c.maxWidth = width
}

//SetMaxHeight sets the max value for height in Config
func (c *Config) SetMaxHeight(height int) {
	c.maxHeight = height
}

//SetMaxSize sets the max value for size which is the total number of cells in Config
func (c *Config) SetMaxSize(size int) {
	c.maxSize = size
}

//SetSquare sets the flag value for whether a new board needs to be square in Config
func (c *Config) SetSquare(square bool) {
	c.square = square
}

//CreateBoard checks all constraints that are set and if they check out
//creates a new board
func (c Config) CreateBoard(width, height int) (*Board, error) {
	if height > c.maxHeight && c.maxHeight > 0 {
		return nil, fmt.Errorf("Height is greater than Max Height: %d", c.maxHeight)
	} else if width > c.maxWidth && c.maxWidth > 0 {
		return nil, fmt.Errorf("Width is greater than Max Width: %d", c.maxWidth)
	} else if height*width > c.maxSize && c.maxSize > 0 {
		return nil, fmt.Errorf("Size is greater than Max Size: %d", c.maxSize)
	} else if c.square && height != width {
		return nil, errors.New("Not a square board")
	}
	return NewBoard(width, height)
}

//MakeCreateBoard make a function to create the new board with configuration
//constraints alternative to using Config object
func MakeCreateBoard(maxHeight, maxWidth, maxSize int, flag bool) func(int, int) (*Board, error) {
	f := NewBoard
	if maxHeight > 0 {
		f = func(h int) func(int, int) (*Board, error) {
			return func(width, height int) (*Board, error) {
				if height > h {
					return nil, fmt.Errorf("Height is greater than Max Height: %d", h)
				}
				return f(width, height)
			}
		}(maxHeight)
	}

	if maxWidth > 0 {
		f = func(w int) func(int, int) (*Board, error) {
			return func(width, height int) (*Board, error) {
				if width > w {
					return nil, fmt.Errorf("Width is greater than Max Width: %d", w)
				}
				return f(width, height)
			}
		}(maxWidth)
	}
	if maxSize > 0 {
		f = func(s int) func(int, int) (*Board, error) {
			return func(width, height int) (*Board, error) {
				if height*width > s {
					return nil, fmt.Errorf("Size is greater than Max Size: %d", s)
				}
				return f(width, height)
			}
		}(maxSize)
	}

	if flag {
		f = func() func(int, int) (*Board, error) {
			return func(height, width int) (*Board, error) {
				if height != width {
					return nil, errors.New("Not a square board")
				}
				return f(width, height)
			}
		}()
	}
	return f
}
