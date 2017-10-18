// Copyright 2017 Dale Farnsworth. All rights reserved.

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
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/dalefarnsworth/codeplug/codeplug"
	"github.com/dalefarnsworth/codeplug/ui"
	"github.com/therecipe/qt/core"
)

const autosaveSuffix = ".autosave"
const maxRecentFiles = 10

var editorOpened = false
var editors []*editor

type editorSettings struct {
	sortAvailableChannels bool
	sortAvailableContacts bool
	codeplugDirectory     string
	autosaveInterval      int
	recentFiles           []string
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

	fInfo, err := os.Stat(filename)
	if err != nil {
		return
	}

	if fInfo.ModTime().After(asInfo.ModTime()) {
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
				ui.WarningPopup("Revert Failed", err.Error())
			}
			ui.ResetWindows(cp)

		default:
			break
		}
	}

	return err
}

func (edt *editor) save() {
	edt.saveAs(edt.codeplug.Filename())
}

func (edt *editor) saveAs(filename string) {
	autosaveFilename := edt.codeplug.Filename() + autosaveSuffix

	if filename == "" {
		dir := settings.codeplugDirectory
		filename = ui.SaveFilename("Save codeplug file", dir)
		if filename == "" {
			return
		}
		settings.codeplugDirectory = filepath.Dir(filename)
		saveSettings()
	}
	err := edt.codeplug.SaveAs(filename)
	if err != nil {
		title := fmt.Sprintf("%s: save failed", filename)
		ui.WarningPopup(title, err.Error())
		return
	}

	os.Remove(autosaveFilename)
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
	app, err := ui.NewApp()
	if err != nil {
		log.Print(err.Error())
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
		newEditor(app, filename)
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

func deleteEditor(editors *[]*editor, i int) {
	copy((*editors)[i:], (*editors)[i+1:])
	(*editors)[len(*editors)-1] = nil
	*editors = (*editors)[:len(*editors)-1]
}

func deleteString(strs *[]string, i int) {
	copy((*strs)[i:], (*strs)[i+1:])
	(*strs)[len(*strs)-1] = ""
	*strs = (*strs)[:len(*strs)-1]
}

func (edt *editor) openCodeplugFile(filename string) {
	if absPath, err := filepath.Abs(filename); err == nil {
		filename = absPath
	}

	fInfo, err := os.Stat(filename)
	if err != nil {
		ui.WarningPopup(filename, err.Error())
		removeRecentFile(filename)
		return
	}

	for _, cp := range codeplug.Codeplugs() {
		xfInfo, err := os.Stat(cp.Filename())
		if err == nil && os.SameFile(xfInfo, fInfo) {
			edt.codeplug = cp
			break
		}
	}

	if edt.codeplug == nil {
		checkAutosave(filename)

		cp, err := codeplug.NewCodeplug(filename, codeplug.CtMd380)
		if err != nil {
			ui.WarningPopup("Codeplug Error", err.Error())
			return
		}

		edt.codeplug = cp
		edt.codeplugHash = edt.codeplug.CurrentHash()
		loadSettings()
		edt.setAutosaveInterval(settings.autosaveInterval)
	}

	addRecentFile(filename)

	highCount := 0
	cp := edt.codeplug
	for _, edt := range editors {
		if edt.codeplug == cp && edt.codeplugCount > highCount {
			highCount = edt.codeplugCount
		}
	}
	edt.codeplugCount = highCount + 1
}

func newEditor(app *ui.App, filename string) {
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

	if filename != "" {
		edt.openCodeplugFile(filename)
	}

	cp := edt.codeplug
	if cp != nil {
		mw.SetCodeplug(cp)
	}

	title := "md-380 Codeplug Editor"
	if cp != nil {
		title = filename + edt.titleSuffix()
	}
	mw.SetTitle(title)

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
					edt.save()

				case ui.PopupDiscard:
					break

				case ui.PopupCancel:
					return false
				}
			}
		}

		for i, editor := range editors {
			if editor == edt {
				deleteEditor(&editors, i)
				break
			}
		}

		asFilename := filename + autosaveSuffix
		os.Remove(asFilename)
		return true
	})

	mw.ConnectChange(func(change *codeplug.Change) {
		updateUndoActions(edt)
	})

	mb := mw.MenuBar()
	mb.Clear()
	menu := mb.AddMenu("File")
	menu.AddAction("Open...", func() {
		dir := settings.codeplugDirectory
		filename = ui.OpenFilename("Open codeplug file", dir)
		if filename == "" {
			return
		}
		settings.codeplugDirectory = filepath.Dir(filename)
		saveSettings()

		newEditor(edt.app, filename)
	})
	recentMenu := menu.AddMenu("Open Recent...")
	recentMenu.ConnectAboutToShow(func() {
		edt.updateRecentMenu(recentMenu)
	})

	menu.AddAction("Revert", func() {
		edt.revertFile()
	}).SetDisabled(cp == nil)

	menu.AddAction("Export to text file...", func() {
		edt.exportText()
	}).SetDisabled(cp == nil)

	menu.AddAction("Import from text file...", func() {
		edt.importText()
	}).SetDisabled(cp == nil)

	menu.AddAction("Save", func() {
		edt.save()
	}).SetDisabled(cp == nil)

	menu.AddAction("Save As...", func() {
		edt.saveAs("")
	}).SetDisabled(cp == nil)

	menu.AddAction("Preferences...", func() {
		edt.preferences()
	})

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
	menu.AddAction("General Settings", func() {
		generalSettings(edt)
	}).SetDisabled(cp == nil)

	menu.AddAction("Channel Information", func() {
		channelInformation(edt)
	}).SetDisabled(cp == nil)

	menu.AddAction("Digital Contacts", func() {
		digitalContacts(edt)
	}).SetDisabled(cp == nil)

	menu.AddAction("Digital Rx Group Lists", func() {
		groupLists(edt)
	}).SetDisabled(cp == nil)

	menu.AddAction("Scan Lists", func() {
		scanLists(edt)
	}).SetDisabled(cp == nil)

	menu.AddAction("Zone Information", func() {
		zoneInformation(edt)
	}).SetDisabled(cp == nil)

	edt.undoAction = menu.AddAction("Undo", func() {
		edt.codeplug.UndoChange()
	})
	edt.undoAction.SetDisabled(true)

	edt.redoAction = menu.AddAction("Redo", func() {
		edt.codeplug.RedoChange()
	})
	edt.redoAction.SetDisabled(true)

	windowsMenu := mb.AddMenu("Windows")
	windowsMenu.ConnectAboutToShow(func() {
		edt.updateWindowsMenu(windowsMenu)
	})

	menu = mb.AddMenu("Help")
	menu.AddAction("About...", func() {
		about()
	})

	row := mw.AddHbox()
	column := row.AddVbox()

	gsButton := column.AddButton("GeneralSettings")
	gsButton.SetDisabled(cp == nil)
	gsButton.ConnectClicked(func() { generalSettings(edt) })

	ciButton := column.AddButton("Channel Information")
	ciButton.SetDisabled(cp == nil)
	ciButton.ConnectClicked(func() { channelInformation(edt) })

	dcButton := column.AddButton("Digital Contacts")
	dcButton.SetDisabled(cp == nil)
	dcButton.ConnectClicked(func() { digitalContacts(edt) })

	glButton := column.AddButton("Digital Rx Group Lists")
	glButton.SetDisabled(cp == nil)
	glButton.ConnectClicked(func() { groupLists(edt) })

	slButton := column.AddButton("Scan Lists")
	slButton.SetDisabled(cp == nil)
	slButton.ConnectClicked(func() { scanLists(edt) })

	ziButton := column.AddButton("Zone Information")
	ziButton.SetDisabled(cp == nil)
	ziButton.ConnectClicked(func() { zoneInformation(edt) })

	column.AddFiller()
	row.AddSeparator()

	column = row.AddVbox()

	edt.undoButton = column.AddButton("Undo")
	edt.undoButton.SetDisabled(true)
	edt.undoButton.ConnectClicked(func() {
		edt.codeplug.UndoChange()
	})

	edt.redoButton = column.AddButton("Redo")
	edt.redoButton.SetDisabled(true)
	edt.redoButton.ConnectClicked(func() {
		edt.codeplug.RedoChange()
	})

	column.AddFiller()

	row.AddFiller()

	mw.Show()

	editorOpened = true
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
			newEditor(edt.app, filename)
		})
	}
	menu.SetEnabled(len(settings.recentFiles) != 0)
}

func addRecentFile(name string) {
	if len(settings.recentFiles) > 0 {
		if name == settings.recentFiles[0] {
			return
		}
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
	filename := ui.SaveFilename("Export to text file", dir)
	if filename == "" {
		return
	}
	settings.codeplugDirectory = filepath.Dir(filename)
	saveSettings()

	err := edt.codeplug.ExportTo(filename)
	if err != nil {
		title := fmt.Sprintf("Export to %s", filename)
		ui.WarningPopup(title, err.Error())
		return
	}
}

func (edt *editor) importText() {
	cp := edt.codeplug
	if cp.Changed() {
		title := fmt.Sprintf("Import from text file")
		msg := cp.Filename() + " has been modified.\n"
		msg += "Are you sure you want to discard the changes?"
		switch ui.YesNoPopup(title, msg) {
		case ui.PopupYes:
			break

		default:
			return
		}
	}

	dir := settings.codeplugDirectory
	filename := ui.OpenFilename("Import from text file", dir)
	if filename == "" {
		return
	}
	settings.codeplugDirectory = filepath.Dir(filename)
	saveSettings()

	_, err := os.Stat(filename)
	if err != nil {
		title := fmt.Sprintf("Import from %s", filename)
		ui.WarningPopup(title, err.Error())
		return
	}

	windows := edt.mainWindow.RecordWindows()
	for _, w := range windows {
		w.Close()
	}

	err = edt.codeplug.ImportFrom(filename)
	if err != nil {
		title := fmt.Sprintf("Import from %s failed", filename)
		msg := err.Error()
		ui.WarningPopup(title, msg)
		return
	}

	ui.ResetWindows(edt.codeplug)
}

func about() {
	msg := fmt.Sprintf("editcp Version %s\n", version)
	msg += `
editcp is free software licensed
under version 3 of the GPL.

Copyright 2017 Dale Farnsworth

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

func updateUndoActions(edt *editor) {
	text := edt.codeplug.UndoString()
	edt.undoAction.SetText("Undo: " + text)
	edt.undoAction.SetDisabled(text == "")

	edt.undoButton.SetText("Undo: " + edt.codeplug.UndoString())
	edt.undoButton.SetDisabled(text == "")

	text = edt.codeplug.RedoString()
	edt.redoAction.SetText("Redo: " + edt.codeplug.RedoString())
	edt.redoAction.SetDisabled(text == "")

	edt.redoButton.SetText("Redo: " + edt.codeplug.RedoString())
	edt.redoButton.SetDisabled(text == "")
}

type fillRecord func(*editor, *ui.HBox)

func (edt *editor) recordWindow(rType codeplug.RecordType, fillRecord fillRecord) {
	windows := edt.mainWindow.RecordWindows()
	w := windows[rType]
	if w != nil {
		w.Show()
		return
	}

	w = edt.mainWindow.NewRecordWindow(rType)
	windows[rType] = w

	w.ConnectClose(func() bool {
		delete(windows, rType)
		return true
	})

	cp := edt.codeplug
	r := cp.Record(rType)
	w.SetTitle(cp.Filename() + edt.titleSuffix() + " " + r.TypeName())

	windowBox := w.AddHbox()
	windowBox.SetContentsMargins(0, 0, 0, 0)
	var rl *ui.RecordList
	var recordFunc func()

	if cp.MaxRecords(rType) == 1 {
		selectorBox := windowBox.AddVbox()
		selectorBox.SetContentsMargins(0, 0, 0, 0)
		recordFunc = func() {
			selectorBox.Clear()
			recordBox := selectorBox.AddHbox()
			recordBox.SetContentsMargins(0, 0, 0, 0)
			fillRecord(edt, recordBox)
			w.EnableWidgets()
		}
	} else {
		rl = windowBox.AddRecordList(rType)
		if rl.Current() < 0 {
			rl.SetCurrent(0)
		}
		selectorBox := windowBox.AddVbox()
		selectorBox.SetContentsMargins(0, 0, 0, 0)
		recordFunc = func() {
			selectorBox.Clear()
			recordBox := selectorBox.AddHbox()
			recordBox.SetContentsMargins(0, 0, 0, 0)
			fillRecord(edt, recordBox)
			addRecordSelector(selectorBox)
			w.EnableWidgets()
		}
	}

	w.SetRecordFunc(recordFunc)
	recordFunc()

	w.Show()
}

func addRecordSelector(box *ui.VBox) {
	w := box.Window()
	cp := w.MainWindow().Codeplug()
	rl := w.RecordList()
	rType := w.RecordType()
	row := box.AddHbox()
	row.SetContentsMargins(0, 0, 0, 0)

	decrement := row.AddButton("<")
	rIndex := rl.Current()
	records := cp.Records(rType)
	row.AddButton(fmt.Sprintf("%d of %d", rIndex+1, len(records)))
	increment := row.AddButton(">")
	row.AddSpace(3)
	add := row.AddButton("Add")
	row.AddSpace(3)
	delete := row.AddButton("Delete")
	row.AddFiller()
	box.AddFiller()

	decrement.ConnectClicked(func() {
		rIndex := rl.Current()
		if rIndex <= 0 {
			return
		}

		rIndex--
		rl.SetCurrent(rIndex)
	})

	increment.ConnectClicked(func() {
		rIndex := rl.Current()
		records := cp.Records(rType)
		if rIndex >= len(records)-1 {
			return
		}

		rIndex++
		rl.SetCurrent(rIndex)
	})

	add.ConnectClicked(func() {
		err := rl.AddSelected()
		if err != nil {
			ui.WarningPopup("Add Record", err.Error())
			return
		}
	})

	delete.ConnectClicked(func() {
		err := rl.RemoveSelected()
		if err != nil {
			ui.WarningPopup("Delete Record", err.Error())
			return
		}
	})
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

func addFieldMembers(vBox *ui.VBox, sortAvailable *bool, nameType codeplug.FieldType, memberType codeplug.FieldType, headerName string) *ui.FieldMembers {
	r := currentRecord(vBox.Window())

	return vBox.AddFieldMembers(r, sortAvailable,
		nameType, memberType, headerName)
}

func loadSettings() {
	as := appSettings
	as.Sync()
	settings.sortAvailableChannels = as.Bool("sortAvailableChannels", false)
	settings.sortAvailableContacts = as.Bool("sortAvailableContacts", false)
	settings.codeplugDirectory = as.String("codeplugDirectory", "")
	settings.autosaveInterval = as.Int("autosaveInterval", 1)
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
	as.SetBool("sortAvailableContacts", settings.sortAvailableContacts)
	as.SetString("codeplugDirectory", settings.codeplugDirectory)
	as.SetInt("autosaveInterval", settings.autosaveInterval)
	as.BeginWriteArray("recentFiles", len(settings.recentFiles))
	for i, name := range settings.recentFiles {
		as.SetArrayIndex(i)
		as.SetString("filename", name)
	}
	as.EndArray()
	as.Sync()
}
