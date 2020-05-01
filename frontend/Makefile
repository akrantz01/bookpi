YARN := yarn

.PHONY: deps
deps:
	$(YARN) install

.PHONY: start
start: deps
	$(YARN) start

.PHONY: build
build: deps
	$(YARN) build

.PHONY: clean
clean:
	rm -rf node_modules
	rm -rf build
