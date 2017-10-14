SHELL = /bin/sh

.PHONY: default linux windows clean clobber upload install tag

EDITCP_SOURCES = *.go
UI_SOURCES = ../ui/*.go
CODEPLUG_SOURCES = ../codeplug/*.go
SOURCES = $(EDITCP_SOURCES) $(UI_SOURCES) $(CODEPLUG_SOURCES)
VERSION = $(shell sed -n '/version =/{s/^[^"]*"//;s/".*//p;q}' <version.go)

default: linux

linux: deploy/linux/editcp deploy/linux/editcp.sh deploy/linux/install

deploy/linux/editcp: $(SOURCES)
	qtdeploy -docker build

.PHONY: deploy/linux/editcp.sh	# Force, in case it's overwritten by install
deploy/linux/editcp.sh: editcp.sh
	cp editcp.sh deploy/linux/editcp.sh

deploy/linux/install: install.sh deploy/linux/editcp
	cp install.sh deploy/linux/install

editcp-$(VERSION).tar.xz: deploy/linux/editcp.sh
	rm -rf editcp-$(VERSION)
	mkdir -p editcp-$(VERSION)
	cp -al deploy/linux/* editcp-$(VERSION)
	tar cJf editcp-$(VERSION).tar.xz editcp-$(VERSION)
	rm -rf editcp-$(VERSION)

install: linux
	cd deploy/linux && ./install .

windows: editcp-$(VERSION).msi

editcp-$(VERSION).msi: deploy/windows/editcp.exe editcp.wxs
	sed 's/VERSION/$(VERSION)/g' editcp.wxs >editcp-$(VERSION).wxs
	wixl --arch x64 -o editcp-$(VERSION).msi editcp-$(VERSION).wxs
	rm editcp-$(VERSION).wxs

deploy/windows/editcp.exe: $(SOURCES)
	qtdeploy -docker build windows_64_static

clean:
	rm -rf editcp editcp-$(VERSION) editcp-*.wxs

clobber: clean
	rm -rf deploy editcp-*.wxs

# The targets below are probably only useful for me. -Dale Farnsworth

upload: linux windows
	rsync editcp-$(VERSION).tar.xz farnsworth.org:
	rsync editcp-$(VERSION).msi farnsworth.org:

tag:
	git tag -s -m "editcp v$(VERSION)" v$(VERSION)
