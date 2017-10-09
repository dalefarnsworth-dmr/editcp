#!/bin/bash
dirname=xyzzy
appname=editcp

if [ ! -f "$dirname/$appname" -o ! -d "$dirname/qml" ]; then
	echo "$dirname/$appname not found." 1>&2
	echo "cd to the $appname installation directory and run ./install"
	exit 1
fi

export LD_LIBRARY_PATH="$dirname"/lib
export QT_PLUGIN_PATH="$dirname"/plugins
export QML_IMPORT_PATH="$dirname"/qml
export QML2_IMPORT_PATH="$dirname"/qml
echo "$dirname/$appname" "$@"
"$dirname/$appname" "$@"
