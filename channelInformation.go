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

func channelInformation(edt *editor) {
	edt.recordWindow(codeplug.RtChannelInformation, ciRecord)
}

func ciRecord(edt *editor, recordBox *ui.HBox) {
	mainBox := recordBox.AddVbox()
	mainBox.SetContentsMargins(0, 0, 0, 0)
	row := mainBox.AddHbox()
	column := row.AddVbox()

	groupBox := column.AddGroupbox("Digital/Analog Data")
	column = groupBox.AddVbox()
	form := column.AddForm()

	r := currentRecord(recordBox.Window())

	form.AddFieldRows(r,
		codeplug.FtChannelMode,
		codeplug.FtBandwidth,
		codeplug.FtScanList,
		codeplug.FtSquelch,
		codeplug.FtRxRefFrequency,
		codeplug.FtTxRefFrequency,
		codeplug.FtTot,
		codeplug.FtTotRekeyDelay,
		codeplug.FtPower)

	column = groupBox.AddVbox()
	form = column.AddForm()

	form.AddFieldRows(r,
		codeplug.FtChannelName,
		codeplug.FtRxFrequency,
		codeplug.FtTxFrequency,
		codeplug.FtAdmitCriteria,
		codeplug.FtAutoscan,
		codeplug.FtRxOnly,
		codeplug.FtLoneWorker,
		codeplug.FtVox,
		codeplug.FtAllowTalkaround)

	column = row.AddVbox()
	groupBox = column.AddGroupbox("Digital Data")
	column = groupBox.AddVbox()
	form = column.AddForm()

	form.AddFieldRows(r,
		codeplug.FtPrivateCallConfirmed,
		codeplug.FtEmergencyAlarmAck,
		codeplug.FtDataCallConfirmed,
		codeplug.FtCompressedUdpDataHeader,
		codeplug.FtContactName,
		codeplug.FtGroupList,
		codeplug.FtColorCode,
		codeplug.FtRepeaterSlot,
		codeplug.FtPrivacy,
		codeplug.FtPrivacyNumber)

	row = mainBox.AddHbox()
	groupBox = row.AddGroupbox("Analog Data")
	row = groupBox.AddHbox()
	column = row.AddVbox()
	form = column.AddForm()

	form.AddFieldRows(r,
		codeplug.FtCtcssDecode,
		codeplug.FtQtReverse,
		codeplug.FtRxSignallingSystem,
		codeplug.FtDisplayPTTID)

	column = row.AddVbox()
	form = column.AddForm()

	form.AddFieldRows(r,
		codeplug.FtCtcssEncode,
		codeplug.FtTxSignallingSystem,
		codeplug.FtReverseBurst)

	column = row.AddVbox()
	form = column.AddForm()

	form.AddFieldRows(r,
		codeplug.FtDecode1,
		codeplug.FtDecode2,
		codeplug.FtDecode3,
		codeplug.FtDecode4)

	column = row.AddVbox()
	form = column.AddForm()

	form.AddFieldRows(r,
		codeplug.FtDecode5,
		codeplug.FtDecode6,
		codeplug.FtDecode7,
		codeplug.FtDecode8)

	mainBox.AddFiller()
}
