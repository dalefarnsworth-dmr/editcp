// Copyright 2017-2019 Dale Farnsworth. All rights reserved.

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

	"github.com/dalefarnsworth-dmr/codeplug"
	"github.com/dalefarnsworth-dmr/ui"
)

func privacySettings(edt *editor) {
	writable := false
	edt.newRecordWindow(codeplug.RtPrivacySettings, writable, psRecord)
}

func psRecord(edt *editor, recordBox *ui.HBox) {
	labelFunc := func(f *codeplug.Field) string {
		return fmt.Sprintf("Key %d", f.Index()+1)
	}
	r := currentRecord(recordBox.Window())

	mainBox := recordBox.AddVbox()
	row := mainBox.AddHbox()

	column := row.AddVbox()
	groupBox := column.AddGroupbox("Key Value (Basic)")
	form := groupBox.AddForm()
	form.AddFieldRows(labelFunc, r.Fields(codeplug.FtPsBasicKey)...)

	row.AddSpace(3)

	column = row.AddVbox()
	groupBox = column.AddGroupbox("Key Value (Enhanced)")
	form = groupBox.AddForm()
	form.AddFieldRows(labelFunc, r.Fields(codeplug.FtPsEnhancedKey)...)
	column.AddFiller()
}
