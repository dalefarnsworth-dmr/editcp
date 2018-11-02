#!/bin/bash
dirname=
appname=editcp

if [ ! -f "$dirname/$appname" ]; then
	echo "$dirname/$appname not found." 1>&2
	echo "cd to the $appname installation directory and run ./install"
	exit 1
fi

export LD_LIBRARY_PATH="$dirname"/lib
export QT_PLUGIN_PATH="$dirname"/plugins
export QML_IMPORT_PATH="$dirname"/qml
export QML2_IMPORT_PATH="$dirname"/qml
export QT_LOGGING_RULES="*=false"
"$dirname/$appname" "$@"
