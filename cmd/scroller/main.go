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

func trySendOnButtonPress[T any](btn machine.Pin, c chan<- T, msg T) error {
	btn.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	return btn.SetInterrupt(
		machine.PinFalling,
		func(machine.Pin) {
			select {
			case c <- msg:
			default:
			}
		},
	)
}

func main() {
	display := mbitm.New()
	display.Configure(mbitm.Config{})

	displayMaxW, displayMaxH := display.Size()
	displayLen := displayMaxW * displayMaxH
	if displayLen == 0 || displayLen <= int16(tailLen) {
		return // exit - nowhere to draw
	}

	tailLenDeltaChan := make(chan int8, 1)
	if err := trySendOnButtonPress(machine.BUTTONA, tailLenDeltaChan, -1); err != nil {
		return
	}
	if err := trySendOnButtonPress(machine.BUTTONB, tailLenDeltaChan, 1); err != nil {
		return
	}

	cursor := float32(tailLen)

	deltat.Loop(
		deltat.DefaultClock(),
		deltat.ByTick(tps),
		func(_ time.Time, elapsedTicks float32) bool {
			brightnessDiv := uint8(int8(255 / tailLen))
			display.ClearDisplay()

			for i := uint8(0); i < tailLen; i++ {
				c := int16(cursor) - int16(i)
				if c < 0 {
					c = displayLen - int16(tailLen-i)
				}

				display.SetPixel(
					c%displayMaxH,
					c/displayMaxH,
					color.RGBA{
						R: 255, G: 255, B: 255,
						A: 255 - brightnessDiv*(tailLen-i),
					},
				)
			}

			if err := display.Display(); err != nil {
				return false
			}

			select {
			case deltaTailLen := <-tailLenDeltaChan:
				if deltaTailLen < 0 {
					// safely cast negatives to positive uint and subtract
					tailLen -= uint8(-1 * deltaTailLen)
					if tailLen < 1 {
						tailLen = 1
					}
				} else {
					tailLen += uint8(deltaTailLen)
					if tailLen >= uint8(displayLen) {
						tailLen = uint8(displayLen) - 1
					}
				}
			default:
			}

			cursor += elapsedTicks
			for cursor >= float32(displayLen) {
				cursor -= float32(displayLen)
			}
			return true
		},
	)
}
