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

	/* Not ready yet
	text := "Select the default model when the codeplug model is unknown.\n"
	column.AddLabel(text)

	cp := edt.codeplug

	models, variantsMap, _ := cp.ModelsVariantsFiles()
	var variants []string

	row = column.AddHbox()
	column.AddSpace(3)
	row2 := column.AddHbox()
	options := append([]string{"Ask"}, models...)
	for i, model := range options {
		b := row.AddRadioButton(model)
		if i == 0 || model == settings.model {
			row2.SetEnabled(i != 0)
			b.SetChecked(true)
		}
		b.ConnectClicked(func(bo bool) {
			text := b.Text()
			row2.SetEnabled(text != "Ask")
			variants = variantsMap[text]
			addVariants(row2, variants)
			settings.model = b.Text()
		})
	}

	addVariants(row2, variants)
	*/

	w.ConnectClose(func() bool {
		saveSettings()
		w.DeleteLater()
		return true
	})

	w.Show()
}

func addVariants(row *ui.HBox, variants []string) {
	row.Clear()
	options := append([]string{"Ask"}, variants...)
	for i, variant := range options {
		b := row.AddRadioButton(variant)
		if i == 0 || variant == settings.variant {
			b.SetChecked(true)
		}
		b.ConnectClicked(func(bo bool) {
			settings.variant = b.Text()
		})
	}
}
