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
	rm -rf frontend.gen.go

.PHONY: deps
deps:
	$(GOGET)

.PHONY: generate
generate: deps
	$(GOGENERATE) .

.PHONY: build
build: generate
	GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD)
