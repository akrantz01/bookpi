PIP := pip3
PYINSTALLER := pyinstaller

.PHONY: deps
deps:
	$(PIP) install -r requirements.txt

.PHONY: build
build: deps
	$(PYINSTALLER) main.spec

.PHONY: clean
clean:
	rm -rf build
	rm -rf dist
