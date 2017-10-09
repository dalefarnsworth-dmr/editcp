#!/bin/bash
dirname=/home/dale/go/src/github.com/dalefarnsworth/codeplug/editcp/deploy/linux
appname=editcp

install() {
	if [ $# -gt 0 ]; then
		dir="$1"
		if [ "$dir" == . ]; then
			exit 0
		fi

		echo "+ ln -sf" "$dirname/$appname.sh" "$dir/$appname"
		ln -sf "$dirname/$appname.sh" "$dir/$appname"
		exit $?
	fi

	i=0
	while read dir; do
		if [ -d "$dir" ]; then
			dirs[$i]="$dir"
			((i++))
		fi
	done <<-END
		$HOME/bin
		/usr/local/bin
		/usr/bin
		$dirname/bin
	END

	while :; do
		echo "Select the directory where a link to $appname will be created."
		PS3="Select a directory number "
		select installdir in "${dirs[@]}" Other Quit; do 
			case "$REPLY" in
			[1-${#dirs[@]}])
				break;;

			$((${#dirs[@]}+1)))
				dir=
				while [ "$dir" = "" ]; do
					echo -n "Enter a directory name: "
					read dir
					[ -z "$dir" ] && continue 3
					[ -d "$dir" ] && break
					echo "$dir: not found" 1>&2
					continue 3
				done
				installdir="$dir"
				break;;

			$((${#dirs[@]}+2)))
				exit;;

			*)
				if [ -z "$installdir" ]; then
					if [ -d "$REPLY" ]; then
						installdir="$REPLY"
						break
					fi
					echo "Invalid selection."
					continue 2
				fi;;
			esac
		done
		break
	done

	install "$installdir"
}

if [ ! -f "./$appname.sh" ]; then
	echo "./$appname.sh is not correctly installed." 1>&2
	echo "cd to the $appname installation directory and run ./$appname.sh --install" 1>&2
	exit 1
fi

dirname=$(pwd)
sed --in-place -e "/^dirname=/cdirname=$dirname" $appname.sh
install "$@"
