#!/bin/bash

GOPICAM_SCRIPTS_FOLDER=$( dirname "${BASH_SOURCE[0]}" )

cd $GOPICAM_SCRIPTS_FOLDER
cd ../html/scss

for SCSSFILE in `ls *.scss`;
do
	echo "transpile " $SCSSFILE

	CSSFILE=$(echo $SCSSFILE | sed -e 's/\.scss$/\.css/')

	sassc -I partials/ $SCSSFILE > ../css/$CSSFILE
done
