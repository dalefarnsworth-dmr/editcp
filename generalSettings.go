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

func generalSettings(edt *editor) {
	edt.recordWindow(codeplug.RtGeneralSettings, gsRecord)
}

func gsRecord(edt *editor, recordBox *ui.HBox) {
	r := currentRecord(recordBox.Window())

	mainBox := recordBox.AddVbox()
	row := mainBox.AddHbox()

	column := row.AddVbox()
	groupBox := column.AddGroupbox("Save")
	form := groupBox.AddForm()
	form.AddFieldRows(r,
		codeplug.FtSavePreamble,
		codeplug.FtSaveModeReceive,
	)

	groupBox = column.AddGroupbox("Alert Tone")
	form = groupBox.AddForm()
	form.AddFieldRows(r,
		codeplug.FtDisableAllTones,
		codeplug.FtChFreeIndicationTone,
		codeplug.FtTalkPermitTone,
		codeplug.FtCallAlertToneDuration,
	)

	groupBox = column.AddGroupbox("Scan")
	form = groupBox.AddForm()
	form.AddFieldRows(r,
		codeplug.FtScanDigitalHangTime,
		codeplug.FtScanAnalogHangTime,
	)

	groupBox = column.AddGroupbox("Lone Worker")
	form = groupBox.AddForm()
	form.AddFieldRows(r,
		codeplug.FtLoneWorkerResponseTime,
		codeplug.FtLoneWorkerReminderTime,
	)

	groupBox = column.AddGroupbox("Power On Password")
	form = groupBox.AddForm()
	form.AddFieldRows(r,
		codeplug.FtPwAndLockEnable,
	)
	form.AddEnabledFieldRows(r, codeplug.FtPwAndLockEnable, "On",
		codeplug.FtPowerOnPassword,
	)

	column = row.AddVbox()
	form = column.AddForm()
	form.AddFieldRows(r,
		codeplug.FtRadioName,
		codeplug.FtRadioID,
		codeplug.FtMonitorType,
		codeplug.FtVoxSensitivity,
		codeplug.FtTxPreambleDuration,
		codeplug.FtRxLowBatteryInterval,
		codeplug.FtPcProgPw,
		codeplug.FtRadioProgPw,
		codeplug.FtSetKeypadLockTime,
		codeplug.FtDisableAllLeds,
	)

	groupBox = column.AddGroupbox("Talkaround")
	form = groupBox.AddForm()
	form.AddFieldRows(r,
		codeplug.FtGroupCallHangTime,
		codeplug.FtPrivateCallHangTime,
	)

	groupBox = column.AddGroupbox("Intro Screen")
	form = groupBox.AddForm()
	form.AddFieldRows(r,
		codeplug.FtIntroScreen,
		codeplug.FtIntroScreenLine1,
		codeplug.FtIntroScreenLine2,
	)

	mainBox.AddFiller()
}
