PYTHON := python3
PIP := pip3
PYINSTALLER := ./venv/bin/pyinstaller

.PHONY: deps
deps: venv

.PHONY: build
build: deps
	$(PYINSTALLER) main.spec

.PHONY: clean
clean:
	rm -rf build
	rm -rf dist

.PHONY: clean-venv
clean-venv:
	rm -rf venv
	rm -rf pyinstaller

venv: venv/bin/activate

venv/bin/activate: clean-venv
	test -d venv || $(PYTHON) -m venv ./venv
	. venv/bin/activate; pip install -r requirements.txt
	$(MAKE) pyinstaller
	touch venv/bin/activate

.PHONY: pyinstaller
pyinstaller:
	git clone https://github.com/pyinstaller/pyinstaller
	cd pyinstaller; git checkout v3.6
	. venv/bin/activate; cd pyinstaller/bootloader; $(PYTHON) ./waf all --no-lsb
	. venv/bin/activate; cd pyinstaller; $(PIP) install -e .
