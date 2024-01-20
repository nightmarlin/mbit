// Command scroller creates a scrolling dot that moves left-to-right,
// top-to-bottom across the mbit display matrix.
package main

import (
	"image/color"
	"machine"
	"time"

	mbitm "tinygo.org/x/drivers/microbitmatrix"

	"github.com/nightmarlin/mbit/deltat"
)

const tps float32 = 15

var tailLen uint8 = 4

func main() {
	display := mbitm.New()
	display.Configure(mbitm.Config{})

	maxW, maxH := display.Size()
	dLen := maxW * maxH
	if dLen == 0 || dLen <= int16(tailLen) {
		return // exit - nowhere to draw
	}

	machine.BUTTONA.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	if err := machine.BUTTONA.SetInterrupt(
		machine.PinFalling,
		func(p machine.Pin) {
			// todo: changing qlen is slightly racy - fix that
			if tailLen != 1 {
				tailLen--
			}
		},
	); err != nil {
		return
	}

	machine.BUTTONB.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	if err := machine.BUTTONB.SetInterrupt(
		machine.PinFalling,
		func(p machine.Pin) {
			if tailLen < uint8(dLen) {
				tailLen++
			}
		},
	); err != nil {
		return
	}

	cursor := float32(tailLen)

	deltat.Loop(
		deltat.DefaultClock(),
		deltat.ByTick(tps),
		func(_ time.Time, elapsedTicks float32) bool {
			curTailLen := tailLen // minimise raciness by reading tailLen at start of loop

			brightnessDiv := uint8(int8(255 / curTailLen))
			display.ClearDisplay()

			for i := uint8(0); i < curTailLen; i++ {
				c := int16(cursor) - int16(i)
				if c < 0 {
					c = dLen - int16(curTailLen-i)
				}

				display.SetPixel(
					c%maxH,
					c/maxH,
					color.RGBA{
						R: 255, G: 255, B: 255,
						A: 255 - brightnessDiv*(curTailLen-i),
					},
				)
			}

			if err := display.Display(); err != nil {
				return false
			}

			cursor += elapsedTicks
			for cursor >= float32(dLen) {
				cursor -= float32(dLen)
			}
			return true
		},
	)
}
