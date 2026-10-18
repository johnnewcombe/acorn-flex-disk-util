package disk

import "github.com/johnnewcombe/System3-disk-util/globals"

/*
DataSector

Every sector on a FLEX diskette (except the two BOOT has the following
format:

	Bytes 0-1 Link to the next sector
	Bytes 2-3 File Logical Record Number
	Bytes 4-255 Data sectors)

If a file occupies more than one sector, the "link to the next sector"
portion contains the track and sector numbers, respectively, of the next
sector in the file. These bytes are zero in the last sector of a file,
indicating that no more data follows (an "end-of-file" condition).
*/

type DataSector struct {
	NextTrack            int    // 1 byte
	NextSector           int    // 1 byte
	SectorSequenceNumber int    // 2 bytes
	Data                 []byte //252 bytes
}

func (d *DataSector) Load(data []byte) error {

	d.NextTrack = int(data[0])
	d.NextSector = int(data[1])
	d.SectorSequenceNumber = BytesToInt(data[2:], true)
	d.Data = data[4:]

	return nil
}

func (d *DataSector) ToByte() ([]byte, error) {
	var (
		sector []byte
	)
	sector = append(sector, byte(d.NextTrack))
	sector = append(sector, byte(d.NextSector))
	sector = append(sector, IntToBytes(d.SectorSequenceNumber, 2, true)...)

	// append sector
	return append(sector, PadOrTrim(d.Data, globals.BytesPerSector-4)...), nil
	//return append(sector, d.Data...), nil
}

// NewDataSector returned a new sector that is linked to the specified one.
// On these disk images the sector numbering starts at 1. i.e. a zero value
// means no link.
func NewDataSector(linkedSector int) (*DataSector, error) {

	ds := &DataSector{}

	// get the link adress for sector
	ds.NextTrack, ds.NextSector = getSectorAddress(linkedSector)
	ds.SectorSequenceNumber = 0
	ds.Data = createEmptySector(0x00)[:globals.BytesPerSector-4]

	return ds, nil
}
