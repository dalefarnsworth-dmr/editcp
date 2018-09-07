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
	"github.com/dalefarnsworth/codeplug/ui"
)

func (edt *editor) preferences() {
	dialog := ui.NewDialog("Preferences")

	loadSettings()

	row := dialog.AddHbox()
	groupBox := row.AddGroupbox("Display Options")
	form := groupBox.AddForm()

	displayGPS := settings.displayGPS

	checked := displayGPS
	checkbox := ui.NewCheckboxWidget(checked, func(checked bool) {
		displayGPS = checked
	})
	form.AddRow("Display GPS fields:", checkbox)
	dialog.AddSpace(2)

	row = dialog.AddHbox()
	groupBox = row.AddGroupbox("AutoSave")
	form = groupBox.AddForm()

	autosaveInterval := settings.autosaveInterval

	spinbox := ui.NewSpinboxWidget(autosaveInterval, 0, 60, func(i int) {
		autosaveInterval = i
	})
	form.AddRow("Auto Save interval (minutes):", spinbox)
	dialog.AddSpace(2)

	row = dialog.AddHbox()

	cancelButton := ui.NewButtonWidget("Cancel", func() {
		dialog.Reject()
	})
	row.AddWidget(cancelButton)

	okButton := ui.NewButtonWidget("Save", func() {
		dialog.Accept()
	})
	row.AddWidget(okButton)

	if !dialog.Exec() {
		return
	}

	settings.displayGPS = displayGPS
	edt.setDisplayGPS(displayGPS)

	settings.autosaveInterval = autosaveInterval
	edt.setAutosaveInterval(autosaveInterval)
	saveSettings()
}
