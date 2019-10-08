#!/bin/sh -l
# Add the path to the profile running this sh
PATH=$PATH:/bin:/usr/bin:/go/bin
export PATH
printenv
# Run the desired action script
echo "Running $1 Pipeline Script"
/tmp/$1.sh