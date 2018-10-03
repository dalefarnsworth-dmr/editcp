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

func gpsSystems(edt *editor) {
	writable := false
	edt.newRecordWindow(codeplug.RtGPSSystems, writable, gpRecord)
}

func gpRecord(edt *editor, recordBox *ui.HBox) {
	r := currentRecord(recordBox.Window())

	column := recordBox.AddVbox()
	form := column.AddForm()
	form.AddFieldTypeRows(r,
		codeplug.FtGpGPSRevertChannel,
		codeplug.FtGpGPSDefaultReportInterval,
		codeplug.FtGpDestinationID)

	recordBox.AddFiller()
}
