package main

import (
	"fmt"
)

type SudokuPuzzle3x3 [9][9]byte

var example = SudokuPuzzle3x3{
	{5, 3, 0, 0, 7, 0, 0, 0, 0},
	{6, 0, 0, 1, 9, 5, 0, 0, 0},
	{0, 9, 8, 0, 0, 0, 0, 6, 0},
	{8, 0, 0, 0, 6, 0, 0, 0, 3},
	{4, 0, 0, 8, 0, 3, 0, 0, 1},
	{7, 0, 0, 0, 2, 0, 0, 0, 6},
	{0, 6, 0, 0, 0, 0, 2, 8, 0},
	{0, 0, 0, 4, 1, 9, 0, 0, 5},
	{0, 0, 0, 0, 8, 0, 0, 7, 9},
}

type SudokuBoard struct {
	width  uint
	height uint
	puzzle [9][9]byte
	board  [9][9]uint64
	sector [3][3]byte
}

type SudokuSolver struct {
	x, y  int
	board *SudokuBoard
}

func (s *SudokuBoard) prepare(puzzle SudokuPuzzle3x3) {
	// Convert the actual 'non-zero' numbers to:
	// value = 1 << number
	s.width = 9
	s.height = 9
	for _, row := range s.sector {
		for x := range row {
			row[x] = 1
		}
	}
	allNumbers := uint64(0)
	for n := uint(1); n <= 9; n++ {
		allNumbers = allNumbers | (uint64(1) << ((n - 1) * 4))
	}
	for _, row := range s.board {
		for x := range row {
			row[x] = allNumbers
		}
	}

	// Start populate the board with the existing numbers from the puzzle
	for y, row := range puzzle {
		for x, n := range row {
			if n != 0 {
				s.fillInNumber(x, y, n)
			}
		}
	}
}

func (v *SudokuSolver) iterate() (hasNumber bool, x int, y int, number byte) {

	return false, 0, 0, 0
}

func (s *SudokuBoard) fillInNumber(x, y int, number byte) {
	// Fill in the number on the board
	s.puzzle[x][y] = number

	// This cell is done
	s.board[x][y] = 0

	// Go left<->right and remove the number from those cells
	numberMask := ^(uint64(0xF) << (number - 1))
	for ix := uint(0); ix <= s.width; ix++ {
		s.board[ix][y] = s.board[ix][y] & numberMask
	}
	// Go up<->down and remove the number from those cells
	for iy := uint(0); iy <= s.height; iy++ {
		s.board[x][iy] = s.board[x][iy] & numberMask
	}
}

func (s *SudokuBoard) Solve(puzzle SudokuPuzzle3x3) (solvedPuzzle SudokuPuzzle3x3) {
	s.prepare(puzzle)
	solver := SudokuSolver{x: 0, y: 0, board: s}
	for solved, x, y, n := solver.iterate(); !solved; {
		s.fillInNumber(x, y, n)
	}

	for y, row := range s.puzzle {
		for x, n := range row {
			solvedPuzzle[x][y] = n
		}
	}
	return
}

func (p SudokuPuzzle3x3) print() {

	// PRINT LIKE THIS
	// 1,2,3 | 4,5,6 | 7,8,9
	// 1,2,3 | 4,5,6 | 7,8,9
	// 1,2,3 | 4,5,6 | 7,8,9
	// ------+-------+------
	// 1,2,3 | 4,5,6 | 7,8,9
	// 1,2,3 | 4,5,6 | 7,8,9
	// 1,2,3 | 4,5,6 | 7,8,9
	// ------+-------+------
	// 1,2,3 | 4,5,6 | 7,8,9
	// 1,2,3 | 4,5,6 | 7,8,9
	// 1,2,3 | 4,5,6 | 7,8,9

	for y, row := range p {
		if y == 3 || y == 6 {
			fmt.Println("------+-------+------")
		}
		for x, n := range row {
			if x == 1 || x == 2 || x == 4 || x == 5 || x == 7 || x == 8 {
				fmt.Print(",")
			} else if x == 3 || x == 6 {
				fmt.Print(" | ")
			}
			fmt.Printf("%v", n)
		}
		fmt.Println()
	}
}

func main() {
	sudoku := &SudokuBoard{}
	solved := sudoku.Solve(example)
	solved.print()
}
