// Command snake plays classic snake! Eat the pips until you can eat no more!
//
// The snake starts with a length of 1, eat pips to add to your length, but
// don't eat your own tail.
//
// Controls:
//
//	Button A: Turn counter-clockwise
//	Button B: Turn clockwise
package main

import (
	"machine"
	"time"

	mbitm "tinygo.org/x/drivers/microbitmatrix"

	"github.com/nightmarlin/mbit/buttons"
	"github.com/nightmarlin/mbit/deltat"
)

const tps float32 = 4

var dead bool

func main() {
	display := mbitm.New()
	display.Configure(mbitm.Config{})

	h, w := display.Size()
	if h*w < 0 {
		return // bad size
	}

	s := snake{
		head: snakeCell{
			p: point{pointX: 2, pointY: 4 /*bottom*/},
		},
		dir: dirUp,
		len: 1,
	}

	rotateChan := make(chan rotation, 1)
	if err := buttons.TrySendOnPress(machine.BUTTONA, rotateChan, anticlockwise); err != nil {
		return
	}
	if err := buttons.TrySendOnPress(machine.BUTTONB, rotateChan, clockwise); err != nil {
		return
	}

	pip := point{pointX: 2, pointY: 1 /*middle-1*/}

	var lastTick float32
	deltat.Loop(
		deltat.DefaultClock(),
		deltat.ByTick(tps),
		func(now time.Time, delta float32) (cont bool) {
			currentTick := lastTick + delta
			defer func() { lastTick = currentTick }()

			if dead {
				if int(currentTick)%2 == 0 {
					display.EnableAll()
				} else {
					display.DisableAll()
				}
				return true
			}

			// process movement every 3 whole ticks
			if int(currentTick)-int(lastTick) != 0 && int(currentTick)%3 == 0 {
				// process rotations
				select {
				case rotate := <-rotateChan:
					s.dir = s.dir.Rotate(rotate)
				default:
				}

				s.Move()

				if s.CollidesSelf() {
					// todo death
					dead = true
					return true
				}

				// if head colliding with pip, grow snake & move pip
				if s.head.p == pip {
					s.Grow()

					for s.Collides(pip) {
						np, err := randPoint(point{5, 5})
						if err != nil {
							return false
						}
						pip = np
					}
				}
			}

			display.ClearDisplay()

			// render snake
			cell := &s.head
			for cell != nil {
				display.SetPixel(int16(cell.p[pointX]), int16(cell.p[pointY]), mbitm.BrightnessFull)
				cell = cell.next
			}

			// render pip on alternating ticks
			if int(currentTick)%2 == 0 {
				// render pip
				display.SetPixel(int16(pip[pointX]), int16(pip[pointY]), mbitm.BrightnessFull)
			}

			if err := display.Display(); err != nil {
				return false
			}
			return true
		},
	)
}

func randPoint(dims point) (point, error) {
	rand, err := machine.GetRNG()
	if err != nil {
		return [2]uint8{}, err
	}

	return point{
		pointX: uint8(rand),
		pointY: uint8(rand >> 8),
	}.norm(dims), nil
}
