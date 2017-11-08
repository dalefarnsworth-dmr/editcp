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
	"github.com/dalefarnsworth/codeplug/ui"
)

func (edt *editor) preferences() {
	dialog := ui.NewDialog("Preferences")

	row := dialog.AddHbox()
	groupBox := row.AddGroupbox("AutoSave")
	form := groupBox.AddForm()

	loadSettings()

	var autoSaveInterval int

	spinbox := ui.NewSpinboxWidget(settings.autosaveInterval, 0, 60, func(i int) {
		autoSaveInterval = i
	})
	form.AddRow("Auto Save interval (minutes):", spinbox)
	row.AddFiller()

	dialog.AddSpace(3)
	row = dialog.AddHbox()

	cancelButton := ui.NewButtonWidget("Cancel", func() {
		dialog.Reject()
	})
	row.AddWidget(cancelButton)

	okButton := ui.NewButtonWidget("Ok", func() {
		dialog.Accept()
	})
	row.AddWidget(okButton)

	dialog.AddFiller()

	if !dialog.Exec() {
		return
	}

	edt.setAutosaveInterval(autoSaveInterval)
	settings.autosaveInterval = autoSaveInterval
	saveSettings()
}
