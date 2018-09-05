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

func channels(edt *editor) {
	writable := true
	edt.recordWindow(codeplug.RtChannels_md380, writable, ciRecord)
}

func ciRecord(edt *editor, recordBox *ui.HBox) {
	mainBox := recordBox.AddVbox()
	row := mainBox.AddHbox()
	column := row.AddVbox()

	groupBox := column.AddGroupbox("Digital/Analog Data")
	column = groupBox.AddVbox()
	form := column.AddForm()

	r := currentRecord(recordBox.Window())

	form.AddFieldRows(r,
		codeplug.FtCiChannelMode,
		codeplug.FtCiBandwidth,
		codeplug.FtCiScanList_md380,
		codeplug.FtCiSquelch,
		codeplug.FtCiRxRefFrequency,
		codeplug.FtCiTxRefFrequency,
		codeplug.FtCiTot,
		codeplug.FtCiTotRekeyDelay,
		codeplug.FtCiPower)

	column = groupBox.AddVbox()
	form = column.AddForm()

	form.AddFieldRows(r,
		codeplug.FtCiName,
		codeplug.FtCiRxFrequency,
		codeplug.FtCiTxFrequency,
		codeplug.FtCiAdmitCriteria,
		codeplug.FtCiAutoscan,
		codeplug.FtCiRxOnly,
		codeplug.FtCiLoneWorker,
		codeplug.FtCiVox,
		codeplug.FtCiAllowTalkaround)

	if settings.displayGPS {
		form.AddFieldRows(r,
			codeplug.FtCiSendGPSInfo,
			codeplug.FtCiReceiveGPSInfo)
	}

	column = row.AddVbox()
	groupBox = column.AddGroupbox("Digital Data")
	column = groupBox.AddVbox()
	form = column.AddForm()

	form.AddFieldRows(r,
		codeplug.FtCiPrivateCallConfirmed,
		codeplug.FtCiEmergencyAlarmAck,
		codeplug.FtCiDataCallConfirmed,
		codeplug.FtCiCompressedUdpDataHeader,
		codeplug.FtCiEmergencySystem,
		codeplug.FtCiContactName,
		codeplug.FtCiGroupList,
		codeplug.FtCiColorCode,
		codeplug.FtCiRepeaterSlot,
		codeplug.FtCiPrivacy,
		codeplug.FtCiPrivacyNumber,
		codeplug.FtCiInCallCriteria)

	if settings.displayGPS {
		form.AddFieldRows(r, codeplug.FtCiGPSSystem)
	}

	row = mainBox.AddHbox()
	groupBox = row.AddGroupbox("Analog Data")
	row = groupBox.AddHbox()
	column = row.AddVbox()
	form = column.AddForm()

	form.AddFieldRows(r,
		codeplug.FtCiCtcssDecode,
		codeplug.FtCiQtReverse,
		codeplug.FtCiRxSignallingSystem,
		codeplug.FtCiDisplayPTTID)

	column = row.AddVbox()
	form = column.AddForm()

	form.AddFieldRows(r,
		codeplug.FtCiCtcssEncode,
		codeplug.FtCiTxSignallingSystem,
		codeplug.FtCiReverseBurst)

	column = row.AddVbox()
	form = column.AddForm()

	form.AddFieldRows(r,
		codeplug.FtCiDecode1,
		codeplug.FtCiDecode2,
		codeplug.FtCiDecode3,
		codeplug.FtCiDecode4)

	column = row.AddVbox()
	form = column.AddForm()

	form.AddFieldRows(r,
		codeplug.FtCiDecode5,
		codeplug.FtCiDecode6,
		codeplug.FtCiDecode7,
		codeplug.FtCiDecode8)
}
