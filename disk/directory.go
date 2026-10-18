package disk

import (
	"github.com/johnnewcombe/System3-disk-util/globals"
)

/*
DirectoryEntry

Each sector in the directory portion of a FLEX diskette contains 10
directory entries. Each entry refers to one file on the diskette. In
each sector, the first four bytes contain the sector linkage information
and the next 12 bytes are not used. When reading information from the
directory using the FMS Get Information Record function, these 16 bytes
are skipped automatically as each sector is read; the user need not be
concerned with them.

Each entry in the directory contains the exact same information that is
stored in the FCB bytes 4-27. See the description of the File Control
Block (FCB) for more details.

A directory entry which has never been used has a zero in the first byte
of the file name. A directory entry which has benn deleted has the
leftmost bit of the name set (i.e. the first byte of the name is negative).

The sectors in the diskette directory on track 0 also have linkage bytes
similar to those in the in the free chain and data files.
*/

type Directory struct {
	DirectorySectors []*DirectorySector
}

// ToByte The directory is 16 sectors in size so 16 * 256 bytes.
func (dir *Directory) ToByte() []byte {
	var (
		directory []byte
	)

	// run through all of the directory sectors adding 10 directory entries to each
	// note that we can't just iterate over the directory entries as each sector has
	// 16 empty bytes at the beginning
	for i := 0; i < globals.DirectorySectors; i++ {

		// add the 16 byte padding at the beginig or each sector
		//padding = make([]byte, 16)
		//directory = append(directory, padding...)
		directory = append(directory, dir.DirectorySectors[i].ToByte()...)
	}

	return directory
}

/*
// getNextFreeDirectoryEntry returns the directory sector id and directory entry id
// of the next free directory entry. If there are no free directories, zero is
// returned for both values. See
func (dir *Directory) getNextFreeDirectoryEntry() (int, int) {
	var ()

	for sectorId, ds := range dir.DirectorySectors {
		for dirEntryId, de := range ds.DirectoryEntries {
			if de.Name == "" || de.Name[0] == 0xFF || de.Name[0] == 0x00 {
				return sectorId + globals.DirectoryStartSector, dirEntryId
			}
		}
	}
	return 0, 0
}
*/

// getNextFreeDirectoryEntry returns the directory sector id and entry id of the
// next free directory entry, and false if the directory is full.
func (dir *Directory) getNextFreeDirectoryEntry() (int, int, bool) {
	for sectorId, ds := range dir.DirectorySectors {
		for dirEntryId, de := range ds.DirectoryEntries {
			if de.Name == "" || de.Name[0] == 0x00 || de.Name[0]&0x80 != 0 {
				return sectorId + globals.DirectoryStartSector, dirEntryId, true
			}
		}
	}
	return 0, 0, false
}

// NewDirectory Creates a new empty Dirctory with 160 empty directory entries.
func NewDirectory() (*Directory, error) {

	var (
		dir          *Directory
		ds           *DirectorySector
		linkedSector int
		err          error
	)

	dir = &Directory{}

	for i := 0; i < globals.DirectorySectors; i++ {

		// Link the directory sectors (directory starts at sectorId 4
		// but we link to the next sector

		if i == globals.DirectorySectors-1 {
			linkedSector = 0
		} else {
			linkedSector = i + 5
		}

		if ds, err = NewDirectorySector(linkedSector); err != nil {
			return nil, err
		}

		for j := 0; j < globals.DirEntriesPerSector; j++ {
			de := &DirectoryEntry{}
			ds.DirectoryEntries[j] = de
		}
		dir.DirectorySectors = append(dir.DirectorySectors, ds)
	}

	// create 16 * 10 directory entries 16 sectors of 10 dir entries each
	return dir, nil
}
