GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOGET := $(GOCMD) get
GOGENERATE := $(GOCMD) generate

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf build
	rm -rf files
	rm -rf database.db
	rm -rf assets/fileb0x*

.PHONY: deps
deps: generate
	$(GOGET)

.PHONY: deps-embed
deps-embed:
	$(GOGET) -u github.com/UnnoTed/fileb0x

.PHONY: generate
generate: deps-embed
	test -d build || ( echo "Frontend must be built first"; exit 1 )
	$(GOGENERATE) .

.PHONY: build
build: deps
	GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD)
