# Acorn System 3 Flex Disk Utility

This utility allows for a simple editing of Acorn System 35, 40 and 80 track single sided Flex disk images. The utility
is designed for use with the Acorn System 3 6809 based Flex computer.


    Usage:
    disk-util [command]
    
    Available Commands:
    add-file    Adds a binary file to the specified Acorn System single sided 40 or 80 track disk image.
    create-disk Creates a blank Acorn System single sided 40 or 80 track disk image.
    help        Help about any command
    version     Displays the program version.
    
    Flags:
    -h, --help   help for disk-util
    
    Use "disk-util [command] --help" for more information about a command.

## Create Disk

Creates a blank Acorn System single sided 40 or 80 track disk image.

    Usage:
    disk-util create-disk [flags]
    
    Flags:
    -b, --boot-option      int     0: Disabled, 1: *LOAD !BOOT, 2: *RUN !BOOT, 3: *EXEC !BOOT.
    -d, --disk-name        string  Name of the disk image file (*.img) to be created. (default "mydisk.ssd")
        --title            string  Title of the disk as displayed in the catalogue.
    -t, --tracks           int     Number of tracks for the single sided disk (40 or 80). (default 40)

The following example creates a 40 track disk image called _mydisk.ssd_ with a title of "Utils".

    disk-util create-disk -d mydisk.ssd --title Utils

## Add File

Adds a binary file to the specified Acorn System single sided 40 or 80 track disk image.

    Usage:
    disk-util add-file [flags]
    
    Flags:
    -d, --disk-name        string  Name of the disk image file (*.img) to be created.
    -e, --execute-address  string  The LOAD address associated with the file.
    -f, --filename         string  The name of a file to add to the disk during the creation process.
    -l, --load-address     string  The EXECUTION address associated with the file, defaults to LOAD address. (default "2800")

The following example adds a file called _somepath/basic.bin_ with a title of "Utils".

    disk-util add-file -d mydisk.ssd -f somepath/basic.bin -l c000 -e c2b2

### A note about filenames...

To keep things simple, the destination filename is derived from the specified filename. It will be an upper case version
of the source name, however, any path will be removed as will the extension. Therefore, in the above example
_somepath/basic.bin_ becomes _BASIC_ on the disk image. This means that the main element of the source file name i.e.
_basic_ must conform to the filename spec of Acorn System DOS and be 7 characters or less.

## Info

Displays disk and catalogue information.

    Usage:
    disk-util info [flags]
    
    Flags:
    -d, --disk-name string   Name of the disk image file (*.img) to be created.

The following command displays the contents of the disk image _./output/prestel.ssd_.

    disk-util info -d prestel.ssd

The output will look something like the following.

    --------------------------------------------------------------------
    Disk Image: ./output/prestel.ssd
    --------------------------------------------------------------------
    Title: PRESTEL     
    Tracks: 40
    Boot Option: 0
    File Offset: 16
    Sectors in Volume: 400
    Files: 2
    --------------------------------------------------------------------
    00 PRESTEL 2800 2800 0AD7 002
    01 BASIC   C000 C2B2 1000 00D
    --------------------------------------------------------------------

The values shown following the filename represent the load address, execution address, length and start sector
respectively.