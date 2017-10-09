SHELL = /bin/sh

.PHONY: default clean clobber build install run tar upload

EDITCP_SOURCES = *.go
UI_SOURCES = ../ui/*.go
CODEPLUG_SOURCES = ../codeplug/*.go
SOURCES = $(EDITCP_SOURCES) $(UI_SOURCES) $(CODEPLUG_SOURCES)
DEPLOYDIR = deploy/linux
VERSION = $(shell sed -n '/^const version =/{s/^[^"]*"//;s/".*//p;q}' <editcp.go)

default: $(DEPLOYDIR)/editcp.sh

$(DEPLOYDIR)/editcp.sh: $(DEPLOYDIR)/editcp editcp.sh $(DEPLOYDIR)/install
	cp editcp.sh $(DEPLOYDIR)/editcp.sh
	@cd $(DEPLOYDIR) && ./install .

$(DEPLOYDIR)/install: install.sh
	cp install.sh $(DEPLOYDIR)/install

$(DEPLOYDIR)/editcp: $(SOURCES)
	qtdeploy -docker build

install: $(DEPLOYDIR)/editcp.sh
	@mkdir -p $(DEPLOYDIR)/bin
	@cd $(DEPLOYDIR) && ./install 

build:
	go build

run: $(DEPLOYDIR)/editcp.sh
	$(DEPLOYDIR)/editcp.sh

tar: $(DEPLOYDIR)/editcp-$(VERSION).tar.xz

editcp-$(VERSION).tar.xz: $(DEPLOYDIR)/editcp.sh
	rm -rf editcp-$(VERSION)
	mkdir -p editcp-$(VERSION)
	cp -al deploy/linux/* editcp-$(VERSION)
	tar cJf editcp-$(VERSION).tar.xz editcp-$(VERSION)
	rm -rf editcp-$(VERSION)

upload: editcp-$(VERSION).tar.xz
	rsync editcp-$(VERSION).tar.xz farnsworth.org:

clean:
	rm -rf editcp editcp-$(VERSION)

clobber: clean
	rm -f $(DEPLOYDIR)/editcp $(DEPLOYDIR)/editcp.sh
