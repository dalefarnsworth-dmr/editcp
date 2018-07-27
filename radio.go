// Copyright 2017-2018 Dale Farnsworth. All rights reserved.

// Dale Farnsworth
// 1007 W Mendoza Ave
// Mesa, AZ  85210
// USA
//
// dale@farnsworth.org

// This file is part of Editcp.
//
// Editcp is free software: you can redistribute it and/or modify
// it under the terms of version 3 of the GNU General Public License
// as published by the Free Software Foundation.
//
// Editcp is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Editcp.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dalefarnsworth/codeplug/codeplug"
	"github.com/dalefarnsworth/codeplug/dfu"
	"github.com/dalefarnsworth/codeplug/ui"
	"github.com/dalefarnsworth/codeplug/userdb"
	"github.com/therecipe/qt/core"
)

type modelUrl struct {
	model string
	url   string
}

func (edt *editor) addRadioMenu(menu *ui.Menu) {
	cp := edt.codeplug
	mb := edt.mainWindow.MenuBar()
	menu = mb.AddMenu("Radio")

	menu.AddAction("Read codeplug from radio", func() {
		err := codeplug.RadioExists()
		if err != nil {
			title := fmt.Sprintf("Read codeplug from radio failed")
			ui.ErrorPopup(title, err.Error())
			return
		}

		edt := newEditor(edt.app, codeplug.FileTypeNew, "")
		if edt == nil || edt.codeplug == nil {
			return
		}

		msgs := []string{
			"Preparing to read codeplug from radio...",
			"Reading codeplug from radio...",
		}
		msgIndex := 0
		pd := ui.NewProgressDialog(msgs[msgIndex])
		err = edt.codeplug.ReadRadio(func(cur int) bool {
			if cur == codeplug.MinProgress {
				pd.SetLabelText(msgs[msgIndex])
				msgIndex++
			}

			pd.SetRange(codeplug.MinProgress, codeplug.MaxProgress)
			pd.SetValue(cur)
			if pd.WasCanceled() {
				return false
			}
			return true
		})
		if err != nil {
			pd.Close()
			title := fmt.Sprintf("Read codeplug from radio failed")
			ui.ErrorPopup(title, err.Error())
			edt.FreeCodeplug()
		}
	})

	menu.AddAction("Write codeplug to radio", func() {
		title := "Write codeplug to radio"
		model := cp.Model()
		freq := cp.FrequencyRange()
		msg := fmt.Sprintf("Write %s %s codeplug to radio?\n", model, freq)
		if ui.YesNoPopup(title, msg) != ui.PopupYes {
			return
		}

		msgs := []string{
			"Preparing to write codeplug to radio...",
			"Writing codeplug to radio...",
		}
		msgIndex := 0

		pd := ui.NewProgressDialog(msgs[msgIndex])
		err := cp.WriteRadio(func(cur int) bool {
			if cur == codeplug.MinProgress {
				pd.SetLabelText(msgs[msgIndex])
				msgIndex++
			}
			pd.SetRange(codeplug.MinProgress, codeplug.MaxProgress)
			pd.SetValue(cur)
			if pd.WasCanceled() {
				return false
			}
			return true
		})
		if err != nil {
			pd.Close()
			title := fmt.Sprintf("Write codeplug to radio failed: %s", err.Error())
			ui.ErrorPopup(title, err.Error())
		}
	}).SetEnabled(cp != nil && cp.Loaded())

	md380toolsMenu := menu.AddMenu("md380tools...")

	md380toolsMenu.AddAction("Write user database to radio...", func() {
		title := "Write user database to radio"
		cancel, download := userdbDialog(title)
		if cancel {
			return
		}

		locType := core.QStandardPaths__CacheLocation
		cacheDir := core.QStandardPaths_WritableLocation(locType)
		tmpFilename := filepath.Join(cacheDir, "users.tmp")

		msgs := []string{
			"Downloading user database from web sites...",
			"Erasing the radio's flash memory for user database...",
			"Writing user database to radio...",
		}
		msgIndex := 0
		if !download {
			msgIndex = 1
		}

		filename := userdbFilename()
		os.MkdirAll(filepath.Dir(filename), os.ModeDir|0755)

		pd := ui.NewProgressDialog(msgs[msgIndex])

		if download {
			err := userdb.WriteMD380ToolsFile(tmpFilename, func(cur int) bool {
				if cur == userdb.MinProgress {
					pd.SetLabelText(msgs[msgIndex])
					msgIndex++
				}
				pd.SetRange(userdb.MinProgress, userdb.MaxProgress)
				pd.SetValue(cur)
				if pd.WasCanceled() {
					return false
				}
				return true
			})
			if err != nil {
				os.Remove(tmpFilename)
				pd.Close()
				title := fmt.Sprintf("Download of user database failed")
				ui.ErrorPopup(title, err.Error())
				return
			}

			os.Rename(tmpFilename, filename)
		}
		dfu, err := dfu.New(func(cur int) bool {
			if cur == dfu.MinProgress {
				pd.SetLabelText(msgs[msgIndex])
				msgIndex++
			}
			pd.SetRange(dfu.MinProgress, dfu.MaxProgress)
			pd.SetValue(cur)
			if pd.WasCanceled() {
				return false
			}
			return true

		})
		if err == nil {
			defer dfu.Close()
			err = dfu.WriteUsers(filename)
		}
		if err != nil {
			pd.Close()
			title := fmt.Sprintf("write of user database failed: %s", err.Error())
			ui.ErrorPopup(title, err.Error())
		}
	})

	md380toolsMenu.AddAction("Write md380tools firmware to radio...", func() {
		path := "https://farnsworth.org/dale/md380tools/"
		nonGpsUrl := path + "firmware/D13.20.bin"
		gpsUrl := path + "firmware/S13.20.bin"

		modelUrls := []modelUrl{
			modelUrl{"MD-380 (D13.20)", nonGpsUrl},
			modelUrl{"MD-380G (S13.20)", gpsUrl},
			modelUrl{"MD-390 (D13.20)", nonGpsUrl},
			modelUrl{"MD-390G (S13.20)", gpsUrl},
			modelUrl{"RT3 (D13.20)", nonGpsUrl},
			modelUrl{"RT8 (S13.20)", gpsUrl},
		}

		title := "Write md380tools firmware to radio..."
		upgrade := true
		canceled, model, url := firmwareDialog(title, modelUrls, upgrade)
		if canceled {
			return
		}

		msgs := []string{
			fmt.Sprintf("Downloading md380tools %s firmware...\n%s", model, url),
			"Erasing the radio's firmware flash memory...",
			fmt.Sprintf("Writing md380tools %s firmware to radio...", model),
		}

		writeFirmware(url, msgs)
	})

	md380toolsMenu.AddAction("Write original firmware to radio...", func() {
		path := "https://farnsworth.org/dale/md380tools/"
		d003Url := path + "original_firmware/D003.020.bin"
		d013Url := path + "original_firmware/D013.020.bin"
		d013_34Url := path + "original_firmware/D013.034.bin"
		s013Url := path + "original_firmware/S013.020.bin"
		d14_04Url := path + "original_firmware/D014.004.bin"

		modelUrls := []modelUrl{
			modelUrl{"MD-380 old (D03.20)", d003Url},
			modelUrl{"MD-380 (D13.20)", d013Url},
			modelUrl{"MD-380 new (D13.34)", d013_34Url},
			modelUrl{"MD-380 newest (D14.04", d14_04Url},
			modelUrl{"MD-380G (S13.20)", s013Url},
			modelUrl{"MD-390 (D13.20)", d013Url},
			modelUrl{"MD-390G (S13.20)", s013Url},
			modelUrl{"RT3 (D03.20)", d003Url},
			modelUrl{"RT8 (S13.20)", s013Url},
		}

		title := "Write original firmware to radio..."
		upgrade := false
		canceled, model, url := firmwareDialog(title, modelUrls, upgrade)
		if canceled {
			return
		}

		msgs := []string{
			fmt.Sprintf("Downloading original %s firmware...\n%s", model, url),
			"Erasing the radio's firmware flash memory...",
			fmt.Sprintf("Writing original %s firmware to radio...", model),
		}

		writeFirmware(url, msgs)
	})
}

func writeFirmware(url string, msgs []string) {
	tmpFile, err := ioutil.TempFile("", "editcp")
	if err != nil {
		title := fmt.Sprintf("temporary file failed: %s", err.Error())
		ui.ErrorPopup(title, err.Error())
		return
	}

	filename := tmpFile.Name()
	defer os.Remove(filename)

	msgIndex := 0
	pd := ui.NewProgressDialog(msgs[msgIndex])

	df, err := dfu.New(func(cur int) bool {
		if cur == dfu.MinProgress {
			pd.SetLabelText(msgs[msgIndex])
			msgIndex++
		}
		pd.SetRange(dfu.MinProgress, dfu.MaxProgress)
		pd.SetValue(cur)
		if pd.WasCanceled() {
			return false
		}
		return true

	})
	if err != nil {
		pd.Close()
		title := "firmware write failed"
		ui.ErrorPopup(title, err.Error())
		return
	}
	defer df.Close()

	err = download(url, filename, func(cur int) bool {
		if cur == dfu.MinProgress {
			pd.SetLabelText(msgs[msgIndex])
			msgIndex++
		}
		pd.SetRange(userdb.MinProgress, userdb.MaxProgress)
		pd.SetValue(cur)
		if pd.WasCanceled() {
			return false
		}
		return true
	})
	if err != nil {
		pd.Close()
		title := "firmware write failed"
		ui.ErrorPopup(title, err.Error())
		return
	}

	err = df.WriteFirmware(filename)
	if err != nil {
		pd.Close()
		title := "write of new firmware failed"
		ui.ErrorPopup(title, err.Error())
		return
	}

	msg := "Turn radio off and back on again."
	ui.InfoPopup("Firmware write complete", msg)
}

func userdbFilename() string {
	locType := core.QStandardPaths__CacheLocation
	cacheDir := core.QStandardPaths_WritableLocation(locType)

	name := "usersDB.bin"

	return filepath.Join(cacheDir, name)
}

func userdbDialog(title string) (canceled, download bool) {
	loadSettings()

	usersFilename := userdbFilename()

	download = true
	if fileYounger(usersFilename, 12*time.Hour) {
		download = false
	}

	downloadCheckbox := ui.NewCheckboxWidget(download, func(checked bool) {
		download = checked
	})
	downloadCheckbox.SetEnabled(fileExists(usersFilename))

	dialog := ui.NewDialog(title)

	filenameBox := ui.NewHbox()
	filenameBox.AddLabel("   " + usersFilename)

	labelText := `
The users database contains DMR ID numbers and callsigns of all registered
users. It can only be be written to radios that have been upgraded to the
md380tools firmware.  See https://github.com/travisgoodspeed/md380tools.`

	dialog.AddLabel(labelText[1:])

	form := dialog.AddForm()
	form.AddRow("Download new users database file", downloadCheckbox)

	dialog.AddLabel("Filename:")
	dialog.AddExistingHbox(filenameBox)

	row := dialog.AddHbox()
	cancelButton := ui.NewButtonWidget("Cancel", func() {
		dialog.Reject()
	})
	row.AddWidget(cancelButton)

	saveButton := ui.NewButtonWidget("Write", func() {
		dialog.Accept()
	})
	row.AddWidget(saveButton)

	saved := dialog.Exec()
	return !saved, download
}

func firmwareDialog(title string, modelUrls []modelUrl, upgrade bool) (canceled bool, model, url string) {

	models := make([]string, len(modelUrls))
	for i, modelUrl := range modelUrls {
		models[i] = modelUrl.model
	}

	model = models[0]
	modelCombobox := ui.NewComboboxWidget(model, models, func(selected string) {
		model = selected
	})

	dialog := ui.NewDialog(title)

	var labelText string
	if upgrade {
		labelText += `
The md380tools firmware only works on MD380, MD380, RT3, and RT8 radios.`
	}

	labelText += `

Before continuing, enable bootloader mode:
	1. Insert a cable into USB.
	2. Connect the cable to the radio.
	3. Power-on the radio by turning volume knob, while holding down
	   the PTT button and the button above PTT.

While in bootloader mode, the LED will flash green and red.`

	if !upgrade {
		labelText += `

Hint: If the display becomes flipped on the md380, try another
md380 variant.`
	}

	dialog.AddLabel(labelText[1:])

	groupBox := dialog.AddGroupbox("Select Radio Model")
	form := groupBox.AddForm()
	form.AddRow("Radio model", modelCombobox)

	row := dialog.AddHbox()

	cancelButton := ui.NewButtonWidget("Cancel", func() {
		dialog.Reject()
	})
	row.AddWidget(cancelButton)

	saveButton := ui.NewButtonWidget("Update Firmware", func() {
		dialog.Accept()
	})
	row.AddWidget(saveButton)

	saved := dialog.Exec()

	for _, modelUrl := range modelUrls {
		if modelUrl.model == model {
			url = modelUrl.url
			break
		}
	}

	return !saved, model, url
}

var timeoutSeconds = 20

var tr = &http.Transport{
	TLSHandshakeTimeout:   time.Duration(timeoutSeconds) * time.Second,
	ResponseHeaderTimeout: time.Duration(timeoutSeconds) * time.Second,
}

var client = &http.Client{
	Transport: tr,
	Timeout:   time.Duration(timeoutSeconds) * time.Second,
}

type downloader struct {
	url               string
	filename          string
	progressCallback  func(progressCounter int) bool
	progressFunc      func() error
	progressIncrement int
	progressCounter   int
}

func newDownloader() *downloader {
	d := &downloader{
		progressFunc: func() error { return nil },
	}

	return d
}

func (d *downloader) setMaxProgressCount(max int) {
	d.progressFunc = func() error { return nil }
	if d.progressCallback != nil {
		d.progressIncrement = MaxProgress / max
		d.progressCounter = 0
		d.progressFunc = func() error {
			d.progressCounter += d.progressIncrement
			curProgress := d.progressCounter
			if curProgress > MaxProgress {
				curProgress = MaxProgress
			}

			if !d.progressCallback(d.progressCounter) {
				return errors.New("")
			}

			return nil
		}
		d.progressCallback(d.progressCounter)
	}
}

func (d *downloader) finalProgress() {
	//fmt.Fprintf(os.Stderr, "\nprogressMax %d\n", d.progressCounter/d.progressIncrement)
	if d.progressCallback != nil {
		d.progressCallback(MaxProgress)
	}
}

const (
	MinProgress = 0
	MaxProgress = 1000000
)

func download(url, filename string, progress func(cur int) bool) error {
	d := newDownloader()
	d.url = url
	d.filename = filename
	d.progressCallback = progress
	return d.download()
}

func (d *downloader) download() (err error) {
	file, err := os.Create(d.filename)
	if err != nil {
		return wrapError("download", err)
	}
	defer func() {
		fErr := file.Close()
		if err == nil {
			err = fErr
		}
		return
	}()

	resp, err := client.Get(d.url)
	if err != nil {
		return wrapError("download", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return wrapError("download", errors.New(resp.Status))
	}
	length := resp.ContentLength
	if length < 0 {
		length = 1024 * 1024
	}

	bufSize := 16 * 1024

	d.setMaxProgressCount(int(length) / bufSize)

	buf := make([]byte, bufSize)
	for {
		err := d.progressFunc()
		if err != nil {
			return wrapError("download", err)
		}

		n, err := resp.Body.Read(buf)
		if n == 0 && err != nil {
			if err == io.EOF {
				break
			}
			return wrapError("download", err)
		}

		n, err = file.Write(buf)
		if err != nil {
			return wrapError("download", err)
		}
	}

	d.finalProgress()

	return nil
}

func wrapError(prefix string, err error) error {
	if err.Error() == "" {
		return err
	}
	return fmt.Errorf("%s: %s", prefix, err.Error())
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func fileYounger(filename string, duration time.Duration) bool {
	fileInfo, err := os.Stat(filename)
	return err == nil && time.Since(fileInfo.ModTime()) < duration
}
