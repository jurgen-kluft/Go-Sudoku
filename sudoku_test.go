package main

import "testing"

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

var emptyPuzzleSector = [3][3]uint64{
	{0x999999999, 0x999999999, 0x999999999},
	{0x999999999, 0x999999999, 0x999999999},
	{0x999999999, 0x999999999, 0x999999999},
}

func CompareSectors(a [3][3]uint64, b [3][3]uint64) bool {
	for y, r := range a {
		for x, v := range r {
			if b[y][x] != v {
				return false
			}
		}
	}
	return true
}

func TestSudokuBoard_prepare(t *testing.T) {
	type fields struct {
		width  uint
		height uint
		puzzle [9][9]uint8
		board  [9][9]uint64
		sector [3][3]uint64
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
				sector: tt.fields.sector,
			}
			s.prepare(tt.args.puzzle)
			if CompareSectors(s.sector, emptyPuzzleSector) == false {
				t.Error("Initialized sectors for empty puzzle is wrong!")
			}
		})
	}
}

func TestSudokuBoard_fillNumber(t *testing.T) {
	s := &SudokuBoard{}
	s.prepare(emptyPuzzle)
	s.fillInNumber(0, 0, 1)
	if s.sector[0][0] != 0x888888880 {
		t.Errorf("Filling in number '1' at '0,0' gives wrong result for sector[0][0] as %x!\n", s.sector[0][0])
	}
	s.fillInNumber(1, 0, 2)
	if s.sector[0][0] != 0x777777700 {
		t.Errorf("Filling in number '2' at '1,0' gives wrong result for sector[0][0] as %x!\n", s.sector[0][0])
	}
	s.fillInNumber(2, 0, 3)
	if s.sector[0][0] != 0x666666000 {
		t.Errorf("Filling in number '3' at '2,0' gives wrong result for sector[0][0] as %x!\n", s.sector[0][0])
	}
}
