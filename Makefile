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

windows: deploy/windows/editcp.exe

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
	rsync deploy/windows/editcp.exe farnsworth.org:editcp-$(VERSION).exe

clean:
	rm -rf editcp editcp-$(VERSION)

clobber: clean
	rm -f deploy/linux/editcp deploy/linux/editcp.sh \
		deploy/windows/editcp.exe
