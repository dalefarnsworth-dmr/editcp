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
	"github.com/dalefarnsworth/codeplug/codeplug"
	"github.com/dalefarnsworth/codeplug/ui"
)

func scanLists(edt *editor) {
	writable := true
	edt.recordWindow(codeplug.RtScanLists_md380, writable, slRecord)
}

func slRecord(edt *editor, recordBox *ui.HBox) {
	column := recordBox.AddVbox()
	addFieldMembers(column, &settings.sortAvailableChannels,
		codeplug.FtSlName, codeplug.FtSlChannel_md380, "Channels")

	row := column.AddHbox()

	r := currentRecord(recordBox.Window())

	column = row.AddVbox()
	form := column.AddForm()
	form.AddFieldRows(r,
		codeplug.FtSlPriorityChannel1_md380,
		codeplug.FtSlPriorityChannel2_md380,
		codeplug.FtSlTxDesignatedChannel_md380)

	column = row.AddVbox()
	form = column.AddForm()
	form.AddFieldRows(r,
		codeplug.FtSlSignallingHoldTime,
		codeplug.FtSlPrioritySampleTime)
}
