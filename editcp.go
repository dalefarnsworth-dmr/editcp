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
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dalefarnsworth/codeplug/codeplug"
	"github.com/dalefarnsworth/codeplug/ui"
	"github.com/therecipe/qt/core"
)

const autosaveSuffix = ".autosave"
const maxRecentFiles = 10

var editorOpened = false
var editors []*editor

type editorSettings struct {
	sortAvailableChannels  bool
	sortAvailableChannelsB bool
	sortAvailableContacts  bool
	codeplugDirectory      string
	autosaveInterval       int
	recentFiles            []string
	model                  string
	freqRange              string
	displayGPS             bool
}

var appSettings *ui.AppSettings
var settings editorSettings

type editor struct {
	app           *ui.App
	codeplug      *codeplug.Codeplug
	mainWindow    *ui.MainWindow
	undoAction    *ui.Action
	redoAction    *ui.Action
	undoButton    *ui.Button
	redoButton    *ui.Button
	prefWindow    *ui.Window
	autosaveTimer *core.QTimer
	codeplugHash  [sha256.Size]byte
	codeplugCount int
}

func checkAutosave(filename string) {
	asFilename := filename + autosaveSuffix
	asInfo, err := os.Stat(asFilename)
	if err != nil {
		return
	}

	fileInfo, err := os.Stat(filename)
	if err != nil {
		return
	}

	if fileInfo.ModTime().After(asInfo.ModTime()) {
		os.Remove(asFilename)
		return
	}

	backupFilename := filename + ".backup"
	title := "Autosave file found"
	msg := "An automatically saved backup for %s exists. " +
		"Would you like to restore the file from its backup? "
	msg = fmt.Sprintf(msg, filename)
	switch ui.YesNoPopup(title, msg) {
	case ui.PopupYes:
		os.Rename(filename, backupFilename)
		os.Rename(asFilename, filename)
		msg := fmt.Sprintf("%s has been saved as %s",
			filename, backupFilename)
		ui.InfoPopup("Backup created", msg)
	default:
		break
	}
}

func (edt *editor) revertFile() error {
	var err error

	cp := edt.codeplug
	if cp.Changed() {
		title := fmt.Sprintf("Revert %s", cp.Filename())
		msg := cp.Filename() + " has been modified.\n"
		msg += "Are you sure you want to discard the changes?"
		switch ui.YesNoPopup(title, msg) {
		case ui.PopupYes:
			err := edt.codeplug.Revert()
			if err != nil {
				ui.ErrorPopup("Revert Failed", err.Error())
			}
			edt.updateMenuBar()
			ui.ResetWindows(cp)

		default:
			break
		}
	}

	return err
}

func (edt *editor) save() string {
	cp := edt.codeplug
	if cp.Filename() == "." || cp.FileType() != codeplug.FileTypeRdt {
		return edt.saveAs("")
	}
	return edt.saveAs(edt.codeplug.Filename())
}

func baseFilename(filename string) string {
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	if ext != "" {
		base = strings.TrimSuffix(base, ext)
	}

	return base
}

func (edt *editor) saveAs(filename string) string {
	cp := edt.codeplug
	if filename == "" {
		dir := settings.codeplugDirectory
		base := baseFilename(edt.codeplug.Filename())
		ext := cp.Ext()
		dir = filepath.Join(dir, base+"."+ext)
		filename = ui.SaveFilename("Save codeplug file", dir, ext)
		if filename == "" {
			return ""
		}
		settings.codeplugDirectory = filepath.Dir(filename)
		saveSettings()
	}

	valid := cp.Valid()
	edt.updateMenuBar()
	if !valid {
		fmtStr := `
%d records with invalid field values were found in the codeplug.

Click on Cancel and then select "Menu->Edit->Show Invalid Fields" to view them.

Or, click on Ignore to continue saving the file.`
		msg := fmt.Sprintf(fmtStr, len(cp.Warnings()))
		title := fmt.Sprintf("%s: save warning", filename)
		rv := ui.WarningPopup(title, msg)
		if rv != ui.PopupIgnore {
			return ""
		}
	}

	err := cp.SaveAs(filename)
	if err != nil {
		title := fmt.Sprintf("%s: save failed", filename)
		ui.ErrorPopup(title, err.Error())
		return ""
	}

	edt.updateFilename()

	autosaveFilename := cp.Filename() + autosaveSuffix
	os.Remove(autosaveFilename)
	return filename
}

func (edt *editor) setDisplayGPS(displayGPS bool) {
	edt.updateButtons()
	w := edt.mainWindow.RecordWindows()[codeplug.RtChannels_md380]
	if w != nil {
		w.RecordFunc()()
	}
}

func (edt *editor) setAutosaveInterval(seconds int) {
	if seconds == 0 {
		edt.autosaveTimer.Stop()
		return
	}
	if edt.autosaveTimer == nil {
		edt.autosaveTimer = core.NewQTimer(nil)
		edt.autosaveTimer.ConnectTimeout(func() {
			edt.autosave()
		})
	}
	edt.autosaveTimer.Start(seconds * 60 * 1000)
}

func (edt *editor) autosave() {
	cp := edt.codeplug
	if cp == nil {
		return
	}

	filename := cp.Filename() + autosaveSuffix

	hash := cp.CurrentHash()
	if hash == edt.codeplugHash {
		return
	}
	edt.codeplugHash = hash

	err := cp.SaveToFile(filename)
	if err != nil {
		os.Remove(filename)
	}
}

func main() {
	//defer func() {
	//	if r := recover(); r != nil {
	//		logFatal(r)
	//	}
	//}()

	app, err := ui.NewApp()
	if err != nil {
		logPrint(err.Error())
		return
	}
	app.SetOrganizationName("codeplug")
	app.SetApplicationName("Codeplug Editor")
	appSettings = app.NewSettings()
	loadSettings()

	filenames := os.Args[1:]
	if len(filenames) == 0 {
		filenames = []string{""}
	}

	for _, filename := range filenames {
		newEditor(app, codeplug.FileTypeNone, filename)
	}

	if len(editors) == 0 {
		return
	}

	app.Exec()

	saveSettings()
}

func (edt *editor) titleSuffix() string {
	suffix := ""
	if edt.codeplugCount > 1 {
		suffix = fmt.Sprintf(" #%d", edt.codeplugCount)
	}
	return suffix
}

func deleteEditor(i int) {
	copy(editors[i:], editors[i+1:])
	editors[len(editors)-1] = nil
	editors = editors[:len(editors)-1]
}

func deleteString(strs *[]string, i int) {
	copy((*strs)[i:], (*strs)[i+1:])
	(*strs)[len(*strs)-1] = ""
	*strs = (*strs)[:len(*strs)-1]
}

func (edt *editor) openCodeplug(fType codeplug.FileType, filename string) {
	if fType == codeplug.FileTypeNone {
		if absPath, err := filepath.Abs(filename); err == nil {
			filename = absPath
		}
		fileInfo, err := os.Stat(filename)
		if err != nil {
			ui.ErrorPopup(filename, err.Error())
			removeRecentFile(filename)
			return
		}

		for _, cp := range codeplug.Codeplugs() {
			xfileInfo, err := os.Stat(cp.Filename())
			if err == nil && os.SameFile(xfileInfo, fileInfo) {
				edt.codeplug = cp
				break
			}
		}
	}

	if edt.codeplug == nil {
		checkAutosave(filename)

		cp, err := codeplug.NewCodeplug(fType, filename)
		if err != nil {
			ui.ErrorPopup("Codeplug Error", err.Error())
			return
		}

		model, freqRange := modelFrequencyRange(cp)

		if model == "" || freqRange == "" {
			return
		}

		err = cp.Load(model, freqRange)
		if err != nil {
			ui.ErrorPopup("Codeplug Load Error", err.Error())
			return
		}
		if !cp.Valid() {
			fmtStr := `
%d records with invalid field values were found in the codeplug.

Select "Menu->Edit->Show Invalid Fields" to view them.`
			msg := fmt.Sprintf(fmtStr, len(cp.Warnings()))
			ui.InfoPopup("codeplug warning", msg)
		}
		edt.updateMenuBar()

		edt.codeplug = cp
		edt.codeplugHash = edt.codeplug.CurrentHash()
		loadSettings()
		edt.setAutosaveInterval(settings.autosaveInterval)
	}

	if fType == codeplug.FileTypeNone {
		addRecentFile(filename)
	}

	highCount := 0
	cp := edt.codeplug
	for _, edt := range editors {
		if edt.codeplug == cp && edt.codeplugCount > highCount {
			highCount = edt.codeplugCount
		}
	}
	edt.codeplugCount = highCount + 1
}

func (edt *editor) FreeCodeplug() {
	if edt.codeplug == nil {
		return
	}
	edt.codeplug.Free()
	if len(editors) > 1 {
		edt.mainWindow.Close()
		return
	}
	edt.codeplug = nil
	edt.codeplugCount--
	edt.updateFilename()
	edt.updateMenuBar()
	edt.updateButtons()
}

func modelFreqRanges(cp *codeplug.Codeplug, model string) (rangesA, rangesB []string) {
	_, freqRangesMap := cp.ModelsFrequencyRanges()
	rangeMapA := make(map[string]bool)
	rangeMapB := make(map[string]bool)

	freqRanges := freqRangesMap[model]
	for _, r := range freqRanges {
		ranges := strings.Split(r, "_")
		rangeMapA[ranges[0]+" MHz"] = true
		if len(ranges) > 1 {
			rangeMapB[ranges[1]+" MHz"] = true
		}
	}

	rangesA = make([]string, 0)
	for r := range rangeMapA {
		rangesA = append(rangesA, r)
	}

	rangesB = make([]string, 0)
	for r := range rangeMapB {
		rangesB = append(rangesB, r)
	}

	return rangesA, rangesB
}

func modelFrequencyRange(cp *codeplug.Codeplug) (model string, freqRange string) {
	models, freqRangesMap := cp.ModelsFrequencyRanges()
	if len(models) == 1 {
		model = models[0]
		ranges := freqRangesMap[model]
		if ranges != nil && len(ranges) == 1 {
			return model, ranges[0]

		}
	}

	model = settings.model

	mOpts := append([]string{"<select model>"}, models...)

	var vOptsA []string
	var rangeA string
	var rangeB string
	var rangesA = make([]string, 0)
	var rangesB = make([]string, 0)

	rangesA, rangesB = modelFreqRanges(cp, model)
	settingRanges := strings.Split(settings.freqRange, "_")

	rangeA = settingRanges[0] + " MHz"
	if len(rangesB) == 0 {
		vOptsA = append([]string{"<select frequency range>"}, rangesA...)
	} else {
		vOptsA = append([]string{"<select frequency range A>"}, rangesA...)
	}
	if len(settingRanges) > 1 {
		rangeB = settingRanges[1] + " MHz"
	}
	vOptsB := append([]string{"<select frequency range B>"}, rangesB...)

	dialog := ui.NewDialog("Select codeplug type")

	cancelButton := ui.NewButtonWidget("Cancel", func() {
		dialog.Reject()
	})
	okButton := ui.NewButtonWidget("Ok", func() {
		dialog.Accept()
	})
	opt := vOptsA[0]
	enable := containsString(rangeA, vOptsA[1:])
	if enable {
		opt = rangeA
	}
	if len(rangesB) != 0 {
		enable = enable && containsString(rangeB, vOptsB[1:])
	}
	okButton.SetEnabled(enable)

	vCbA := ui.NewComboboxWidget(opt, vOptsA, func(s string) {
		rangeA = s
		enable := containsString(rangeA, vOptsA[1:])
		rangesA, rangesB = modelFreqRanges(cp, model)
		if len(rangesB) != 0 {
			enable = enable && containsString(rangeB, rangesB)
		}
		okButton.SetEnabled(enable)
	})
	vCbA.SetEnabled(containsString(model, mOpts[1:]))

	opt = vOptsB[0]
	if containsString(rangeB, vOptsB[1:]) {
		opt = rangeB
	}

	vCbB := ui.NewComboboxWidget(opt, vOptsB, func(s string) {
		rangeB = s
		enable := containsString(rangeA, vOptsA[1:])
		rangesA, rangesB = modelFreqRanges(cp, model)
		if len(rangesB) != 0 {
			enable = enable && containsString(rangeB, rangesB)
		}
		okButton.SetEnabled(enable)
	})
	vCbB.SetEnabled(containsString(model, mOpts[1:]))

	if len(models) == 1 {
		mOpts = models
	}

	var form *ui.Form
	var mCb *ui.Widget

	mCb = ui.NewComboboxWidget(model, mOpts, func(s string) {
		model = s

		rangesA, rangesB = modelFreqRanges(cp, model)
		settingRanges := strings.Split(settings.freqRange, "_")
		rangeA = settingRanges[0] + " MHz"
		if len(rangesB) == 0 {
			vOptsA = append([]string{"<select frequency range>"}, rangesA...)
		} else {
			vOptsA = append([]string{"<select frequency range A>"}, rangesA...)
		}
		if len(settingRanges) > 1 {
			rangeB = settingRanges[1] + " MHz"
		}
		vOptsB := append([]string{"<select frequency range B>"}, rangesB...)
		vCbA.SetEnabled(containsString(model, mOpts[1:]))

		opt := vOptsA[0]
		enable := containsString(rangeA, vOptsA[1:])
		if enable {
			opt = rangeA
		}
		if len(rangesB) > 1 {
			enable = enable && containsString(rangeB, vOptsB[1:])
		}
		okButton.SetEnabled(enable)

		ui.UpdateComboboxWidget(vCbA, opt, vOptsA)

		vCbA.SetLabel("")
		if len(rangesB) > 1 {
			vCbA.SetLabel("A")
		}

		opt = vOptsB[0]
		if containsString(rangeB, vOptsB[1:]) {
			opt = rangeB
		}

		if len(rangesB) > 1 {
			vCbB.SetEnabled(containsString(model, mOpts[1:]))
			ui.UpdateComboboxWidget(vCbB, opt, vOptsB)
			vCbB.SetVisible(true)
		} else {
			vCbB.SetVisible(false)
		}
	})

	dialog.AddLabel("Select the codeplug model and frequency range.")
	form = dialog.AddForm()
	form.AddRow("", mCb)
	form.AddRow("", vCbA)
	if len(rangesB) > 1 {
		vCbA.SetLabel("A")
	}
	form.AddRow("B", vCbB)
	vCbB.SetVisible(len(rangesB) > 1)

	row := dialog.AddHbox()
	row.AddWidget(cancelButton)
	row.AddWidget(okButton)

	if !dialog.Exec() {
		return "", ""
	}

	freqRange = rangeA
	if len(rangesB) > 1 {
		freqRange += "_" + rangeB
	}
	freqRange = strings.Replace(freqRange, " MHz", "", -1)

	if containsString(model, models) {
		settings.model = model
		settings.freqRange = freqRange
	}
	saveSettings()

	return model, freqRange
}

func containsString(str string, strs []string) bool {
	found := false
	for _, s := range strs {
		if s == str {
			found = true
		}
	}
	return found
}

func newEditor(app *ui.App, fType codeplug.FileType, filename string) *editor {
	var edt *editor
	for _, ed := range editors {
		if ed.codeplug == nil {
			edt = ed
			break
		}
	}

	if edt == nil {
		edt = new(editor)
		edt.app = app
		editors = append(editors, edt)
	}

	mw := edt.mainWindow
	if mw == nil {
		mw = ui.NewMainWindow()
		edt.mainWindow = mw
	}

	if filename != "" || fType != codeplug.FileTypeNone {
		edt.openCodeplug(fType, filename)
	}

	cp := edt.codeplug
	if cp != nil {
		mw.SetCodeplug(cp)
	}

	edt.updateFilename()

	mw.ConnectClose(func() bool {
		if cp != nil {
			count := 0
			for _, edt := range editors {
				if edt.codeplug == cp {
					count++
				}
			}

			if count == 1 && cp.Changed() {
				title := fmt.Sprintf("Save %s", cp.Filename())
				msg := cp.Filename() + " has been modified.\n"
				msg += "Do you want to save the changes?"
				switch ui.SavePopup(title, msg) {
				case ui.PopupSave:
					if edt.save() == "" {
						return false
					}

				case ui.PopupDiscard:
					break

				case ui.PopupCancel:
					return false
				}
			}
		}

		for i, editor := range editors {
			if editor == edt {
				deleteEditor(i)
				break
			}
		}

		if cp != nil {
			asFilename := cp.Filename() + autosaveSuffix
			os.Remove(asFilename)
		}
		return true
	})

	mw.ConnectChange(func(change *codeplug.Change) {
		updateUndoActions(edt)
	})

	edt.updateMenuBar()

	edt.updateButtons()

	mw.Show()

	if len(editors) > 1 && cp == nil {
		mw.Close()
		return nil
	}

	editorOpened = true
	return edt
}

func (edt *editor) updateMenuBar() {
	cp := edt.codeplug
	mb := edt.mainWindow.MenuBar()
	mb.Clear()
	menu := mb.AddMenu("File")
	menu.AddAction("New...", func() {
		newEditor(edt.app, codeplug.FileTypeNew, "")
	})
	menu.AddAction("Open...", func() {
		dir := settings.codeplugDirectory
		exts := edt.codeplug.AllExts()
		filenames := ui.OpenCPFilenames("Open codeplug file", dir, exts)
		for _, filename := range filenames {
			if filename != "" {
				newEditor(edt.app, codeplug.FileTypeNone, filename)
			}
		}
	})
	recentMenu := menu.AddMenu("Open Recent...")
	recentMenu.ConnectAboutToShow(func() {
		edt.updateRecentMenu(recentMenu)
	})
	recentMenu.SetEnabled(len(settings.recentFiles) != 0)

	menu.AddAction("Revert", func() {
		edt.revertFile()
	}).SetEnabled(cp != nil)

	menu.AddSeparator()

	importMenu := menu.AddMenu("Import...")
	importMenu.AddAction("Import text file...", func() {
		edt.importText()
	})

	importMenu.AddAction("Import Spreadsheet file...", func() {
		edt.importXLSX()
	})

	importMenu.AddAction("Import JSON file...", func() {
		edt.importJSON()
	})

	exportMenu := menu.AddMenu("Export...")
	exportMenu.SetEnabled(cp != nil)

	exportMenu.AddAction("Export to text...", func() {
		edt.exportText()
	})

	exportMenu.AddAction("Export to Spreadsheet...", func() {
		edt.exportXLSX()
	})

	exportMenu.AddAction("Export to JSON...", func() {
		edt.exportJSON()
	})

	menu.AddSeparator()

	menu.AddAction("Save", func() {
		edt.save()
	}).SetEnabled(cp != nil)

	menu.AddAction("Save As...", func() {
		edt.saveAs("")
	}).SetEnabled(cp != nil)

	menu.AddSeparator()

	menu.AddAction("Close", func() {
		edt.mainWindow.Close()
	})

	menu.AddAction("Quit", func() {
		for i := len(editors) - 1; i >= 0; i-- {
			editors[i].mainWindow.Close()
		}
	})

	menu = mb.AddMenu("Edit")
	menu.AddAction("Basic Information", func() {
		basicInformation(edt)
	}).SetEnabled(cp != nil)

	menu.AddAction("General Settings", func() {
		generalSettings(edt)
	}).SetEnabled(cp != nil)

	menu.AddAction("Menu Items", func() {
		menuItems(edt)
	}).SetEnabled(cp != nil)

	menu.AddAction("Privacy Settings", func() {
		privacySettings(edt)
	}).SetEnabled(cp != nil)

	menu.AddAction("Channels", func() {
		channels(edt)
	}).SetEnabled(cp != nil)

	menu.AddAction("Contacts", func() {
		contacts(edt)
	}).SetEnabled(cp != nil)

	menu.AddAction("RX Group Lists", func() {
		groupLists(edt)
	}).SetEnabled(cp != nil)

	menu.AddAction("Scan Lists", func() {
		scanLists(edt)
	}).SetEnabled(cp != nil)

	menu.AddAction("Zones", func() {
		zones(edt)
	}).SetEnabled(cp != nil)

	if cp != nil && cp.HasRecordType(codeplug.RtGPSSystems) {
		menu.AddAction("GPS Systems", func() {
			gpsSystems(edt)
		}).SetEnabled(cp != nil && settings.displayGPS)
	}

	menu.AddSeparator()

	edt.undoAction = menu.AddAction("Undo", func() {
		edt.codeplug.UndoChange()
	})
	edt.undoAction.SetEnabled(false)

	edt.redoAction = menu.AddAction("Redo", func() {
		edt.codeplug.RedoChange()
	})
	edt.redoAction.SetEnabled(false)

	menu.AddSeparator()

	menu.AddAction("Show Invalid Fields", func() {
		checkCodeplug(edt)
	}).SetEnabled(cp != nil && len(cp.Warnings()) != 0)

	menu.AddSeparator()

	menu.AddAction("Preferences...", func() {
		edt.preferences()
	})

	edt.addRadioMenu(menu)

	windowsMenu := mb.AddMenu("Windows")
	windowsMenu.ConnectAboutToShow(func() {
		edt.updateWindowsMenu(windowsMenu)
	})

	menu = mb.AddMenu("Help")
	menu.AddAction("About...", func() {
		about()
	})
	menu.AddAction("Thanks...", func() {
		thanks()
	})
}

func (edt *editor) updateButtons() {
	cp := edt.codeplug

	row := edt.mainWindow.AddHbox()
	row.Clear()
	column := row.AddVbox()

	biButton := column.AddButton("Basic Information")
	biButton.SetEnabled(cp != nil)
	biButton.ConnectClicked(func() { basicInformation(edt) })

	gsButton := column.AddButton("General Settings")
	gsButton.SetEnabled(cp != nil)
	gsButton.ConnectClicked(func() { generalSettings(edt) })

	miButton := column.AddButton("Menu Items")
	miButton.SetEnabled(cp != nil)
	miButton.ConnectClicked(func() { menuItems(edt) })

	psButton := column.AddButton("Privacy Settings")
	psButton.SetEnabled(cp != nil)
	psButton.ConnectClicked(func() { privacySettings(edt) })

	ciButton := column.AddButton("Channels")
	ciButton.SetEnabled(cp != nil)
	ciButton.ConnectClicked(func() { channels(edt) })

	dcButton := column.AddButton("Contacts")
	dcButton.SetEnabled(cp != nil)
	dcButton.ConnectClicked(func() { contacts(edt) })

	glButton := column.AddButton("RX Group Lists")
	glButton.SetEnabled(cp != nil)
	glButton.ConnectClicked(func() { groupLists(edt) })

	slButton := column.AddButton("Scan Lists")
	slButton.SetEnabled(cp != nil)
	slButton.ConnectClicked(func() { scanLists(edt) })

	ziButton := column.AddButton("Zones")
	ziButton.SetEnabled(cp != nil)
	ziButton.ConnectClicked(func() { zones(edt) })

	if cp != nil && cp.HasRecordType(codeplug.RtGPSSystems) {
		gpButton := column.AddButton("GPS Systems")
		gpButton.SetEnabled(cp != nil && settings.displayGPS)
		gpButton.ConnectClicked(func() { gpsSystems(edt) })
	}

	row.AddSeparator()

	column = row.AddVbox()

	edt.undoButton = column.AddButton("Undo")
	edt.undoButton.SetFixedHeight()
	edt.undoButton.SetEnabled(false)
	edt.undoButton.ConnectClicked(func() {
		edt.codeplug.UndoChange()
	})

	edt.redoButton = column.AddButton("Redo")
	edt.redoButton.SetFixedHeight()
	edt.redoButton.SetEnabled(false)
	edt.redoButton.ConnectClicked(func() {
		edt.codeplug.RedoChange()
	})

	column.AddFiller()
	row.AddFiller()
}

func (edt *editor) updateFilename() {
	title := "Codeplug Editor"
	cp := edt.codeplug
	if cp != nil {
		filename := cp.Filename()
		title = filename + edt.titleSuffix()
		if _, err := os.Stat(filename); err == nil {
			settings.codeplugDirectory = filepath.Dir(filename)
			saveSettings()
		}
		addRecentFile(filename)
	}

	edt.mainWindow.SetTitle(title)
}

func (edt *editor) updateWindowsMenu(menu *ui.Menu) {
	menu.Clear()

	mainWindows := ui.MainWindows()
	sort.Slice(mainWindows, func(i, j int) bool {
		return mainWindows[i].Title() < mainWindows[j].Title()
	})
	for i := range mainWindows {
		mw := mainWindows[i]
		menu.AddAction(mw.Title(), func() {
			mw.Show()
		})
		windows := make([]*ui.Window, 0, 16)
		for _, w := range mw.RecordWindows() {
			windows = append(windows, w)
		}
		sort.Slice(windows, func(i, j int) bool {
			return windows[i].Title() < windows[j].Title()
		})
		for i := range windows {
			w := windows[i]
			menu.AddAction(w.Title(), func() {
				w.Show()
			})
		}
	}
}

func (edt *editor) updateRecentMenu(menu *ui.Menu) {
	menu.Clear()

	loadSettings()

	for i := range settings.recentFiles {
		filename := settings.recentFiles[i]
		menu.AddAction(filename, func() {
			newEditor(edt.app, codeplug.FileTypeNone, filename)
		})
	}
	menu.SetEnabled(len(settings.recentFiles) != 0)
}

func addRecentFile(name string) {
	if _, err := os.Stat(name); err != nil {
		return
	}

	if len(settings.recentFiles) > 0 && name == settings.recentFiles[0] {
		return
	}

	removeRecentFile(name)

	settings.recentFiles = append([]string{name}, settings.recentFiles...)

	if len(settings.recentFiles) > maxRecentFiles {
		settings.recentFiles = settings.recentFiles[:maxRecentFiles]
	}

	saveSettings()
}

func removeRecentFile(name string) {
	for i, n := range settings.recentFiles {
		if n == name {
			deleteString(&settings.recentFiles, i)
			break
		}
	}
	saveSettings()
}

func (edt *editor) exportText() {
	dir := settings.codeplugDirectory
	base := baseFilename(edt.codeplug.Filename())
	ext := "txt"
	dir = filepath.Join(dir, base+"."+ext)
	filename := ui.SaveFilename("Export to text file", dir, ext)
	if filename == "" {
		return
	}
	settings.codeplugDirectory = filepath.Dir(filename)
	saveSettings()

	err := edt.codeplug.ExportText(filename)
	if err != nil {
		title := fmt.Sprintf("Export to %s", filename)
		ui.ErrorPopup(title, err.Error())
		return
	}
}

func (edt *editor) importText() {
	dir := settings.codeplugDirectory
	filename := ui.OpenTextFilename("Import text file", dir)
	if filename == "" {
		return
	}
	settings.codeplugDirectory = filepath.Dir(filename)
	saveSettings()

	newEditor(edt.app, codeplug.FileTypeText, filename)
}

func (edt *editor) importXLSX() {
	dir := settings.codeplugDirectory
	filename := ui.OpenXLSXFilename("Import Spreadsheet file", dir)
	if filename == "" {
		return
	}
	settings.codeplugDirectory = filepath.Dir(filename)
	saveSettings()

	newEditor(edt.app, codeplug.FileTypeXLSX, filename)
}

func (edt *editor) exportXLSX() {
	dir := settings.codeplugDirectory
	base := baseFilename(edt.codeplug.Filename())
	ext := "xlsx"
	dir = filepath.Join(dir, base+"."+ext)
	filename := ui.SaveFilename("Export to Spreadsheet file", dir, ext)
	if filename == "" {
		return
	}
	settings.codeplugDirectory = filepath.Dir(filename)
	saveSettings()

	err := edt.codeplug.ExportXLSX(filename)
	if err != nil {
		title := fmt.Sprintf("Export to %s", filename)
		ui.ErrorPopup(title, err.Error())
		return
	}
}

func (edt *editor) importJSON() {
	dir := settings.codeplugDirectory
	filename := ui.OpenJSONFilename("Import JSON file", dir)
	if filename == "" {
		return
	}
	settings.codeplugDirectory = filepath.Dir(filename)
	saveSettings()

	newEditor(edt.app, codeplug.FileTypeJSON, filename)
}

func (edt *editor) exportJSON() {
	dir := settings.codeplugDirectory
	base := baseFilename(edt.codeplug.Filename())
	ext := "json"
	dir = filepath.Join(dir, base+"."+ext)
	filename := ui.SaveFilename("Export to JSON file", dir, ext)
	if filename == "" {
		return
	}
	settings.codeplugDirectory = filepath.Dir(filename)
	saveSettings()

	err := edt.codeplug.ExportJSON(filename)
	if err != nil {
		title := fmt.Sprintf("Export to %s", filename)
		ui.ErrorPopup(title, err.Error())
		return
	}
}

func about() {
	msg := fmt.Sprintf("editcp Version %s\n", version)
	msg += `
editcp is free software licensed
under version 3 of the GPL.

Copyright 2017-2018 Dale Farnsworth

Dale Farnsworth
1007 W Mendoza Ave
Mesa, AZ  85210
USA

dale@farnsworth.org

The source code for editcp may be found at
https://github.com/dalefarnsworth/codeplug
`
	ui.InfoPopup("About editcp", msg)
}

func thanks() {
	msgs := []string{
		"A big thank you to:",
		"  José Melo, CT4TX, for creating the nice logo",
		"  Ron McMurdy, W5QLD, for reporting bugs",
		"  Markus Lenggenhager, HB9BRJ, for reporting bugs",
		"  Roy G. Jackson, KW4G, for reporting bugs",
		"  Kevin Otte, N8VNR, for reporting bugs",
		"  Andreas Krüger, DJ3EI, for reporting bugs",
		"  Martin Jones, KI0KO, for reporting bugs",
		"  Marco Carrara, IW2KWD, for suggesting improvements",
		"  Bob Finch, W9YA, for reporting bugs",
		"",
		"Please let me know of errors or omissions in the above list.",
		"Dale Farnsworth, NO7K, dale@farnsworth.org",
	}

	msg := strings.Join(msgs, "\n")
	ui.InfoPopup("Thanks", msg)
}

func updateUndoActions(edt *editor) {
	text := edt.codeplug.UndoString()
	edt.undoAction.SetText("Undo: " + text)
	edt.undoAction.SetEnabled(text != "")

	edt.undoButton.SetText("Undo: " + edt.codeplug.UndoString())
	edt.undoButton.SetEnabled(text != "")

	text = edt.codeplug.RedoString()
	edt.redoAction.SetText("Redo: " + edt.codeplug.RedoString())
	edt.redoAction.SetEnabled(text != "")

	edt.redoButton.SetText("Redo: " + edt.codeplug.RedoString())
	edt.redoButton.SetEnabled(text != "")
}

type fillRecord func(*editor, *ui.HBox)

func (edt *editor) recordWindow(rType codeplug.RecordType, writable bool, fillRecord fillRecord) {
	windows := edt.mainWindow.RecordWindows()
	w := windows[rType]
	if w != nil {
		w.Show()
		return
	}

	w = edt.mainWindow.NewRecordWindow(rType, writable)
	windows[rType] = w

	w.ConnectClose(func() bool {
		delete(windows, rType)
		return true
	})

	cp := edt.codeplug
	r := cp.Record(rType)
	w.SetTitle(cp.Filename() + edt.titleSuffix() + " " + r.TypeName())

	windowBox := w.AddHbox()
	var rl *ui.RecordList
	var recordFunc func()

	if cp.MaxRecords(rType) == 1 {
		selectorBox := windowBox.AddVbox()
		recordFunc = func() {
			selectorBox.Clear()
			recordBox := selectorBox.AddHbox()
			fillRecord(edt, recordBox)
			w.EnableWidgets()
		}
	} else {
		rl = windowBox.AddRecordList(rType)
		if rl.Current() < 0 {
			rl.SetCurrent(0)
		}
		selectorBox := windowBox.AddVbox()
		recordFunc = func() {
			selectorBox.Clear()
			recordBox := selectorBox.AddHbox()
			fillRecord(edt, recordBox)
			addRecordSelector(selectorBox, writable)
			w.EnableWidgets()
		}
	}

	w.SetRecordFunc(recordFunc)
	recordFunc()

	w.Show()
}

func addRecordSelector(box *ui.VBox, writable bool) {
	w := box.Window()
	cp := w.MainWindow().Codeplug()
	rl := w.RecordList()
	rType := w.RecordType()
	row := box.AddHbox()
	row.SetFixedHeight()

	decrement := row.AddButton("<")
	decrement.ConnectClicked(func() {
		rIndex := rl.Current()
		if rIndex <= 0 {
			return
		}

		rIndex--
		rl.SetCurrent(rIndex)
	})

	rIndex := rl.Current()
	records := cp.Records(rType)
	row.AddButton(fmt.Sprintf("%d of %d", rIndex+1, len(records)))
	increment := row.AddButton(">")
	increment.ConnectClicked(func() {
		rIndex := rl.Current()
		records := cp.Records(rType)
		if rIndex >= len(records)-1 {
			return
		}

		rIndex++
		rl.SetCurrent(rIndex)
	})

	if writable {
		row.AddSpace(3)
		add := row.AddButton("Add")
		add.ConnectClicked(func() {
			err := rl.AddSelected()
			if err != nil {
				ui.ErrorPopup("Add Record", err.Error())
				return
			}
		})

		dup := row.AddButton("Dup")
		dup.ConnectClicked(func() {
			err := rl.DupSelected()
			if err != nil {
				ui.ErrorPopup("Dup Record", err.Error())
				return
			}
		})

		row.AddSpace(3)
		delete := row.AddButton("Delete")
		delete.ConnectClicked(func() {
			err := rl.RemoveSelected()
			if err != nil {
				ui.ErrorPopup("Delete Record", err.Error())
				return
			}
		})
	}

	row.AddFiller()
}

func currentRecord(w *ui.Window) *codeplug.Record {
	rIndex := 0
	rl := w.RecordList()
	if rl != nil {
		rIndex = rl.Current()
	}
	records := w.MainWindow().Codeplug().Records(w.RecordType())

	return records[rIndex]
}

func loadSettings() {
	as := appSettings
	as.Sync()
	settings.sortAvailableChannels = as.Bool("sortAvailableChannels", false)
	settings.sortAvailableChannelsB = as.Bool("sortAvailableChannelsB", false)
	settings.sortAvailableContacts = as.Bool("sortAvailableContacts", false)
	settings.codeplugDirectory = as.String("codeplugDirectory", "")
	settings.autosaveInterval = as.Int("autosaveInterval", 1)
	settings.model = as.String("model", "")
	settings.freqRange = as.String("frequencyRange", "")
	settings.displayGPS = as.Bool("displayGPS", true)

	size := as.BeginReadArray("recentFiles")
	settings.recentFiles = make([]string, size)
	for i := 0; i < size; i++ {
		as.SetArrayIndex(i)
		settings.recentFiles[i] = as.String("filename", "")
	}
	as.EndArray()
}

func saveSettings() {
	as := appSettings
	as.SetBool("sortAvailableChannels", settings.sortAvailableChannels)
	as.SetBool("sortAvailableChannelsB", settings.sortAvailableChannelsB)
	as.SetBool("sortAvailableContacts", settings.sortAvailableContacts)
	as.SetString("codeplugDirectory", settings.codeplugDirectory)
	as.SetInt("autosaveInterval", settings.autosaveInterval)
	as.SetString("model", settings.model)
	as.SetString("frequencyRange", settings.freqRange)
	as.SetBool("displayGPS", settings.displayGPS)

	as.BeginWriteArray("recentFiles", len(settings.recentFiles))
	for i, name := range settings.recentFiles {
		as.SetArrayIndex(i)
		as.SetString("filename", name)
	}
	as.EndArray()

	as.Sync()
}
