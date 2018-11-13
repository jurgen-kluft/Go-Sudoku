package main

import (
	"fmt"
	"testing"
)

var emptyPuzzle = SudokuPuzzle3x3{
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
}

func TestSudokuBoard_prepare(t *testing.T) {
	type fields struct {
		width  uint
		height uint
		puzzle [9][9]uint8
		board  [9][9][9]uint8
	}
	type args struct {
		puzzle SudokuPuzzle3x3
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"empty puzzle", fields{width: 9, height: 9}, args{emptyPuzzle}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SudokuBoard{
				width:  tt.fields.width,
				height: tt.fields.height,
				puzzle: tt.fields.puzzle,
				board:  tt.fields.board,
			}
			s.prepare(tt.args.puzzle)
			if s.width != 9 || s.height != 9 {
				t.Error("Initialized sudoku board is wrong!")
			}
		})
	}
}

func TestSudokuBoard_fillNumber(t *testing.T) {
	s := &SudokuBoard{}
	s.prepare(emptyPuzzle)

	n := uint(0)
	for x := uint(0); x < 3; x++ {
		for y := uint(0); y < 3; y++ {
			s.fillInNumber(x, y, n+1)
			if s.board[y][x][n] != 0 {
				t.Errorf("Filling in number %d at '(%d,%d)' gives wrong result as %d!\n", n, x, y, s.board[y][x][n])
			}
			if s.board[0][0][n] != 0 {
				t.Errorf("Filling in number %d at '(%d,%d)' gives wrong result at (%d,%d)=%d!\n", n, x, y, 0, 0, s.board[y][x][n])
			}
			n++
		}
	}
}

func TestSudokuBoard_sectorIterator2(t *testing.T) {
	for y := uint(0); y < 9; y++ {
		for x := uint(0); x < 9; x++ {
			p := point{x, y}
			iter := p.iterateOverAllCellsInSector()
			c := uint(0)
			for p, ok := iter(); ok; p, ok = iter() {
				if p.x != (x/3)*3+(c%3) {
					t.Errorf("Sector cell (%d,%d) iterator has wrong x, has %d and should be %d\n", x/3, y/3, p.x, (x/3)*3+(c%3))
				}
				if p.y != (y/3)*3+(c/3) {
					t.Errorf("Sector cell (%d,%d) iterator has wrong y, has %d and should be %d\n", x/3, y/3, p.y, (y/3)*3+(c/3))
				}
				c++
			}
			if c != 9 {
				t.Errorf("Sector cell iterator runs more than 9 times\n")
			}
		}
	}
}

func TestSudokuBoard_sectorIteratorAligned(t *testing.T) {
	p := point{0, 0}
	iter := p.iterateOverAllAlignedCellsInSector()
	c := uint(0)
	for p1, p2, ok := iter(); ok; p1, p2, ok = iter() {
		fmt.Printf("Aligned cells (%d,%d) <-> (%d,%d)\n", p1.x, p1.y, p2.x, p2.y)
		if p1.x != p2.x && p1.y != p2.y {
			t.Errorf("Error; Aligned cells (%d,%d) <-> (%d,%d)\n", p1.x, p1.y, p2.x, p2.y)
		}
		c++
	}
	if c != 18 {
		t.Errorf("Sector cell iterator runs an incorrect amount (%d?)\n", c)
	}
}
