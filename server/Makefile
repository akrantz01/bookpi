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
deps: generate
	$(GOGET)

.PHONY: deps-broccoli
deps-broccoli:
	$(GOGET) -u aletheia.icu/broccoli

.PHONY: generate
generate: deps-broccoli
	test -d build || ( echo "Frontend must be built first"; exit 1 )
	$(GOGENERATE) .

.PHONY: build
build: deps
	GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD)
