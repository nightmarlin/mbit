TARGET := $(shell go run github.com/nightmarlin/mbitselect@v1 -fallback=microbit)

deploy:
	tinygo flash -target=$(TARGET) .

build:
	tinygo build -target=$(TARGET) -o=snake.hex ./main.go

.PHONY: deploy build
