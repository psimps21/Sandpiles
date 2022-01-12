package main

import (
	"canvas"
	"fmt"
	"image"
)

//DrawToCanvas generates the image corresponding to a canvas after drawing a Board
//object's bodies on a square canvas that is canvasWidth pixels x canvasWidth pixels.
func (b *Board) DrawToCanvas() image.Image {

	cellWidth := 5
	height := b.Rows() * cellWidth
	width := b.Cols() * cellWidth

	cv := canvas.CreateNewPalettedCanvas(width, height, nil)

	black := canvas.MakeColor(0, 0, 0)
	g1 := canvas.MakeColor(85, 85, 85)
	g2 := canvas.MakeColor(170, 170, 170)
	white := canvas.MakeColor(255, 255, 255)
	ran := canvas.MakeColor(214, 159, 109)

	for r := range b.board {
		for c := range b.board {
			val := b.Count(r, c)
			if val == 0 {
				cv.SetFillColor(black)
			} else if val == 1 {
				cv.SetFillColor(g1)
			} else if val == 2 {
				cv.SetFillColor(g2)
			} else if val == 3 {
				cv.SetFillColor(white)
			} else {
				cv.SetFillColor(ran)
			}
			x := r * cellWidth
			y := c * cellWidth
			cv.ClearRect(x, y, x+cellWidth, y+cellWidth)
			cv.Fill()
		}
	}
	return canvas.GetImage(cv)
}

// DrawBoards creates an image for a GameBoards in a given list
func DrawBoards(boards []*Board) []image.Image {
	numSteps := len(boards)
	imgList := make([]image.Image, numSteps)
	for i := 0; i < numSteps; i++ { //  := range boards {
		// if i%500 == 0 {
		// fmt.Println(i, "boards drawn")
		imgList[i] = boards[i].DrawToCanvas()
		// }
		// imgList[i] = boards[i].DrawToCanvas()
	}
	fmt.Println("Drawing Done")
	return imgList
}
