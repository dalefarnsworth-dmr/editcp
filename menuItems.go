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

func menuItems(edt *editor) {
	writable := false
	edt.recordWindow(codeplug.RtMenuItems, writable, miRecord)
}

func miRecord(edt *editor, recordBox *ui.HBox) {
	r := currentRecord(recordBox.Window())

	mainBox := recordBox.AddVbox()
	row := mainBox.AddHbox()

	column := row.AddVbox()
	row2 := column.AddHbox()
	column2 := row2.AddVbox()
	form := column2.AddForm()
	form.AddFieldTypeRows(r,
		codeplug.FtMiHangTime,
		codeplug.FtMiTextMessage,
	)

	column2 = row2.AddVbox()

	groupBox := column.AddGroupbox("Contacts")
	column2 = groupBox.AddVbox()
	form = column2.AddForm()
	form.AddFieldTypeRows(r,
		codeplug.FtMiCallAlert,
		codeplug.FtMiManualDial,
		codeplug.FtMiRemoteMonitor,
		codeplug.FtMiRadioEnable,
	)

	column2 = groupBox.AddVbox()
	form = column2.AddForm()
	form.AddFieldTypeRows(r,
		codeplug.FtMiEdit,
		codeplug.FtMiRadioCheck,
		codeplug.FtMiProgramKey,
		codeplug.FtMiRadioDisable,
	)

	groupBox = column.AddGroupbox("Call Log")
	column2 = groupBox.AddVbox()
	form = column2.AddForm()
	form.AddFieldTypeRows(r,
		codeplug.FtMiMissed,
		codeplug.FtMiOutgoingRadio,
	)

	column2 = groupBox.AddVbox()
	form = column2.AddForm()
	form.AddFieldTypeRows(r,
		codeplug.FtMiAnswered,
	)

	column = row.AddVbox()
	groupBox = column.AddGroupbox("Utilities")
	column2 = groupBox.AddVbox()
	form = column2.AddForm()
	form.AddFieldTypeRows(r,
		codeplug.FtMiTalkaround,
		codeplug.FtMiPower,
		codeplug.FtMiIntroScreen,
		codeplug.FtMiLedIndicator,
		codeplug.FtMiPasswordAndLock,
		codeplug.FtMiDisplayMode,
		codeplug.FtMiGps,
	)

	column2 = groupBox.AddVbox()
	form = column2.AddForm()
	form.AddFieldTypeRows(r,
		codeplug.FtMiToneOrAlert,
		codeplug.FtMiBacklight,
		codeplug.FtMiKeyboardLock,
		codeplug.FtMiSquelch,
		codeplug.FtMiVox,
		codeplug.FtMiProgramRadio,
	)

	groupBox = column.AddGroupbox("Scan")
	column2 = groupBox.AddVbox()
	form = column2.AddForm()
	form.AddFieldTypeRows(r,
		codeplug.FtMiScan,
	)

	column2 = groupBox.AddVbox()
	form = column2.AddForm()
	form.AddFieldTypeRows(r,
		codeplug.FtMiEditList,
	)
}
