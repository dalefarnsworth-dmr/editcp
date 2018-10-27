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
	"strings"

	"github.com/dalefarnsworth/codeplug/codeplug"
	"github.com/dalefarnsworth/codeplug/ui"
)

func basicInformation(edt *editor) {
	writable := false
	edt.newRecordWindow(codeplug.RtBasicInformation_md380, writable, biRecord)
}

func biRecord(edt *editor, recordBox *ui.HBox) {
	w := recordBox.Window()
	r := currentRecord(w)
	cp := edt.codeplug

	column := recordBox.AddVbox()
	form := column.AddForm()

	model, types := cp.ModelTypes()
	if len(types) > 0 && (len(types) != 1 || types[0] != model) {
		model += " (" + strings.Join(types, ", ") + ")"
	}
	modelWidget := ui.NewLineEditWidget(model)
	modelWidget.SetLabel("Model Name")
	form.AddWidget(modelWidget)

	form.AddReadOnlyFieldTypeRows(r,
		codeplug.FtBiFrequencyRange_md380,
		codeplug.FtBiFrequencyRangeA,
		codeplug.FtBiFrequencyRangeB,
		codeplug.FtBiLastProgrammedTime,
	)
}
