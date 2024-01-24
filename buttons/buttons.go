package buttons

import (
	"machine"
)

func TrySendOnPress[T any](btn machine.Pin, c chan<- T, msg T) error {
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
