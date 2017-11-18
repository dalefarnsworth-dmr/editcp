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
	edt.recordWindow(codeplug.RtGeneralSettings_md380, gsRecord)
}

func gsRecord(edt *editor, recordBox *ui.HBox) {
	r := currentRecord(recordBox.Window())

	mainBox := recordBox.AddVbox()
	row := mainBox.AddHbox()

	column := row.AddVbox()
	groupBox := column.AddGroupbox("Save")
	form := groupBox.AddForm()
	form.AddFieldRows(r,
		codeplug.FtGsSavePreamble,
		codeplug.FtGsSaveModeReceive,
	)

	groupBox = column.AddGroupbox("Alert Tone")
	form = groupBox.AddForm()
	form.AddFieldRows(r,
		codeplug.FtGsDisableAllTones,
		codeplug.FtGsChFreeIndicationTone,
		codeplug.FtGsTalkPermitTone,
		codeplug.FtGsCallAlertToneDuration,
	)

	groupBox = column.AddGroupbox("Scan")
	form = groupBox.AddForm()
	form.AddFieldRows(r,
		codeplug.FtGsScanDigitalHangTime,
		codeplug.FtGsScanAnalogHangTime,
	)

	groupBox = column.AddGroupbox("Lone Worker")
	form = groupBox.AddForm()
	form.AddFieldRows(r,
		codeplug.FtGsLoneWorkerResponseTime,
		codeplug.FtGsLoneWorkerReminderTime,
	)

	groupBox = column.AddGroupbox("Power On Password")
	form = groupBox.AddForm()
	form.AddFieldRows(r,
		codeplug.FtGsPwAndLockEnable,
	)
	form.AddFieldRows(r,
		codeplug.FtGsPowerOnPassword,
	)

	column = row.AddVbox()
	form = column.AddForm()
	form.AddFieldRows(r,
		codeplug.FtGsRadioName,
		codeplug.FtGsRadioID,
		codeplug.FtGsMonitorType,
		codeplug.FtGsVoxSensitivity,
		codeplug.FtGsTxPreambleDuration,
		codeplug.FtGsRxLowBatteryInterval,
		codeplug.FtGsBacklightColor,
		codeplug.FtGsFreqChannelMode,
		codeplug.FtGsModeSelect,
		codeplug.FtGsLockUnlock,
		codeplug.FtGsPcProgPassword,
		codeplug.FtGsRadioProgPassword,
		codeplug.FtGsBacklightTime,
		codeplug.FtGsSetKeypadLockTime,
		codeplug.FtGsDisableAllLeds,
	)

	groupBox = column.AddGroupbox("Talkaround")
	form = groupBox.AddForm()
	form.AddFieldRows(r,
		codeplug.FtGsGroupCallHangTime,
		codeplug.FtGsPrivateCallHangTime,
	)

	groupBox = column.AddGroupbox("Intro Screen")
	form = groupBox.AddForm()
	form.AddFieldRows(r,
		codeplug.FtGsIntroScreen,
		codeplug.FtGsIntroScreenLine1,
		codeplug.FtGsIntroScreenLine2,
	)

	mainBox.AddFiller()
}
