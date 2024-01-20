// Command showall enables every LED in the display matrix
package main

import (
	mbitm "tinygo.org/x/drivers/microbitmatrix"
)

func main() {
	display := mbitm.New()
	display.Configure(mbitm.Config{})
	display.EnableAll()
}
