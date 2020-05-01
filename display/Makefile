PIP := pip3
PYINSTALLER := pyinstaller

.PHONY: deps
deps:
	$(PIP) install adafruit-circuitpython-ssd1306 pyinstaller ifaddr

.PHONY: build
build: deps
	$(PYINSTALLER) main.spec

.PHONY: clean
clean:
	rm -rf build
	rm -rf dist
