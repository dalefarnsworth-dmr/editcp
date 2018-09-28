SHELL = /bin/sh

.PHONY: default linux windows clean clobber upload install docker-usb tag

EDITCP_SRC = *.go
RADIO_SRC = ../dmrRadio/*.go
UI_SRC = ../ui/*.go
CODEPLUG_SRC = ../codeplug/*.go
DFU_SRC = ../dfu/*.go
STDFU_SRC = ../stdfu/*.go
USERDB_SRC = ../userdb/*.go
SOURCES = $(EDITCP_SRC) $(UI_SRC) $(CODEPLUG_SRC) $(DFU_SRC) $(STDFU_SRC) $(USERDB_SRC)
RADIO_SRCS =  $(RADIO_SRC) $(CODEPLUG_SRC) $(DFU_SRC) $(STDFU_SRC) $(USERDB_SRC)
VERSION = $(shell sed -n '/version =/{s/^[^"]*"//;s/".*//p;q}' <version.go)

default: linux dmrRadio

linux: deploy/linux/editcp deploy/linux/editcp.sh deploy/linux/install deploy/linux/99-md380.rules ../dmrRadio/dmrRadio

deploy/linux/editcp: $(SOURCES)
	qtdeploy -docker build
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

windows: editcp-$(VERSION)-installer.exe dmrRadio-$(VERSION)-installer.exe

editcp-$(VERSION)-installer.exe: deploy/win32/editcp.exe editcp.nsi dll/*.dll
	makensis -DVERSION=$(VERSION) editcp.nsi

dmrRadio-$(VERSION)-installer.exe: ../dmrRadio/dmrRadio.exe editcp.nsi dll/*.dll
	makensis -DVERSION=$(VERSION) dmrRadio.nsi

deploy/win32/editcp.exe: $(SOURCES)
	qtdeploy -docker build windows_32_static
	mkdir -p deploy/win32
	cp deploy/windows/editcp.exe deploy/win32

docker-usb:
	docker rmi -f therecipe/qt:windows_32_static
	cd ../docker/windows32-with-usb && \
		docker build -t therecipe/qt:windows_32_static .
	docker rmi -f therecipe/qt:linux
	cd ../docker/linux-with-usb && \
		docker build -t therecipe/qt:linux .

../dmrRadio/dmrRadio: $(RADIO_SRCS)
	cd ../dmrRadio && go build

../dmrRadio/dmrRadio.exe: $(RADIO_SRCS)
	cd ../dmrRadio && GOOS=windows GOARCH=386 go build

dmrRadio: dmrRadio-$(VERSION).tar.xz

dmrRadio-$(VERSION).tar.xz: ../dmrRadio/dmrRadio
	cd ../dmrRadio && go build
	rm -rf dmrRadio-$(VERSION)
	mkdir -p dmrRadio-$(VERSION)
	cp -al ../dmrRadio/dmrRadio dmrRadio-$(VERSION)
	tar cJf dmrRadio-$(VERSION).tar.xz dmrRadio-$(VERSION)
	rm -rf dmrRadio-$(VERSION)

FORCE:

changelog.txt: FORCE
	sh generateChangelog >changelog.txt

clean:
	rm -rf ../dmrRadio/dmrRadio

clobber: clean
	rm -rf editcp-* deploy/* dmrRadio-*

# The targets below are probably only useful for me. -Dale Farnsworth

upload-no-push: changelog.txt editcp-$(VERSION).tar.xz editcp-$(VERSION)-installer.exe dmrRadio-$(VERSION).tar.xz dmrRadio-$(VERSION)-installer.exe
	rsync editcp-$(VERSION).tar.xz farnsworth.org:
	rsync editcp-$(VERSION)-installer.exe farnsworth.org:
	rsync dmrRadio-$(VERSION).tar.xz farnsworth.org:
	rsync dmrRadio-$(VERSION)-installer.exe farnsworth.org:
	rsync changelog.txt farnsworth.org:

upload: tag upload-no-push
	git push
	git push --tags
	rsh farnsworth.org uploaded

tag:
	git log --no-decorate -1 --oneline | grep -q -i "update.*version"
	git tag -f -s -m "editcp v$(VERSION)" v$(VERSION)
