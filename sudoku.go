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

var alignmentChecks = [3][3]uint8{
	{1, 1, 0},
	{0, 1, 1},
	{1, 0, 1},
}

type point struct {
	x uint
	y uint
}

// Use like this:
// p := point{0,0}
// iter := p.iterateOverAllCellsInSector()
// for p, ok := iter(); ok; p, ok = iter() {
//     fmt.Printf("Sector iterator: cell(%d,%d)", p.x,p.y)
// }
func (p point) iterateOverAllCellsInSector() func() (point, bool) {
	anchor := point{x: (p.x / 3) * 3, y: (p.y / 3) * 3}
	iter := -1
	return func() (point, bool) {
		iter++
		return point{x: anchor.x + uint(iter%3), y: anchor.y + uint(iter/3)}, (iter < 9)
	}
}

// Sector
// 0,1,2
// 3,4,5
// 6,7,8
// Possible aligned pairs(p1,p2)
// 0 1
// 0 2
// 0 3
// 0 6
// 1 2
// 1 4
// 1 7
// 2 5
// 2 8
// 3 4
// 3 5
// 3 6
// 4 5
// 4 7
// 5 8
// 6 7
// 6 8
// 7 8
func (p point) iterateOverAllAlignedCellsInSector() func() (point, point, bool) {
	anchor := point{x: (p.x / 3) * 3, y: (p.y / 3) * 3}
	c1 := -1
	c2 := 6
	d := 3
	return func() (point, point, bool) {
		if d == 1 && (c2%3) == 2 {
			c2 = c1
			d = 3
		}
		if d == 3 && (c2+d) > 8 {
			c1++
			if c1%3 == 2 {
				d = 3
			} else {
				d = 1
			}
			c2 = c1
		}
		c2 += d
		p1 := point{x: anchor.x + uint(c1%3), y: anchor.y + uint(c1/3)}
		p2 := point{x: anchor.x + uint((c2)%3), y: anchor.y + uint((c2)/3)}
		return p1, p2, (c1 != 8)
	}
}

func (s *SudokuBoard) countOccurencesInSector(p point, number uint, ignored func(p point) bool) (normalCount uint8, ignoredCount uint8) {
	normalCount = 0
	ignoredCount = 0
	iter := p.iterateOverAllCellsInSector()
	for p, ok := iter(); ok; p, ok = iter() {
		if ignored(p) == false {
			normalCount += s.board[p.y][p.x][number-1]
		} else {
			ignoredCount += s.board[p.y][p.x][number-1]
		}
	}
	return
}

func (s *SudokuBoard) checkAlignmentRule(x uint, y uint, n uint) (found bool, sp point, sd [2]uint) {
	p := point{x: x, y: y}
	iter := p.iterateOverAllAlignedCellsInSector()
	for p1, p2, ok := iter(); ok; p1, p2, ok = iter() {
		ignore := func(p point) bool {
			// This is to ignore these 2 cells from computing the occurences
			return (p.x == p1.x && p.y == p1.y) || (p.x == p2.x && p.y == p2.y)
		}
		no, ni := s.countOccurencesInSector(p, n, ignore)
		if no == 0 && ni == 2 {
			// There are no normal occurences of this number but the two cells we checked
			// both have this number. This is what we are looking for.
			return true, point{x: p1.x, y: p1.y}, [2]uint{p2.x - p1.x, p2.y - p1.y}
		}
	}

	return false, point{}, [2]uint{0, 0}
}

func (s *SudokuBoard) eliminateNumberFromAllCellsOnAxisExceptCurrentSector(p point, axis [2]uint, number uint) {
	if axis[0] == 1 {
		psx := p.x / 3
		for sx := uint(0); sx < 3; sx++ {
			if sx != psx {
				for cx := uint(0); cx < 3; cx++ {
					bx := sx*3 + cx
					by := p.y
					s.board[by][bx][number-1] = 0
				}
			}
		}
	} else if axis[1] == 1 {
		psy := p.y / 3
		for sy := uint(0); sy < 3; sy++ {
			if sy != psy {
				for cy := uint(0); cy < 3; cy++ {
					by := sy*3 + cy
					bx := p.x
					s.board[by][bx][number-1] = 0
				}
			}
		}
	}
}

func (s *SudokuBoard) applyAlignmentRule() {
	for sy := uint(0); sy < 3; sy++ {
		for sx := uint(0); sx < 3; sx++ {
			for n := uint(0); n < 9; n++ {
				bx := sx * 3
				by := sy * 3
				found, p, axis := s.checkAlignmentRule(bx, by, n)
				if found {
					s.eliminateNumberFromAllCellsOnAxisExceptCurrentSector(p, axis, n)
				}
			}
		}
	}
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
	//s.printBoard()

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
	//s.printBoard()
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
