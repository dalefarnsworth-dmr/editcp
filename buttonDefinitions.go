// Copyright 2018 Dale Farnsworth. All rights reserved.

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
	"fmt"

	"github.com/dalefarnsworth/codeplug/codeplug"
	"github.com/dalefarnsworth/codeplug/ui"
)

func buttonDefinitions(edt *editor) {
	writable := true
	edt.newRecordWindow(codeplug.RtButtonDefinitions, writable, bdRecord)
}

func fieldTypeWidget(w *ui.Window, r *codeplug.Record, fType codeplug.FieldType) *ui.FieldWidget {
	return w.NewFieldWidget("", r.Field(fType))
}

func stack(w *ui.Window, r *codeplug.Record, fTypes ...codeplug.FieldType) *ui.StackedWidget {
	sw := ui.NewStackedWidget(w)
	for _, fType := range fTypes {
		field := r.Field(fType)
		widget := w.NewFieldWidget("", field)
		sw.AddWidget(widget)
	}

	return sw
}

func bdRecord(edt *editor, recordBox *ui.HBox) {
	cp := edt.codeplug
	r := currentRecord(recordBox.Window())

	mainColumn := recordBox.AddVbox()
	row := mainColumn.AddHbox()
	row.SetFixedHeight()
	column := row.AddVbox()
	column.SetFixedWidth()
	form := column.AddForm()
	form.AddFieldTypeRows(r, codeplug.FtBdLongPressDuration)

	row.AddFiller()

	mainColumn.AddSpace(1)

	mainRow := mainColumn.AddHbox()
	column1 := mainRow.AddVbox()
	column1.SetFixedWidth()

	groupBox := column1.AddGroupbox("Radio Buttons")
	form = groupBox.AddForm()

	column1.AddFiller()

	records := cp.Records(codeplug.RtRadioButtons)
	for _, r := range records {
		form.AddFieldTypeRows(r, codeplug.FtRbButton)
	}

	if cp.HasRecordType(codeplug.RtRadioButtons2) {
		records = cp.Records(codeplug.RtRadioButtons2)
		for _, r := range records {
			form.AddFieldTypeRows(r, codeplug.FtRbButton)
		}
	}

	column2 := mainRow.AddVbox()

	row = column2.AddHbox()
	column = row.AddVbox()
	row.AddFiller()
	column.SetFixedWidth()
	groupBox = column.AddGroupbox("One Touch Access")
	table := groupBox.AddTable()
	w := edt.recordWindow(r.Type())
	records = w.Records(codeplug.RtOneTouch)
	labels := make([]string, 0)

	for _, r := range records {
		table.AddRow(
			fieldTypeWidget(w, r, codeplug.FtOtMode),
			fieldTypeWidget(w, r, codeplug.FtOtCall),

			stack(w, r, codeplug.FtOtCallType, codeplug.FtOtDTMF),
			stack(w, r, codeplug.FtOtTextMessage, codeplug.FtOtEncode),
		)
		labels = append(labels, fmt.Sprintf("%d", r.Index()+1))
	}

	table.AddLeftLabels(labels)

	r = records[0]
	labels = []string{
		r.Field(codeplug.FtOtMode).TypeName(),
		r.Field(codeplug.FtOtCall).TypeName(),
		"Call Type",
		"Message/Encode",
	}
	table.AddTopLabels(labels)

	table.ResizeToContents()

	column2.AddSpace(1)

	row = column2.AddHbox()
	column = row.AddVbox()
	column.SetFixedWidth()
	row.AddFiller()

	groupBox = column.AddGroupbox("Number Key Quick Contact Access")
	table = groupBox.AddTable()
	records = w.Records(codeplug.RtNumberKey)
	labels = make([]string, 0)

	for _, r := range records {
		table.AddRow(
			fieldTypeWidget(w, r, codeplug.FtNkContact),
		)
		labels = append(labels, fmt.Sprintf(" Number key %d", r.Index()))
	}

	table.AddLeftLabels(labels)

	r = records[0]
	labels = []string{r.Field(codeplug.FtNkContact).TypeName()}
	table.AddTopLabels(labels)
	table.ResizeToContents()
	column.SetFixedWidth()
	column2.AddFiller()
}
