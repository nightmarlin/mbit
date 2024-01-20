// Command scroller creates a scrolling dot that moves left-to-right,
// top-to-bottom across the mbit display matrix.
package main

import (
	"time"

	mbitm "tinygo.org/x/drivers/microbitmatrix"
)

const (
	qLen int16   = 4
	tps  float32 = 10
)

func main() {
	display := mbitm.New()
	display.Configure(mbitm.Config{})

	maxW, maxH := display.Size()
	dLen := maxW * maxH
	if dLen == 0 || dLen <= qLen {
		return // exit - nowhere to draw
	}

	cursor := float32(qLen)

	prevT := time.Now()
	for {
		now := time.Now()
		deltaT := float32(now.Sub(prevT)) / float32(time.Second)
		prevT = now

		display.ClearDisplay()

		for i := int16(0); i < qLen; i++ {
			c := int16(cursor) - i
			if c < 0 {
				c = dLen
			}

			// todo: make dynamic
			var brightness = mbitm.BrightnessFull
			switch i {
			case 0:
				brightness = mbitm.Brightness9
			case 1:
				brightness = mbitm.Brightness7
			case 2:
				brightness = mbitm.Brightness5
			case 3:
				brightness = mbitm.Brightness3
			}

			display.SetPixel(c%maxH, c/maxH, brightness)
		}

		if err := display.Display(); err != nil {
			_ = 0 // brok. debug assignment for breakpoint.
		}

		cursor += deltaT * tps
		for cursor >= float32(dLen) {
			cursor -= float32(dLen)
		}
	}
}