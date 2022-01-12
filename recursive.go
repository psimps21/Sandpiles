package main

import (
	"runtime"
)

func (b *Board) RecTopple(r, c int) {
	if b.Count(r, c) < 4 {
		return
	}
	b.toppleCount++
	if b.toppleCount%1 == 0 {
		b.gifList = append(b.gifList, CopyBoard(b))
	}

	toppled, rm := b.Count(r, c)/4, b.Count(r, c)%4
	b.SetToPositionRC(r, c, rm)

	if r < b.Rows()-1 {
		b.AddToPositionRC(r+1, c, toppled) // bottom
		b.RecTopple(r+1, c)
	}
	if r > 0 {
		b.AddToPositionRC(r-1, c, toppled) // top
		b.RecTopple(r-1, c)
	}
	if c < b.Cols()-1 {
		b.AddToPositionRC(r, c+1, toppled) // right
		b.RecTopple(r, c+1)
	}
	if c > 0 {
		b.AddToPositionRC(r, c-1, toppled) // left
		b.RecTopple(r, c-1)
	}
}

func (b *Board) RecToppleDiag(r, c int) {
	if b.Count(r, c) < 4 {
		return
	}
	b.toppleCount++
	// if b.toppleCount%1 == 0 {
	// b.gifList = append(b.gifList, CopyBoard(b))
	// }

	toppled, rm := b.Count(r, c)/4, b.Count(r, c)%4
	b.SetToPositionRC(r, c, rm)

	if r < b.Rows()-1 && c < b.Cols()-1 {
		b.AddToPositionRC(r+1, c+1, toppled) // SE
		b.RecTopple(r+1, c+1)
	}
	if r > 0 && c < b.Cols()-1 {
		b.AddToPositionRC(r-1, c+1, toppled) // NE
		b.RecTopple(r-1, c+1)
	}
	if r > 0 && c > 0 {
		b.AddToPositionRC(r-1, c-1, toppled) // NW
		b.RecTopple(r-1, c-1)
	}
	if r < b.Rows()-1 && c > 0 {
		b.AddToPositionRC(r+1, c-1, toppled) // SW
		b.RecTopple(r+1, c-1)
	}
}

func (b *Board) RecToppleDouble(r, c int) {
	if b.Count(r, c) < 4 {
		return
	}
	b.toppleCount++
	// if b.toppleCount%1 == 0 {
	// 	b.gifList = append(b.gifList, CopyBoard(b))
	// }

	toppled, rm := b.Count(r, c)/4, b.Count(r, c)%4
	b.SetToPositionRC(r, c, rm)

	if r < b.Rows()-2 && c < b.Cols()-2 {
		b.AddToPositionRC(r+2, c+2, toppled) // SE
		b.RecTopple(r+2, c+2)
	}
	if r > 1 && c < b.Cols()-2 {
		b.AddToPositionRC(r-2, c+2, toppled) // NE
		b.RecTopple(r-2, c+2)
	}
	if r > 1 && c > 1 {
		b.AddToPositionRC(r-2, c-2, toppled) // NW
		b.RecTopple(r-2, c-2)
	}
	if r < b.Rows()-2 && c > 1 {
		b.AddToPositionRC(r+2, c-2, toppled) // SW
		b.RecTopple(r+2, c-2)
	}
}

func (b *Board) ToppleBoardRec() {
	for r := range b.board {
		for c := range b.board {
			if b.Count(r, c) >= 4 {
				b.RecToppleDouble(r, c)
			}
		}
	}
}

func (b *Board) ToppleRecParallel() {
	div := runtime.NumCPU()
	subsize := b.Rows() / div
	for i := 0; i < div; i++ {
		if i == div-1 { // bottom BoardStrip
			for r := i * subsize; r < b.Rows(); r++ {
				for c := 0; c < b.Cols(); c++ {
					if b.Count(r, c) > 3 {
						go b.RecTopple(r, c)
					}
				}
			}
		} else {
			for r := i * subsize; r < (i+1)*subsize; r++ {
				for c := 0; c < b.Cols(); c++ {
					if b.Count(r, c) > 3 {
						go b.RecTopple(r, c)
					}
				}
			}
		}
	}
}
