SHELL = /bin/sh

.PHONY: default clean clobber build install run tar upload windows

EDITCP_SOURCES = *.go
UI_SOURCES = ../ui/*.go
CODEPLUG_SOURCES = ../codeplug/*.go
SOURCES = $(EDITCP_SOURCES) $(UI_SOURCES) $(CODEPLUG_SOURCES)
VERSION = $(shell sed -n '/version =/{s/^[^"]*"//;s/".*//p;q}' <version.go)

default: deploy/linux/editcp.sh

deploy/linux/editcp.sh: deploy/linux/editcp editcp.sh deploy/linux/install
	cp editcp.sh deploy/linux/editcp.sh
	@cd deploy/linux && ./install .

deploy/linux/install: install.sh
	cp install.sh deploy/linux/install

deploy/linux/editcp: $(SOURCES)
	qtdeploy -docker build

windows: editcp-$(VERSION).msi

editcp-$(VERSION).msi: deploy/windows/editcp.exe editcp.wxs
	sed 's/VERSION/$(VERSION)/g' editcp.wxs >editcp-$(VERSION).wxs
	wixl --arch x64 -o editcp-$(VERSION).msi editcp-$(VERSION).wxs
	rm editcp-$(VERSION).wxs

deploy/windows/editcp.exe: $(SOURCES)
	qtdeploy -docker build windows_64_static

install: deploy/linux/editcp.sh
	@mkdir -p deploy/linux/bin
	@cd deploy/linux && ./install 

build:
	go build

run: deploy/linux/editcp.sh
	deploy/linux/editcp.sh

tar: deploy/linux/editcp-$(VERSION).tar.xz

editcp-$(VERSION).tar.xz: deploy/linux/editcp.sh
	rm -rf editcp-$(VERSION)
	mkdir -p editcp-$(VERSION)
	cp -al deploy/linux/* editcp-$(VERSION)
	tar cJf editcp-$(VERSION).tar.xz editcp-$(VERSION)
	rm -rf editcp-$(VERSION)

upload: editcp-$(VERSION).tar.xz windows
	rsync editcp-$(VERSION).tar.xz farnsworth.org:
	rsync editcp-$(VERSION).msi farnsworth.org:

tag:
	git tag -s -m "editcp v$(VERSION)" v$(VERSION)

clean:
	rm -rf editcp editcp-$(VERSION) editcp-*.wxs

clobber: clean
	rm -f deploy editcp-*.wxs
