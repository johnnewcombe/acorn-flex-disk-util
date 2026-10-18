package disk

import (
	"slices"

	"github.com/johnnewcombe/System3-disk-util/globals"
)

type DirectorySector struct {
	NextTrack        int
	NextSector       int
	DirectoryEntries []*DirectoryEntry
}

func (ds *DirectorySector) Load(sector []byte) {

	ds.NextTrack = int(sector[0])
	ds.NextSector = int(sector[1])

	// 10 dir entries per sector but the first 16 bytes are ignored
	for dirEntryData := range slices.Chunk(sector[16:], globals.BytesPerDirEntry) {
		de := new(DirectoryEntry)
		de.Load(dirEntryData)
		ds.DirectoryEntries = append(ds.DirectoryEntries, de)
	}
}

func (ds *DirectorySector) ToByte() []byte {

	var (
		sector []byte
	)

	sector = make([]byte, 16)
	sector[0] = byte(ds.NextTrack)
	sector[1] = byte(ds.NextSector)

	for _, de := range ds.DirectoryEntries {
		sector = append(sector, de.ToByte()...)
	}
	return sector
}

func NewDirectorySector(linkedSector int) (*DirectorySector, error) {

	ds := &DirectorySector{}

	// get the link address for sector
	ds.NextTrack, ds.NextSector = getSectorAddress(linkedSector)

	for i := 0; i < globals.DirEntriesPerSector; i++ {
		ds.DirectoryEntries = append(ds.DirectoryEntries, &DirectoryEntry{})
	}
	return ds, nil
}
