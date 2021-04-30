SHELL = /bin/sh

.PHONY: default linux windows windows64 clean clobber install docker_usb docker_usb_windows docker_usb_windows64 docker_usb_linux tag

EDITCP_SRC = *.go
UI_SRC = ../ui/*.go
CODEPLUG_SRC = ../codeplug/*.go
DFU_SRC = ../dfu/*.go
STDFU_SRC = ../stdfu/*.go
USERDB_SRC = ../userdb/*.go
DEBUG_SRC = ../debug/*.go
SOURCES = $(EDITCP_SRC) $(UI_SRC) $(CODEPLUG_SRC) $(DFU_SRC) $(STDFU_SRC) $(USERDB_SRC) $(DEBUG_SRC)
VERSION = $(shell sed -n '/version =/{s/^[^"]*"//;s/".*//p;q}' <version.go)

default: linux

linux: clean deploy/linux/editcp deploy/linux/editcp.sh deploy/linux/install deploy/linux/99-md380.rules

deploy/linux/editcp: $(SOURCES)
	go mod tidy
	go mod vendor
	qtdeploy -docker build linux
	rm -rf deploy/linux/qml

.PHONY: deploy/linux/editcp.sh	# Force, in case it's overwritten by install
deploy/linux/editcp.sh: editcp.sh
	cp editcp.sh deploy/linux/editcp.sh

deploy/linux/install: install.sh deploy/linux/editcp 99-md380.rules
	cp install.sh deploy/linux/install

deploy/linux/99-md380.rules: 99-md380.rules
	cp 99-md380.rules deploy/linux/

editcp-$(VERSION).tar.xz: linux
	rm -rf editcp-$(VERSION)
	mkdir -p editcp-$(VERSION)
	cp -al deploy/linux/* editcp-$(VERSION)
	tar cJf editcp-$(VERSION).tar.xz editcp-$(VERSION)
	rm -rf editcp-$(VERSION)

install: linux
	cd deploy/linux && ./install .

windows: clean editcp-$(VERSION)-installer.exe

windows64: clean editcp64-$(VERSION)-installer.exe

editcp-$(VERSION)-installer.exe: deploy/win32/editcp.exe editcp.nsi dll/*.dll
	makensis -DVERSION=$(VERSION) editcp.nsi

editcp64-$(VERSION)-installer.exe: deploy/win64/editcp64.exe editcp64.nsi dll/*.dll
	makensis -DVERSION=$(VERSION) editcp64.nsi

deploy/win32/editcp.exe: $(SOURCES)
	go mod tidy
	go mod vendor
	qtdeploy -docker build windows_32_static
	mkdir -p deploy/win32
	cp deploy/windows/editcp.exe deploy/win32

deploy/win64/editcp64.exe: $(SOURCES)
	go mod tidy
	go mod vendor
	qtdeploy -docker build windows_64_static
	mkdir -p deploy/win64
	cp deploy/windows/editcp.exe deploy/win64/editcp64.exe

macOS: clean darwin
darwin: FORCE
	@echo "Under macOS/darwin, it is ok to ignore the sed warnings about extra characters"
	go mod tidy
	go mod vendor
	qtdeploy build desktop

docker_usb_windows:
	docker rmi -f therecipe/qt:windows_32_static >/dev/null 2>&1
	docker pull therecipe/qt:windows_32_static
	cd ../docker/windows32-with-usb && \
		docker build -t therecipe/qt:windows_32_static_usb .
	docker rmi -f therecipe/qt:windows_32_static
	docker tag therecipe/qt:windows_32_static_usb therecipe/qt:windows_32_static

docker_usb_windows64:
	docker rmi -f therecipe/qt:windows_64_static >/dev/null 2>&1
	docker pull therecipe/qt:windows_64_static
	cd ../docker/windows64-with-usb && \
		docker build -t therecipe/qt:windows_64_static_usb .
	docker rmi -f therecipe/qt:windows_64_static
	docker tag therecipe/qt:windows_64_static_usb therecipe/qt:windows_64_static

docker_usb_linux:
	docker rmi -f therecipe/qt:linux >/dev/null 2>&1
	docker pull therecipe/qt:linux
	cd ../docker/linux-with-usb && \
		docker build -t therecipe/qt:linux_usb .
	docker rmi -f therecipe/qt:linux
	docker tag therecipe/qt:linux_usb therecipe/qt:linux

docker_usb: docker_usb_linux docker_usb_windows docker_usb_windows64

FORCE:

editcp-changelog: editcp-changelog.txt

editcp-changelog.txt: FORCE
	sh generateChangelog >editcp-changelog.txt

clean:
	rm -rf deploy/*

clobber: clean
	rm -rf editcp-*
