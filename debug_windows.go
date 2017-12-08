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
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"unsafe"
)

func display(title string, args ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			dprint("display panicked", r)
			os.Exit(1)
		}
	}()

	var mod = syscall.NewLazyDLL("user32.dll")
	var proc = mod.NewProc("MessageBoxW")
	var MB_OK = 0x00000000

	str := " "
	for _, arg := range args {
		str += fmt.Sprintf("%#v\n", arg)
	}

	_, _, _ = proc.Call(0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(str))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(MB_OK))
}

func dprint(v ...interface{}) {
	skip := 1

	_, filename, line, ok := runtime.Caller(skip)

	str := ""
	if ok {
		dir := filepath.Base(filepath.Dir(filename))
		filename = filepath.Base(filename)
		filename = filepath.Join(dir, filename)
		str = fmt.Sprintf("%s:%d", filename, line)
	}

	v = append([]interface{}{str}, v...)

	display("Debug Messsage", v...)
}

func logFatalf(s string, v ...interface{}) {
	display("Fatal error", fmt.Sprintf(s, v...))
	os.Exit(1)
}

func logFatal(v ...interface{}) {
	display("Fatal error", fmt.Sprint(v...))
	os.Exit(1)
}

func logPrintf(s string, v ...interface{}) {
	display("Information", fmt.Sprintf(s, v...))
}

func logPrint(v ...interface{}) {
	display("Information", fmt.Sprint(v...))
}
