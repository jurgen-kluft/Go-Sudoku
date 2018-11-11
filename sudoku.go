package main

import (
	"fmt"
)

type SudokuPuzzle3x3 [9][9]uint8

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
	puzzle [9][9]uint8
	board  [9][9]uint64
	sector [3][3]uint64
}

type SudokuSolver struct {
	board *SudokuBoard
}

func (s *SudokuBoard) prepare(puzzle SudokuPuzzle3x3) {
	s.width = 9
	s.height = 9

	// Every sector acts as a higher level ADD of every cell
	// in that sector
	everyNumberOnce := uint64(0)
	everyNumberMax := uint64(0)
	for n := uint(0); n < 9; n++ {
		everyNumberOnce = everyNumberOnce | uint64(1)<<(n*4)
		everyNumberMax = everyNumberMax | uint64(9)<<(n*4)
	}
	for y, row := range s.sector {
		for x := range row {
			s.sector[y][x] = everyNumberMax
		}
	}
	for y, row := range s.board {
		for x := range row {
			s.board[y][x] = everyNumberOnce
		}
	}

	// Start populate the board with the existing numbers from the puzzle
	for y, row := range puzzle {
		for x, n := range row {
			if n != 0 {
				//fmt.Printf("Fill in number '%v' at (%v,%v)\n", n, x, y)
				s.fillInNumber(uint(x), uint(y), n)
			}
		}
	}
}

func (s *SudokuBoard) removeNumber(x, y uint, number uint8) {
	maskNumberIsolate := (uint64(0xF) << (4 * (number - 1)))
	if s.board[y][x]&maskNumberIsolate != 0 {
		// Remove this number possibility from the board-cell and sector-cell
		s.board[y][x] &= ^maskNumberIsolate
		sx := x / 3
		sy := y / 3
		s.sector[sy][sx] -= 1 << (4 * (number - 1))
	}
}

func (s *SudokuBoard) fillInNumber(x, y uint, number uint8) {
	// Sector [x,y]
	sx := x / 3
	sy := y / 3

	// Remove this number from every cell in its sector
	for iy := uint(0); iy < 3; iy++ {
		by := sy*3 + iy
		for ix := uint(0); ix < 3; ix++ {
			bx := sx*3 + ix
			s.removeNumber(bx, by, number)
		}
	}
	// Go left<->right and remove the number from those cells
	for ix := uint(0); ix < s.width; ix++ {
		s.removeNumber(ix, y, number)
	}
	// Go up<->down and remove the number from those cells
	for iy := uint(0); iy < s.height; iy++ {
		s.removeNumber(x, iy, number)
	}

	// This cell on the board is set, so all other possibilities that
	// this cell contained should be removed from the sector cell.
	s.sector[sy][sx] -= s.board[y][x]
	s.board[y][x] = 0

	// Fill in the 'found' number on the puzzle
	s.puzzle[y][x] = number
}

func (s *SudokuBoard) printSector() {
	// PRINT LIKE THIS
	// 1,2,3
	// 1,2,3
	// 1,2,3

	for _, row := range s.sector {
		for x, n := range row {
			if x == 1 || x == 2 {
				fmt.Print(",")
			}
			fmt.Printf("%09x", n)
		}
		fmt.Println()
	}
}

func (v *SudokuSolver) iterate() (hasNumber bool, x uint, y uint, number uint8) {
	// Find a 'changed' sector and scan all cells in that sector to see if
	// we can find a cell with a 'single' number
	for sy, srow := range v.board.sector {
		for sx, sn := range srow {
			if sn != 0 {
				for n := uint(0); n < 9; n++ {
					if (sn & (uint64(0xF) << (n * 4))) == (uint64(1) << (n * 4)) {
						// Here we have a sector that has a number only occuring once
						for iy := uint(0); iy < 3; iy++ {
							for ix := uint(0); ix < 3; ix++ {
								if (v.board.board[uint(sy*3)+iy][uint(sx*3)+ix] & (uint64(0xF) << (n * 4))) == (uint64(1) << (n * 4)) {
									return true, uint(sx*3) + ix, uint(sy*3) + iy, uint8(n + 1)
								}
							}
						}
					}
				}
			}
		}
	}

	return false, 0, 0, 0
}

// Solve does all the work and will return a solved puzzle
func (s *SudokuBoard) Solve(puzzle SudokuPuzzle3x3) (solvedPuzzle SudokuPuzzle3x3) {
	s.prepare(puzzle)
	solver := SudokuSolver{board: s}
	for true {
		found, x, y, n := solver.iterate()
		if found {
			s.fillInNumber(x, y, n)
		} else {
			break
		}
	}

	for y, row := range s.puzzle {
		for x, n := range row {
			solvedPuzzle[y][x] = n
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
	example.print()
	fmt.Println()
	sudoku := &SudokuBoard{}
	solved := sudoku.Solve(example)
	solved.print()
}
