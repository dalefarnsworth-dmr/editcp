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

func zones(edt *editor) {
	writable := true
	edt.recordWindow(codeplug.RtZones_md380, writable, ziRecord)
}

func ziRecord(edt *editor, recordBox *ui.HBox) {
	column := recordBox.AddVbox()
	addFieldMembers(column, &settings.sortAvailableChannels,
		codeplug.FtZiName, codeplug.FtZiChannel_md380, "Channels")
}
