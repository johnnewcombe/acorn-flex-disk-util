#  Acorn System 3 Flex Disk Utility
This utility allows for a simple editing of Acorn System 35, 40 and 80 track single sided Flex disk images. The utility
is designed for use with the Acorn System 3 6809 based Flex computer.


## Usage

```text
Acorn System DOS disk Utility for Acorn System DOS. Note that this utility only supports single sided (.ssd) disk images.

Usage:
  disk-util [command]

Available Commands:
  add-binary-file Adds a raw or Flex binary file to the specified Acorn Flex single sided 35, 40 or 80 track disk image.
  cat             Displays disk and catalogue information.
  completion      Generate the autocompletion script for the specified shell
  delete-file     Deletes a file from the specified Acorn Flex single sided 35, 40 or 80 track disk image.
  get-file        Gets a file from the specified Acorn Flex single sided 35, 40 or 80 track disk image.
  help            Help about any command
  new-disk        Creates a blank Acorn Flex single sided 35, 40 or 80 track disk image.
  validate        Check the structure of a FLEX disk image
  validate        Check the structure of a FLEX disk image
  version         Displays the program version.

Flags:
  -h, --help   help for disk-util

Use "disk-util [command] --help" for more information about a command.
```

## Add Binary File

```text
Adds a binary file to the specified Acorn Flex single sided 35, 40 or 80 track disk image. 
raw binary files are automatically converted to Flex CMD format files.

Usage:
  disk-util add-binary-file [flags]

Flags:
  -d, --disk-image-filename string   Name of the disk image file (*.img) to be used.
  -f, --filename string              The name of the local file to be processed.
  -h, --help                         help for add-binary-file
  -l, --load-address uint16          The address where the file should be loaded into memory. (default 256)
  -t, --transfer-address uint16      The execute address. (default 65535)
```

## Catalogue

```text
Displays disk and catalogue information.

Usage:
  disk-util cat [flags]

Flags:
  -d, --disk-image-filename string   Name of the disk image file (*.img) to be used.
  -h, --help                         help for cat
```

## Delete File

```text
Deletes a file from the specified Acorn Flex single sided 35, 40 or 80 track disk image.

Usage:
  disk-util delete-file [flags]

Flags:
  -d, --disk-image-filename string   Name of the disk image file (*.img) to be used.
  -f, --filename string              The name of the flex filename to process.
  -h, --help                         help for delete-file
```

## Get File

```text
Deletes a file from the specified Acorn Flex single sided 35, 40 or 80 track disk image.

Usage:
  disk-util get-file [flags]

Flags:
  -d, --disk-image-filename string   Name of the disk image file (*.img) to be used.
  -f, --filename string              The name of the flex filename to process.
  -h, --help                         help for get-file
```

## New Disk

```text
Creates a blank Acorn Flex single sided 35, 40 or 40 or 80 track disk image.

Usage:
  disk-util new-disk [flags]

Flags:
  -d, --disk-image-filename string   Name of the disk image file (*.img) to be used. (default "mydisk.img")
  -h, --help                         help for new-disk
  -t, --tracks int                   Number of tracks for the single sided disk (35, 40 or 80). (default 80)
  -v, --volume-label string          Volume Label of the disk as displayed in the catalogue. (default "ACORNFLEX")
  -n, --volume-number int            Volume Number Nof the disk as displayed in the catalogue. (default 1)
```

## Validate

```text
Checks a FLEX disk image for structural damage:

  - the image size and the geometry recorded in the system information record
  - the directory and every file's sector chain (links, record numbers,
    sector counts and end sectors against the directory)
  - the free chain against the system information record
  - sectors belonging to more than one chain, or to none
  - that .CMD and .BIN files contain valid FLEX binary records

Errors mean the disk is damaged and should not be written to. Warnings are
suspicious but not necessarily damage; use --strict to treat them as errors.
The exit status is non-zero if the disk fails.

Usage:
  disk-util validate <disk-image> [flags]

Aliases:
  validate, verify, check

Flags:
  -d, --disk-image-filename string   Name of the disk image file (*.img) to be used.
  -h, --help                         help for validate
      --strict                       treat warnings as errors
```

## Version

```text
Displays the program version.

Usage:
  disk-util version [flags]

Flags:
  -h, --help   help for version
```

