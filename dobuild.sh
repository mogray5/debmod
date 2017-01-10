#!/bin/bash

cd $1
debuild -us -uc -b
cd ..
if ls *.deb 1> /dev/null 2>&1; then
	reprepro -b $2 include $3 *.changes
fi
