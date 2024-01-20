install-tinygo-macos:
	brew tap tinygo-org/tools
	brew install tinygo

deploy-scroller:
	$(MAKE) -C ./cmd/scroller deploy

.PHONY: deploy build build-stripped
