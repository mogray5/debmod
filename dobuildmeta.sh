cd $1
equivs-build ns-control

if ls *.deb 1> /dev/null 2>&1; then
	reprepro -b $2 includedeb $3 *.deb
fi