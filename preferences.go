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
	"github.com/dalefarnsworth/codeplug/codeplug"
	"github.com/dalefarnsworth/codeplug/ui"
)

func (edt *editor) preferences() {
	w := edt.mainWindow.NewWindow()

	column := w.AddVbox()
	row := column.AddHbox()
	groupBox := row.AddGroupbox("AutoSave")
	form := groupBox.AddForm()

	loadSettings()

	spinbox := ui.NewSpinbox(settings.autosaveInterval, 0, 60, func(i int) {
		edt.setAutosaveInterval(i)
		settings.autosaveInterval = i
	})
	form.AddRow("Auto Save interval (minutes):", spinbox)
	row.AddFiller()

	text := "Select the action to be taken when the codeplug type " +
		"is ambiguous\n" +
		"  Ask: ask each time a codplug file is opened.\n" +
		"  Type: the file is opened as that type of codeplug."
	column.AddLabel(text)

	nameMap := make(map[string]bool)
	for _, name := range settings.ambiguousNames {
		nameMap[name] = true
	}

	ambigs := codeplug.AmbiguousCodeplugNames()
	buttons := make([]*ui.RadioButton, 0)
	for _, names := range ambigs {
		row := column.AddHbox()
		b := row.AddRadioButton("Ask")
		b.SetChecked(true)
		for _, name := range names {
			b := row.AddRadioButton(name)
			if nameMap[name] {
				b.SetChecked(true)
			}
			buttons = append(buttons, b)
		}
		column.AddSpace(3)
	}

	w.ConnectClose(func() bool {
		names := make([]string, 0)
		for _, b := range buttons {
			if b.IsChecked() {
				names = append(names, b.Text())
			}
		}
		settings.ambiguousNames = names
		saveSettings()

		w.DeleteLater()
		return true
	})

	w.Show()
}
