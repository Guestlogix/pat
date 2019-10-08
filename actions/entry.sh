#!/bin/sh -l
# Add the path to the profile running this sh
PATH=$PATH:/bin:/usr/bin:/go/bin
export PATH
# Run the desired action script
echo "Seeing if PAT works"
pat --help

echo "Checking current dir"
pwd
ls -al
echo "Checking up a dir"
cd ..
pwd
ls -al