#!/bin/bash
# Creates the README.md file based on the utilities built in help.

echo "#  Acorn System 3 Flex Disk Utility"
echo "This utility allows for a simple editing of Acorn System 35, 40 and 80 track single sided Flex disk images. The utility
is designed for use with the Acorn System 3 6809 based Flex computer.
"
echo
echo "## Usage"
echo
echo '```text'
./bin/linux-amd64/aflex-du
echo '```'
echo
echo "## Add Binary File"
echo
echo '```text'
./bin/linux-amd64/aflex-du add-binary-file -h
echo '```'
echo
echo "## Catalogue"
echo
echo '```text'
./bin/linux-amd64/aflex-du cat -h
echo '```'
echo
echo "## Delete File"
echo
echo '```text'
./bin/linux-amd64/aflex-du delete-file -h
echo '```'
echo
echo "## Get File"
echo
echo '```text'
./bin/linux-amd64/aflex-du get-file -h
echo '```'
echo
echo "## New Disk"
echo
echo '```text'
./bin/linux-amd64/aflex-du new-disk -h
echo '```'
echo
echo "## Validate"
echo
echo '```text'
./bin/linux-amd64/aflex-du validate -h
echo '```'
echo
echo "## Version"
echo
echo '```text'
./bin/linux-amd64/aflex-du version -h
echo '```'
echo
