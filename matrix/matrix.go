package matrix

import (
	"errors"
	"fmt"
)

// ScanDirection scan matrix direction
type ScanDirection uint

const (
	// ROW for row first
	ROW ScanDirection = 1

	// COLUMN for column first
	COLUMN ScanDirection = 2
)

// State value of matrix map[][]
type State uint16

const (
	// StateInit represents the initial block state of the matrix
	StateInit State = iota
	// ZERO represents the initial state.
	// Deprecated: use StateInit instead
	ZERO = StateInit

	// StateFalse represents the block has been set to false
	StateFalse State = 0x1

	// StateTrue represents the block has been set to TRUE
	StateTrue State = 0x2

	// StateVersion indicates the version block of matrix
	StateVersion State = 0x3

	// StateFormat indicates the format block of matrix
	StateFormat State = 0x4

	// StateFinder indicates the finder block of matrix
	StateFinder State = 0x5
)

func (s State) String() string {
	return fmt.Sprintf("0x%X", uint16(s))
}

var (
	// ErrorOutRangeOfW x out of range of Width
	ErrorOutRangeOfW = errors.New("out of range of width")

	// ErrorOutRangeOfH y out of range of Height
	ErrorOutRangeOfH = errors.New("out of range of height")
)

// StateSliceMatched should be
// Deprecated: since rule3_backup removed
func StateSliceMatched(ss1, ss2 []State) bool {
	if len(ss1) != len(ss2) {
		return false
	}
	for idx := range ss1 {
		if (ss1[idx] ^ ss2[idx]) != 0 {
			return false
		}
	}

	return true
}

// New generate a matrix with map[][]bool
func New(width, height int) *Matrix {
	mat := make([][]State, width)
	for w := 0; w < width; w++ {
		mat[w] = make([]State, height)
	}

	m := &Matrix{
		mat:    mat,
		width:  width,
		height: height,
	}

	m.init()
	return m
}

// Matrix is a matrix data type
// width:3 height: 4 for [3][4]int
type Matrix struct {
	mat    [][]State
	width  int
	height int
}

// do some init work
func (m *Matrix) init() {
	for w := 0; w < m.width; w++ {
		for h := 0; h < m.height; h++ {
			m.mat[w][h] = StateInit
		}
	}
}

// Print to stdout
func (m *Matrix) print() {
	m.Iterate(ROW, func(x, y int, s State) {
		fmt.Printf("%2d ", s)
		if (x + 1) == m.width {
			fmt.Println()
		}
	})
}

func (m *Matrix) Print() {
	m.print()
}

// Copy matrix into a new Matrix
func (m *Matrix) Copy() *Matrix {
	mat2 := make([][]State, m.width)
	for w := 0; w < m.width; w++ {
		mat2[w] = make([]State, m.height)
		copy(mat2[w], m.mat[w])
	}

	m2 := &Matrix{
		width:  m.width,
		height: m.height,
		mat:    mat2,
	}

	return m2
}

// Width ... width
func (m *Matrix) Width() int {
	return m.width
}

// Height ... height
func (m *Matrix) Height() int {
	return m.height
}

// Set [w][h] as true
func (m *Matrix) Set(w, h int, c State) error {
	if w >= m.width || w < 0 {
		return ErrorOutRangeOfW
	}
	if h >= m.height || h < 0 {
		return ErrorOutRangeOfH
	}
	m.mat[w][h] = c
	return nil
}

// Get state value from matrix with position {x, y}
func (m *Matrix) Get(w, h int) (State, error) {
	if w >= m.width || w < 0 {
		return StateInit, ErrorOutRangeOfW
	}
	if h >= m.height || h < 0 {
		return StateInit, ErrorOutRangeOfH
	}
	return m.mat[w][h], nil
}

// IterateFunc ...
type IterateFunc func(int, int, State)

// Iterate the Matrix with loop direction ROW major or COLUMN major.
// COLUMN is recommended.
func (m *Matrix) Iterate(dir ScanDirection, f IterateFunc) {
	// row direction first
	if dir == ROW {
		for h := 0; h < m.height; h++ {
			for w := 0; w < m.width; w++ {
				f(w, h, m.mat[w][h])
			}
		}
		return
	}

	// column direction first
	if dir == COLUMN {
		for w := 0; w < m.width; w++ {
			for h := 0; h < m.height; h++ {
				f(w, h, m.mat[w][h])
			}
		}
		return
	}
}

// XOR ...
func XOR(s1, s2 State) State {
	if s1 != s2 {
		return StateTrue
	}
	return StateFalse
}

// Row return a row of matrix, cur should be y dimension.
func (m *Matrix) Row(cur int) []State {
	if cur >= m.height || cur < 0 {
		return nil
	}

	col := make([]State, m.height)
	for w := 0; w < m.width; w++ {
		col[w] = m.mat[w][cur]
	}
	return col
}

// Col return a slice of column, cur should be x dimension.
func (m *Matrix) Col(cur int) []State {
	if cur >= m.width || cur < 0 {
		return nil
	}

	return m.mat[cur]
}
