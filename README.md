## Codeplug editor for the MD380/MD390/MD40 DMR Radios

## Introduction
This program is similar in purpose to the [MD380 CPS program](
http://www.tyt888.com/?mod=download) provided by TYT Electronics
Technology Co., LTD.  It provides several features that CPS lacks,
while not implementing all features of CPS.
I wrote `editcp` because I wanted to be able to edit codeplugs in Linux.

This codeplug editor works for the Tytera MD380 and MD390 as well as the
Alinco DJ-MD40 radios.  Support for additional radio models is likely,
but is not scheduled at this time.

### Features
* `Editcp` permits the editing of General Settings, Channels, Contacts, Zones,
Group Lists, and Scan Lists.
* It supports reordering list items via drag-and-drop.
* Multiple codeplugs may be opened simultaneously and
items may be copied from one code plug to another via drag-and-drop.
* `Editcp` provides unlimited undo/redo.
* `Editcp` performs extensive input validation and codeplug entry validation.
* Codeplug information may be exported to and imported from human readable
text files.
* `Editcp` can edit .rdt files as well as the .bin files produced
by [md380tools](https://github.com/travisgoodspeed/md380tools).

## Building from Source
`Editcp` development has been done on Linux (specifically Ubuntu 17.04),
so that is the recommended platform for building from source.

1. `Editcp` is written in [go](https://golang.org/).  You must download
and install [go](https://golang.org/dl/) version 1.8 or later.

2. Install [git](https://git-scm.com/). On Debian, Ubuntu, and other
Debian-derived systems that may be done by:

```bash
$ sudo apt-get install git
```

3. `Editcp` uses the QT GUI library. You'll need to install
the [Qt binding for Go](https://github.com/therecipe/qt). I recommend
the docker installation described at
https://github.com/therecipe/qt/wiki/Deploying-Linux-to-Linux
and
https://github.com/therecipe/qt/wiki/Deploying-Linux-to-Windows-32-bit-Static

4. `Editcp` uses the libusb-1.0-0-dev package. You'll need to install it.
Also, if you're using the docker qt installation, you'll need to install
libusb-1.0-0-dev in the docker images it uses. This can be done by running
```bash
$ sudo docker tag therecipe/qt:windows_32_static therecipe/qt:windows_32_static-orig
$ sudo docker tag therecipe/qt:linux therecipe/qt:linux-orig
$ sudo make docker_usb
```
If you are building under macOS, then there is no need to use docker,
see
[Deploying Linux to macOS](https://github.com/therecipe/qt/wiki/Deploying-Linux-to-macOS).

5. Get the source code:
```bash
$ go get -d github.com/dalefarnsworth-dmr/editcp/...
$ go get github.com/google/gousb
$ go get github.com/tealeg/xlsx
```

6. Change to the `editcp` source directory:
```bash
$ cd $GOPATH/src/github.com/dalefarnsworth-dmr/editcp
```

7. Build `editcp`:
```bash
$ make
```
Under macOS, instead of running `make`, run `qtdeploy build desktop`

8. Install `editcp`:
```bash
$ make install
```
You will be prompted for a directory name where a symbolic link to
the `editcp` executable will be placed. If you don't have write permissions
for that directory, you will need to run this command as root.
Under macOS, skip running `make install`

9. You man now run `editcp`, optionally passing the name of a codeplug file.
```bash
$ editcp
```
or
```bash
$ editcp file.rdt
```
Under macOS, the app is at `deploy/darwin/editcp.app`.  To run from
the Terminal so as to see stdout and stderr, run
`./deploy/darwin/editcp.app/Contents/MacOS/editcp`

## Installing Pre-built Executables
Instructions for downloading pre-built executables for Windows and Linux are
available at https://www.farnsworth.org/dale/editcp.

## Disclaimer
`Editcp` has only been used by a small number of people at present. While
no problems have been observed in radios after loading codeplugs edited by
`editcp`, I can't guarantee that such will never occur. Use `editcp` at
your own risk.

## Contributing
Contributions to `editcp` are welcome. If you've fixed a bug or implemented
a cool new feature that you would like to share, please feel free to open
a pull request here.

## Author
Dale Farnsworth

<dale@farnsworth.org>

IRC: freenode channel: #md380, user: dfarnsworth
