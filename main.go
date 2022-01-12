package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	size, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("Cannot convert string to integer")
	}

	pile, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic("Cannot convert string to integer")
	}

	placement := os.Args[3]

	/*
		Serial Solution
	*/

	serialBoard := SetStartBoard(size, pile, placement)
	// parallelBoard := CopyBoard(serialBoard)

	// iterStatus := true
	// count := 0

	start := time.Now()
	serialBoard.ToppleBoardRec()

	// for iterStatus {
	// 	iterStatus = false
	// 	count++
	// 	iterStatus = serialBoard.UpdateSerial()
	// }
	elap1 := time.Since(start).Seconds()

	fmt.Printf("Serial took %f seconds and %d topples\n", elap1, serialBoard.toppleCount)

	serialBoard.gifList = append(serialBoard.gifList, CopyBoard(serialBoard))
	serialImage := serialBoard.DrawToCanvas()
	f, err3 := os.Create("serial.png")
	if err3 != nil {
		log.Fatal(err3)
	}
	err4 := png.Encode(f, serialImage)
	if err4 != nil {
		f.Close()
		log.Fatal(err4)
	}

	/*
		Gif drawing
	*/

	// imgList := DrawBoards(serialBoard.gifList)
	// outputFile := "evoGif"
	// gifhelper.ImagesToGIF(imgList, outputFile)

	/*
		Parallel Solution
	*/
	// start := time.Now()
	// parallelBoard := SetStartBoard(size, pile, placement)
	// elap4 := time.Since(start).Seconds()
	// fmt.Printf("Took %f seconds to initialize board\n", elap4)
	// // Update board until it no longer changes
	// status := true
	// count = 0
	// start = time.Now()
	// parallelBoard.ToppleRecParallel()
	// for status {
	// 	count++
	// 	status = parallelBoard.UpdateBoardParallel()
	// }
	// elap2 := time.Since(start).Seconds()

	// fmt.Printf("Parallel took %f seconds with %d topples\n", elap2, parallelBoard.toppleCount)

	// start = time.Now()
	// parallelImage := parallelBoard.DrawToCanvas()
	// elap3 := time.Since(start).Seconds()
	// fmt.Printf("Took %f seconds to draw final board\n", elap3)

	// f, err3 = os.Create("parallel.png")
	// if err3 != nil {
	// 	log.Fatal(err3)
	// }

	// err4 = png.Encode(f, parallelImage)
	// if err4 != nil {
	// 	f.Close()
	// 	log.Fatal(err4)
	// }
}
