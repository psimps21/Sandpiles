package main

import (
	"fmt"
	"math/rand"
	"runtime"
)

type Board struct {
	board       [][]int
	toppleCount int
	gifList     []*Board
}

// InitializeBoard returns a board struct with a board of given size
func InitializeBoard(rowSize, colSize int) *Board {
	board := make([][]int, rowSize)
	for i := range board {
		board[i] = make([]int, colSize)
	}
	return &Board{board: board}
}

//SetStartBoard returns a pointer to a board with the given starting conditions
func SetStartBoard(size, pile int, placement string) *Board {
	board := InitializeBoard(size, size)
	switch placement {
	case "central":
		cent := (size / 2) - 1
		board.SetToPositionRC(cent, cent, pile)
	case "random":
		if size >= 100 {
			// coinPerPos, extra := pile/100, pile%100
			randRows, randCols := rand.Perm(size), rand.Perm(size)
			var count int
			for i := range randRows {
				if i == 99 {
					board.AddToPositionRC(randRows[i], randCols[i], pile-count)
					break
				}
				if count < pile-1 {
					v := rand.Intn(pile - count + 1)
					board.AddToPositionRC(randRows[i], randCols[i], v)
					count += v
				} else {
					board.AddToPositionRC(randRows[i], randCols[i], 1)
					break
				}
			}
		} else {
			randRows, randCols := make([]int, 100), make([]int, 100)
			// Get 100 randomly sampled positions
			for i := range randRows {
				randRows[i], randCols[i] = rand.Intn(size), rand.Intn(size)
			}

			// Add one coin to each position
			var count int
			for i := range randRows {
				if i == 99 {
					board.AddToPositionRC(randRows[i], randCols[i], pile-count)
					break
				}
				if count < pile-1 {
					v := rand.Intn(pile - count + 1)
					board.AddToPositionRC(randRows[i], randCols[i], v)
					count += v
				} else {
					board.AddToPositionRC(randRows[i], randCols[i], 1)
					break
				}
			}
		}
	}
	return board
}

// PrintBoard Prints the values of a board
func PrintBoard(b *Board) {
	for i := range b.board {
		fmt.Println(b.board[i])
	}
}

// CopyBoard returns a copy of a board
func CopyBoard(b *Board) *Board {
	newBoard := InitializeBoard(b.Rows(), b.Cols())
	for i := range newBoard.board {
		for j := range newBoard.board {
			newBoard.board[i][j] = b.board[i][j]
		}
	}
	return newBoard
}

// SplitBoard returns a list of substrips spanning the entire board
func SplitBoard(b *Board, div int) []*BoardStrip {
	subSize := b.Rows() / div
	if subSize < 2 {
		panic("Subdivision of board is too small")
	}
	// Create the sub boards
	boardStrips := make([]*BoardStrip, div)
	for i := 0; i < div; i++ {
		if i == div-1 { // bottom BoardStrip
			subboard := b.board[i*subSize:]
			boardStrips[i] = InitializeBoardStrip(subboard, "bottom")
		} else {
			subboard := b.board[(i * subSize) : (i+1)*subSize]
			if i == 0 {
				boardStrips[i] = InitializeBoardStrip(subboard, "top")
			} else {
				boardStrips[i] = InitializeBoardStrip(subboard, "internal")
			}
		}
	}

	// Set BoardStrip neighbors
	for i, bStrp := range boardStrips {
		if i == 0 { // top stip
			bStrp.SetBottomStripNeighbor(boardStrips[i+1])
		} else if i == len(boardStrips)-1 {
			bStrp.SetTopStripNeighbor(boardStrips[i-1])
		} else {
			bStrp.SetInternalStripNeighbor(boardStrips[i-1], boardStrips[i+1])
		}
	}

	return boardStrips
}

// CombineBoardStrips returns a board from a list of BoardStrips
func CombineBoardStrips(strips []*BoardStrip) [][]int {
	var rowCount int
	newBoard := make([][]int, strips[0].Cols())
	for _, bs := range strips {
		for r := range bs.board {
			newBoard[rowCount] = bs.board[r]
			rowCount++
		}
	}
	return newBoard
}

// UpdateBoardParallel returns true if the board changed after updating all cells in parallel
func (b *Board) UpdateBoardParallel() bool {
	status, div := false, runtime.NumCPU()
	// Spilt board into BoardStrips
	boardStrips := SplitBoard(b, div)

	// Update BoardStrips
	stripStatus := make(chan bool, 4)
	for _, bs := range boardStrips {
		go bs.UpdateBoardStrip(stripStatus)
	}

	for range boardStrips {
		if <-stripStatus {
			status = true
		}
	}

	for _, bs := range boardStrips {
		if len(bs.topCells) != 0 && bs.topNeighbor != nil {
			bs.topNeighbor.UpdateBorder(bs.topCells, "top")
		}
		if len(bs.botCells) != 0 && bs.bottomNeighbor != nil {
			bs.bottomNeighbor.UpdateBorder(bs.botCells, "bottom")
		}
	}

	b.board = CombineBoardStrips(boardStrips)

	return status
}

/*
	Board Methods
*/

// Rows returns the number of rows in a board
func (bs *Board) Rows() int {
	return len(bs.board)
}

// Cols returns the number of columns in a Board
func (bs *Board) Cols() int {
	return len(bs.board[0])
}

// Count returns the count at position r,c
func (b *Board) Count(r, c int) int {
	return b.board[r][c]
}

// SetToPositionRC sets the number of coins for a position on the board (r,c)
func (b *Board) SetToPositionRC(r, c, count int) {
	b.board[r][c] = count
}

// AddToPositionRC adds a number of coins, count, to a position on the board (r,c)
func (b *Board) AddToPositionRC(r, c, count int) {
	b.board[r][c] += count
}

// UpdateCell updates a cell based on the number of coins present
func (b *Board) UpdateCell(r, c int) bool {
	if b.Count(r, c) < 4 {
		return false
	}

	toppled, rm := b.Count(r, c)/4, b.Count(r, c)%4
	for _, n := range b.FindNeighbors(r, c) {
		neighborR, neighborC := n[0], n[1]
		b.AddToPositionRC(neighborR, neighborC, toppled)
	}
	b.SetToPositionRC(r, c, rm)

	return true
}

// FindNeighbors returns a 2d array of ints corresponding to the coordinates of all viable neighbors
func (b *Board) FindNeighbors(r, c int) [][]int {
	var neighbors *[][]int
	if 0 < r && r < b.Rows()-1 && 0 < c && c < b.Cols()-1 { // in central board
		neighbors = &[][]int{{r - 1, c}, {r, c + 1}, {r + 1, c}, {r, c - 1}} // top,right,bottom,left
	} else if 0 < r && r < b.Rows()-1 && c == 0 { // on left border
		neighbors = &[][]int{{r - 1, c}, {r, c + 1}, {r + 1, c}}
	} else if 0 < r && r < b.Rows()-1 && c == b.Cols()-1 { // on right border
		neighbors = &[][]int{{r - 1, c}, {r + 1, c}, {r, c - 1}}
	} else if r == 0 && 0 < c && c < b.Cols()-1 { // on top border
		neighbors = &[][]int{{r, c + 1}, {r + 1, c}, {r, c - 1}}
	} else if r == b.Rows()-1 && 0 < c && c < b.Cols()-1 { // on bottom border
		neighbors = &[][]int{{r - 1, c}, {r, c + 1}, {r, c - 1}}
	} else if r == 0 && c == 0 { // in top left corner
		neighbors = &[][]int{{r, c + 1}, {r + 1, c}}
	} else if r == 0 && c == b.Cols()-1 { // in top right corner
		neighbors = &[][]int{{r + 1, c}, {r, c - 1}}
	} else if r == b.Rows()-1 && c == 0 { // in bottom left corner
		neighbors = &[][]int{{r - 1, c}, {r, c + 1}}
	} else if r == b.Rows()-1 && c == b.Cols()-1 { // in bottom right corner
		neighbors = &[][]int{{r - 1, c}, {r, c - 1}}
	}
	return *neighbors
}

// UpdateSerial
func (b *Board) UpdateSerial() bool {
	iterStatus := false
	for r := range b.board {
		for c := range b.board {
			if b.Count(r, c) >= 4 {
				iterStatus = b.UpdateCell(r, c)
			}
		}
	}
	return iterStatus
}
