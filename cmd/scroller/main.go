// Command scroller creates a scrolling dot that moves left-to-right,
// top-to-bottom across the mbit display matrix.
package main

import (
	"image/color"
	"time"

	mbitm "tinygo.org/x/drivers/microbitmatrix"

	"github.com/nightmarlin/mbit/deltat"
)

const (
	qLen           uint8   = 4
	brightnessDivs         = uint8(int8(255 / qLen))
	tps            float32 = 10
)

func main() {
	display := mbitm.New()
	display.Configure(mbitm.Config{})

	maxW, maxH := display.Size()
	dLen := maxW * maxH
	if dLen == 0 || dLen <= int16(qLen) {
		return // exit - nowhere to draw
	}

	cursor := float32(qLen)

	deltat.Loop(
		deltat.DefaultClock(),
		deltat.ByTick(tps),
		func(_ time.Time, elapsedTicks float32) bool {
			display.ClearDisplay()

			for i := uint8(0); i < qLen; i++ {
				c := int16(cursor) - int16(i)
				if c < 0 {
					c = dLen - int16(qLen-i)
				}

				display.SetPixel(
					c%maxH,
					c/maxH,
					color.RGBA{
						R: 255, G: 255, B: 255,
						A: 255 - brightnessDivs*(qLen-i),
					},
				)
			}

			if err := display.Display(); err != nil {
				_ = 0 // bogus assignment for debugger to break on.
			}

			cursor += elapsedTicks
			for cursor >= float32(dLen) {
				cursor -= float32(dLen)
			}
			return true
		},
	)
}
