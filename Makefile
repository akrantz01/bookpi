BINARY := bookpi
VERSION ?= vlatest

all: help

.PHONY: help
help:
	@echo
	@echo "List of commands:"
	@echo
	@echo "  make help	      - display this message"
	@echo "  make dist        - create a tarball of all necessary files"
	@echo "  make build	      - build all projects"
	@echo "  make clean	      - remove all generated files"
	@echo "  make clean-dist  - remove files for distribution"

.PHONY: dist
dist: build
	mkdir -p bookpi
	cp releases/* bookpi/
	cp services/* bookpi/
	cp install.sh bookpi/
	cp bookpi.sh bookpi/
	tar caf bookpi.tar.gz bookpi/*

.PHONY: build
build:
	mkdir -p releases
	$(MAKE) -C display build
	mv display/dist/main releases/$(BINARY)-$(VERSION)-display
	$(MAKE) -C frontend build
	mv -f frontend/build server/build
	$(MAKE) -C server build
	mv server/server releases/$(BINARY)-$(VERSION)-server

.PHONY: clean
clean: clean-dist
	$(MAKE) -C server clean
	$(MAKE) -C frontend clean
	$(MAKE) -C display clean

.PHONY: clean-dist
clean-dist:
	rm -rf releases
	rm -rf bookpi
	rm -rf bookpi.tar.gz
