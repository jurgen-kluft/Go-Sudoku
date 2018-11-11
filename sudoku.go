package main

import (
	"fmt"
)

type SudokuPuzzle3x3 [9][9]uint8

var exampleMed = SudokuPuzzle3x3{
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

var exampleHard = SudokuPuzzle3x3{
	{0, 0, 0, 0, 0, 0, 6, 8, 0},
	{0, 0, 0, 0, 7, 3, 0, 0, 9},
	{3, 0, 9, 0, 0, 0, 0, 4, 5},
	{4, 9, 0, 0, 0, 0, 0, 0, 0},
	{8, 0, 3, 0, 5, 0, 9, 0, 2},
	{0, 0, 0, 0, 0, 0, 0, 3, 6},
	{9, 6, 0, 0, 0, 0, 3, 0, 8},
	{7, 0, 0, 6, 8, 0, 0, 0, 0},
	{0, 2, 8, 0, 0, 0, 0, 0, 0},
}

type SudokuBoard struct {
	width  uint
	height uint
	puzzle [9][9]uint8
	board  [9][9][9]uint8
}

type SudokuSolver struct {
	board *SudokuBoard
}

func (s *SudokuBoard) printBoard() {
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

	for y, row := range s.board {
		if y == 3 || y == 6 {
			fmt.Println("------------------------------+-------------------------------+------------------------------")
		}
		for x, p := range row {
			if x == 1 || x == 2 || x == 4 || x == 5 || x == 7 || x == 8 {
				fmt.Print(",")
			} else if x == 3 || x == 6 {
				fmt.Print(" | ")
			}
			for n := int(8); n >= 0; n-- {
				if p[n] == 0 {
					fmt.Print("0")
				} else {
					fmt.Print("1")
				}
			}
		}
		fmt.Println()
	}
	fmt.Println()
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
	fmt.Println()
}

func (s *SudokuBoard) prepare(puzzle SudokuPuzzle3x3) {
	s.width = 9
	s.height = 9

	for sy := uint(0); sy < s.height; sy++ {
		for sx := uint(0); sx < s.width; sx++ {
			for n := uint(0); n < 9; n++ {
				s.board[sy][sx][n] = 1
			}
		}
	}
	s.printBoard()

	// Start populate the board with the existing numbers from the puzzle
	for sy := uint(0); sy < s.height; sy++ {
		for sx := uint(0); sx < s.width; sx++ {
			n := uint(puzzle[sy][sx])
			if n != 0 {
				//fmt.Printf("Fill in number '%v' at (%v,%v)\n", n, x, y)
				s.fillInNumber(sx, sy, n)
			}
		}
	}
	s.printBoard()
}

func (s *SudokuBoard) fillInNumber(x, y uint, number uint) {
	if s.puzzle[y][x] == 0 {
		// Go left<->right and remove the number from those cells
		for ix := uint(0); ix < s.width; ix++ {
			s.board[y][ix][number-1] = 0
		}
		// Go up<->down and remove the number from those cells
		for iy := uint(0); iy < s.height; iy++ {
			s.board[iy][x][number-1] = 0
		}

		//TODO: There is another rule that states that is that if you
		// have 2 cells in a sector that are horizontal or vertical
		// aligned then any same number in both cells should be eliminated
		// in the adjacent cells in the direction of their alignment.

		// Sector [x,y]
		sx := (x / 3) * 3
		sy := (y / 3) * 3

		// Remove this number from any cell in this sector
		for iy := uint(0); iy < 3; iy++ {
			for ix := uint(0); ix < 3; ix++ {
				s.board[sy+iy][sx+ix][number-1] = 0
			}
		}

		// Empty this cell
		for n := uint(0); n < 9; n++ {
			s.board[y][x][n] = 0
		}
		// Fill in the 'found' number on the puzzle
		s.puzzle[y][x] = uint8(number)
	}
}

func (s *SudokuBoard) getCellInSectorWithOnlyOneNumberOccurance(sx, sy uint) (has bool, x uint, y uint, number uint) {
	sx = (sx / 3) * 3
	sy = (sy / 3) * 3

	coverage := [9]uint8{}
	for iy := uint(0); iy < 3; iy++ {
		for ix := uint(0); ix < 3; ix++ {
			cell := s.board[sy+iy][sx+ix]
			for n := uint(0); n < 9; n++ {
				coverage[n] = coverage[n] + cell[n]
			}
		}
	}

	for n := uint(0); n < 9; n++ {
		if coverage[n] == 1 {
			// Here we have a sector that has a number only occuring once
			// Find the sector-cell that has this number
			for iy := uint(0); iy < 3; iy++ {
				for ix := uint(0); ix < 3; ix++ {
					cell := s.board[sy+iy][sx+ix]
					if cell[n] == 1 {
						return true, sx + ix, sy + iy, uint(n + 1)
					}
				}
			}
		}
	}

	return false, 0, 0, 0
}

func (s *SudokuBoard) getCellInSectorWithOnlyOneNumber(sx, sy uint) (has bool, x uint, y uint, number uint) {
	sx = (sx / 3) * 3
	sy = (sy / 3) * 3

	for iy := uint(0); iy < 3; iy++ {
		for ix := uint(0); ix < 3; ix++ {
			cell := s.board[sy+iy][sx+ix]
			ones := 0
			for n := uint(0); n < 9; n++ {
				if cell[n] == 1 {
					ones++
					number = n + 1
				}
			}
			if ones == 1 {
				return true, sx + ix, sy + iy, number
			}
		}
	}
	return false, 0, 0, 0
}

func (s *SudokuBoard) solveStep() (hasNumber bool, x uint, y uint, number uint) {
	// Find a 'changed' sector and scan all cells in that sector to see if
	// we can find a cell with a 'single' number
	for sy := uint(0); sy < 3; sy++ {
		for sx := uint(0); sx < 3; sx++ {
			has, cx, cy, cn := s.getCellInSectorWithOnlyOneNumberOccurance(sx*3, sy*3)
			if has {
				return true, cx, cy, cn
			}
			one, ox, oy, on := s.getCellInSectorWithOnlyOneNumber(sx*3, sy*3)
			if one {
				return true, ox, oy, on
			}
		}
	}

	return false, 0, 0, 0
}

// Solve does all the work and will return a solved puzzle
func (s *SudokuBoard) Solve(puzzle SudokuPuzzle3x3) (solvedPuzzle SudokuPuzzle3x3) {
	s.prepare(puzzle)

	maxsteps := 50

	for true {
		found, x, y, n := s.solveStep()
		if found {
			fmt.Printf("Fill number '%d' at (%d,%d)\n", n, x+1, y+1)
			s.fillInNumber(x, y, n)
			if maxsteps == 0 {
				break
			}
			maxsteps--
		} else {
			s.printBoard()
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

func main() {
	puzzle := exampleHard
	/// ----------------------------
	puzzle.print()
	/// ----------------------------
	sudoku := &SudokuBoard{}
	solved := sudoku.Solve(puzzle)
	/// ----------------------------
	solved.print()
}
