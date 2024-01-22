TARGET := $(shell go run github.com/nightmarlin/mbitselect@v1 -fallback=microbit)

deploy:
	tinygo flash -target=$(TARGET) ./main.go

build:
	tinygo build -target=$(TARGET) -o=showall.hex ./main.go

.PHONY: deploy build
