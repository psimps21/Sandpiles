package main

// Cell holds the column and value of a toppled cell
type Cell struct {
	c, v int
}

// BoardStrip is a subsection of the board that has been padded with zeros to handle edge updating between subboards
type BoardStrip struct {
	Board
	topNeighbor    *BoardStrip
	bottomNeighbor *BoardStrip
	topCells       []*Cell
	botCells       []*Cell
}

/*
	BoardStrip Methods
*/

// SetTopStripNeighbor updates the bottom border for a BoardStrip
func (bs *BoardStrip) SetTopStripNeighbor(bottomN *BoardStrip) {
	bs.bottomNeighbor = bottomN
	bs.botCells = make([]*Cell, 0, bs.Cols())
}

// SetBottomStripNeighbor updates the bottom border for a BoardStrip
func (bs *BoardStrip) SetBottomStripNeighbor(topN *BoardStrip) {
	bs.topNeighbor = topN
	bs.topCells = make([]*Cell, 0, bs.Cols())
}

// SetInternalStripNeighbor updates the bottom border for a BoardStrip
func (bs *BoardStrip) SetInternalStripNeighbor(topN, bottomN *BoardStrip) {
	bs.bottomNeighbor = bottomN
	bs.botCells = make([]*Cell, 0, bs.Cols())
	bs.topNeighbor = topN
	bs.topCells = make([]*Cell, 0, bs.Cols())
}

// UpdateBSCell updates a cell in a BoardStrip based on the number of coins present
func (bs *BoardStrip) UpdateBSCell(r, c int) bool {
	if bs.Count(r, c) < 4 {
		return false
	}

	toppled, rm := bs.Count(r, c)/4, bs.Count(r, c)%4
	bs.SetToPositionRC(r, c, rm)
	for _, n := range bs.FindNeighbors(r, c) {
		neighborR, neighborC := n[0], n[1]
		bs.AddToPositionRC(neighborR, neighborC, toppled)
	}

	if r == 0 {
		if bs.topNeighbor != nil {
			bs.topCells = append(bs.topCells, &Cell{c, toppled})
		}
	} else if r == bs.Rows()-1 {
		if bs.bottomNeighbor != nil {
			bs.botCells = append(bs.botCells, &Cell{c, toppled})
		}
	}

	return true
}

// UpdateBorder updates the cells from that toppled into other subboards
func (bs *BoardStrip) UpdateBorder(cells []*Cell, loc string) {
	if loc == "bottom" {
		for _, cell := range cells {
			bs.AddToPositionRC(0, cell.c, cell.v)
		}
	} else if loc == "top" {
		for _, cell := range cells {
			bs.AddToPositionRC(bs.Rows()-1, cell.c, cell.v)
		}
	}
}

// UpdateParallel updates a subboard and returns true if any coins toppled
func (bs *BoardStrip) UpdateParallel() bool {
	iterStatus := false
	for r := 0; r < bs.Rows(); r++ {
		for c := 0; c < bs.Cols(); c++ {
			if bs.Count(r, c) >= 4 {
				iterStatus = bs.UpdateBSCell(r, c)
			}
		}
	}
	return iterStatus
}

// UpdateBoardStrip returns true if any cells in a BoardStrip have been updated in a generation
func (bs *BoardStrip) UpdateBoardStrip(finish chan bool) {
	updateStatus := false
	updateStatus = bs.UpdateParallel()
	finish <- updateStatus
}

/*
	Functions
*/

// InitializeBoardStrip retuens a pointer to a BoardStrip with a subboard
func InitializeBoardStrip(board [][]int, stripType string) *BoardStrip {
	bST := &BoardStrip{}
	bST.board = board
	return bST
}
