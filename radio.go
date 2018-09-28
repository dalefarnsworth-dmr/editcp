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

type modelURL struct {
	model string
	url   string
}

func writeMD380toolsUsers() {
	title := "Write user database to radio"
	text := `
The users database contains DMR ID numbers and callsigns of all registered
users. It can only be be written to radios that have been upgraded to the
md380tools firmware.  See https://github.com/travisgoodspeed/md380tools.

WARNING: Corruption may occur if a signal is received while writing to the
radio.  The radio should be tuned to an unprogrammed (or at least quiet)
channel while writing the new user database.`

	cancel, download := userdbDialog(title, text)
	if cancel {
		return
	}

	locType := core.QStandardPaths__CacheLocation
	cacheDir := core.QStandardPaths_WritableLocation(locType)
	tmpFilename := filepath.Join(cacheDir, "users.tmp")

	msgs := []string{
		"Downloading user database from web sites...",
		"Erasing the radio's user database...",
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
		db := userdb.New()
		err := db.WriteMD380ToolsFile(tmpFilename, func(cur int) error {
			if cur == userdb.MinProgress {
				pd.SetLabelText(msgs[msgIndex])
				msgIndex++
			}
			pd.SetRange(userdb.MinProgress, userdb.MaxProgress)
			pd.SetValue(cur)
			if pd.WasCanceled() {
				return errors.New("cancelled")
			}
			return nil
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
	df, err := dfu.New(func(cur int) error {
		if cur == dfu.MinProgress {
			pd.SetLabelText(msgs[msgIndex])
			msgIndex++
		}
		pd.SetRange(dfu.MinProgress, dfu.MaxProgress)
		pd.SetValue(cur)
		if pd.WasCanceled() {
			return errors.New("cancelled")
		}
		return nil

	})
	if err == nil {
		defer df.Close()
		err = df.WriteUsers(filename)
	}
	if err != nil {
		pd.Close()
		title := fmt.Sprintf("write of user database failed: %s", err.Error())
		ui.ErrorPopup(title, err.Error())
	}
}

func writeExpandedUsers(title, text string) {
	cancel, download := userdbDialog(title, text)
	if cancel {
		return
	}

	locType := core.QStandardPaths__CacheLocation
	cacheDir := core.QStandardPaths_WritableLocation(locType)
	tmpFilename := filepath.Join(cacheDir, "users.tmp")

	msgs := []string{
		"Downloading user database from web sites...",
		"Preparing to write user database to radio...",
		"Erasing the radio's user database...",
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
		db := userdb.New()
		err := db.WriteMD380ToolsFile(tmpFilename, func(cur int) error {
			if cur == userdb.MinProgress {
				pd.SetLabelText(msgs[msgIndex])
				msgIndex++
			}
			pd.SetRange(userdb.MinProgress, userdb.MaxProgress)
			pd.SetValue(cur)
			if pd.WasCanceled() {
				return errors.New("cancelled")
			}
			return nil
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
	df, err := dfu.New(func(cur int) error {
		if cur == dfu.MinProgress {
			pd.SetLabelText(msgs[msgIndex])
			msgIndex++
		}
		pd.SetRange(dfu.MinProgress, dfu.MaxProgress)
		pd.SetValue(cur)
		if pd.WasCanceled() {
			return errors.New("cancelled")
		}
		return nil

	})
	if err == nil {
		defer df.Close()
		file, err := os.Open(filename)
		if err == nil {
			defer file.Close()
			users := dfu.ParseUsers(file)
			err = df.WriteMD2017Users(users)
		}
	}
	if err != nil {
		pd.Close()
		title := fmt.Sprintf("write of user database failed: %s", err.Error())
		ui.ErrorPopup(title, err.Error())
	}
}

func writeMD2017Users() {
	title := "Write user database to radio"
	text := `
The users database contains DMR ID numbers and callsigns of all registered
users. It can only be be written to MD-2017 radios.

WARNING: This only works on MD-2017 radios with the "CSV" firmware versions.`

	writeExpandedUsers(title, text)
}

func writeUV380Users() {
	title := "Write user database to radio"
	text := `
The users database contains DMR ID numbers and callsigns of all registered
users. It can only be be written to MD-UV380 radios.

WARNING: This only works on MD-UV380 radios with the "CSV" firmware versions.`

	writeExpandedUsers(title, text)
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

		cp := edt.codeplug

		msgs := []string{
			"Preparing to read codeplug from radio...",
			"Reading codeplug from radio...",
		}
		msgIndex := 0
		pd := ui.NewProgressDialog(msgs[msgIndex])
		err = cp.ReadRadio(func(cur int) error {
			if cur == codeplug.MinProgress {
				pd.SetLabelText(msgs[msgIndex])
				msgIndex++
			}

			pd.SetRange(codeplug.MinProgress, codeplug.MaxProgress)
			pd.SetValue(cur)
			if pd.WasCanceled() {
				return errors.New("cancelled")
			}
			return nil
		})
		if err != nil {
			pd.Close()
			title := fmt.Sprintf("Read codeplug from radio failed")
			ui.ErrorPopup(title, err.Error())
			edt.FreeCodeplug()
		}

		if !cp.Valid() {
			fmtStr := `
%d records with invalid field values were found in the codeplug.

Select "Menu->Edit->Show Invalid Fields" to view them.`
			msg := fmt.Sprintf(fmtStr, len(cp.Warnings()))
			ui.InfoPopup("codeplug warning", msg)
		}
		edt.updateMenuBar()
	})

	menu.AddAction("Write codeplug to radio", func() {
		valid := cp.Valid()
		edt.updateMenuBar()
		if !valid {
			fmtStr := `
%d records with invalid field values were found in the codeplug.

Click on Cancel and then select "Menu->Edit->Show Invalid Fields" to view them.

Or, click on Ignore to continue writing to the radio.`
			msg := fmt.Sprintf(fmtStr, len(cp.Warnings()))
			title := "write warning"
			rv := ui.WarningPopup(title, msg)
			if rv != ui.PopupIgnore {
				return
			}
		}

		title := "Write codeplug to radio"
		model := cp.Model()
		freq := cp.FrequencyRange()
		warn := `

WARNING: Corruption may occur if a signal is received
while writing to the radio.  The radio should be tuned
to an unprogrammed (or at least quiet) channel while
writing the new codeplug.`
		msg := fmt.Sprintf("%s\n\nWrite %s %s codeplug to radio?\n", warn, model, freq)
		if ui.YesNoPopup(title, msg) != ui.PopupYes {
			return
		}

		msgs := []string{
			"Preparing to write codeplug to radio...",
			"Erasing the radio's codeplug...",
			"Writing codeplug to radio...",
		}
		msgIndex := 0

		pd := ui.NewProgressDialog(msgs[msgIndex])
		err := cp.WriteRadio(func(cur int) error {
			if cur == codeplug.MinProgress {
				pd.SetLabelText(msgs[msgIndex])
				msgIndex++
			}
			pd.SetRange(codeplug.MinProgress, codeplug.MaxProgress)
			pd.SetValue(cur)
			if pd.WasCanceled() {
				return errors.New("cancelled")
			}
			return nil
		})
		if err != nil {
			pd.Close()
			title := fmt.Sprintf("Write codeplug to radio failed: %s", err.Error())
			ui.ErrorPopup(title, err.Error())
		}
	}).SetEnabled(cp != nil && cp.Loaded())

	menu.AddSeparator()

	menu.AddAction("Write factory firmware to radio...", func() {
		path := "https://farnsworth.org/dale/md380tools/"
		d003URL := path + "original_firmware/D003.020.bin"
		d013URL := path + "original_firmware/D013.020.bin"
		d013_34URL := path + "original_firmware/D013.034.bin"
		s013URL := path + "original_firmware/S013.020.bin"
		d14_04URL := path + "original_firmware/D014.004.bin"

		modelURLs := []modelURL{
			modelURL{"MD-380 old (D03.20)", d003URL},
			modelURL{"MD-380 (D13.20)", d013URL},
			modelURL{"MD-380 new (D13.34)", d013_34URL},
			modelURL{"MD-380 newest (D14.04", d14_04URL},
			modelURL{"MD-380G (S13.20)", s013URL},
			modelURL{"MD-390 (D13.20)", d013URL},
			modelURL{"MD-390G (S13.20)", s013URL},
			modelURL{"RT3 (D03.20)", d003URL},
			modelURL{"RT8 (S13.20)", s013URL},
		}

		title := "Write factory firmware to radio..."
		upgrade := false
		canceled, model, url := firmwareDialog(title, modelURLs, upgrade)
		if canceled {
			return
		}

		msgs := []string{
			fmt.Sprintf("Downloading original %s firmware...\n%s", model, url),
			"Erasing the radio's firmware...",
			fmt.Sprintf("Writing factory %s firmware to radio...", model),
		}

		writeFirmware(url, msgs)
	})

	menu.AddSeparator()
	writeUsersMenu := menu.AddMenu("Write user database to radio...")

	writeUsersMenu.AddAction("Write md380tools user database to radio...", writeMD380toolsUsers)
	writeUsersMenu.AddAction("Write MD2017 user database to radio...", writeMD2017Users)
	writeUsersMenu.AddAction("Write MD-UV380 user database to radio...", writeUV380Users)

	menu.AddSeparator()

	md380toolsMenu := menu.AddMenu("md380tools...")

	md380toolsMenu.AddAction("Write user database to radio...", writeMD380toolsUsers)

	md380toolsMenu.AddAction("Write md380tools firmware to radio...", func() {
		path := "https://farnsworth.org/dale/md380tools/firmware/"
		nonGpsURL := path + "D13.20.bin"
		gpsURL := path + "S13.20.bin"

		modelURLs := []modelURL{
			modelURL{"MD-380 (D13.20)", nonGpsURL},
			modelURL{"MD-380G (S13.20)", gpsURL},
			modelURL{"MD-390 (D13.20)", nonGpsURL},
			modelURL{"MD-390G (S13.20)", gpsURL},
			modelURL{"RT3 (D13.20)", nonGpsURL},
			modelURL{"RT8 (S13.20)", gpsURL},
		}

		title := "Write md380tools firmware to radio..."
		upgrade := true
		canceled, model, url := firmwareDialog(title, modelURLs, upgrade)
		if canceled {
			return
		}

		msgs := []string{
			fmt.Sprintf("Downloading md380tools %s firmware...\n%s", model, url),
			"Erasing the radio's firmware...",
			fmt.Sprintf("Writing md380tools %s firmware to radio...", model),
		}

		writeFirmware(url, msgs)
	})
	md380toolsMenu.AddAction("Write KD4Z md380tools firmware to radio...", func() {
		path := "https://farnsworth.org/dale/md380tools/kd4z/"
		nonGpsURL := path + "firmware-noGPS.bin"
		gpsURL := path + "firmware-GPS.bin"

		modelURLs := []modelURL{
			modelURL{"MD-380 (D13.20)", nonGpsURL},
			modelURL{"MD-380G (S13.20)", gpsURL},
			modelURL{"MD-390 (D13.20)", nonGpsURL},
			modelURL{"MD-390G (S13.20)", gpsURL},
			modelURL{"RT3 (D13.20)", nonGpsURL},
			modelURL{"RT8 (S13.20)", gpsURL},
		}

		title := "Write KD4Z md380tools firmware to radio..."
		upgrade := true
		canceled, model, url := firmwareDialog(title, modelURLs, upgrade)
		if canceled {
			return
		}

		msgs := []string{
			fmt.Sprintf("Downloading KD4Z md380tools %s firmware...\n%s", model, url),
			"Erasing the radio's firmware...",
			fmt.Sprintf("Writing KD4Z md380tools %s firmware to radio...", model),
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

	df, err := dfu.New(func(cur int) error {
		if cur == dfu.MinProgress {
			pd.SetLabelText(msgs[msgIndex])
			msgIndex++
		}
		pd.SetRange(dfu.MinProgress, dfu.MaxProgress)
		pd.SetValue(cur)
		if pd.WasCanceled() {
			return errors.New("cancelled")
		}
		return nil

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

	file, err := os.Open(filename)
	if err != nil {
		logFatalf("writeFirmware: %s", err.Error())
	}

	defer file.Close()

	err = df.WriteFirmware(file)
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

func userdbDialog(title string, labelText string) (canceled, download bool) {
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

func firmwareDialog(title string, modelURLs []modelURL, upgrade bool) (canceled bool, model, url string) {

	models := make([]string, len(modelURLs))
	for i, modelURL := range modelURLs {
		models[i] = modelURL.model
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

	for _, modelURL := range modelURLs {
		if modelURL.model == model {
			url = modelURL.url
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

// Minimum and maximum progress values
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
