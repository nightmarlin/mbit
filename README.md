# mbit

An experiment into tinygo and how it works on the BBC Micro:Bit platform.

> [!important]
> tinygo only runs on a single thread. goroutines must either perform IO or make
> intrinsically blocking calls to trigger other goroutines. an alternative is
> calling `runtime.Gosched()`

## Useful links

- MBit docs: https://tinygo.org/docs/reference/microcontrollers/microbit/
- TinyGo tips & gotchas: https://tinygo.org/docs/guides/tips-n-tricks/
- TinyGo bluetooth: https://github.com/tinygo-org/bluetooth/
- TinyGo concepts: https://tinygo.org/docs/concepts/

## Setup

```sh
brew tap tinygo-org/tools
brew install tinygo
```

point your env's go installation to the tinygo install.

> IJ IDEs can use the TinyGo plugin to automate this

Plug in a macrobit, and then deploy away! It Just Works!
