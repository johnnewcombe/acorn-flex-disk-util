package disk

import (
	"fmt"
)

type DirectoryEntry struct {
	Name           string // 8 chars
	Ext            string // 3 chars
	StartTrack     int    // 1 byte
	StartSector    int    // 1 byte
	EndTrack       int    // 1 byte
	EndSector      int    // 1 byte
	TotalSectors   int    // 2 bytes
	RandomFileFlag int    // 1 byte
	Month          int    // 1 byte
	Day            int    // 1 byte
	Year           int    // 1 byte
}

func (de *DirectoryEntry) Load(data []byte) error {

	// the data parameter contains a sector of data, the first 16 bytes are ignored

	de.Name = string(data[:8])
	de.Ext = string(data[8:11])
	// two bytes not used
	de.StartTrack = int(data[13])
	de.StartSector = int(data[14])
	de.EndTrack = int(data[15])
	de.EndSector = int(data[16])
	de.TotalSectors = BytesToInt(data[17:19], true)
	de.RandomFileFlag = int(data[19])
	// one byte not used
	de.Month = int(data[21])
	de.Day = int(data[22])
	de.Year = int(data[23])

	//   for each directory entry
	//     create directory entry

	//dir.DirectoryEntries = make([]DirectoryEntry, globals.DirEntriesPerSector*globals.DirectorySectors)

	return nil
}

// ToByte Each directory entry is 24 bytes long and there are 10 in a sector
func (de *DirectoryEntry) ToByte() []byte {
	var (
		dirEntry []byte
	)
	// each dir entry is 24 bytes long
	// TODO: what if strings are empty, can we use a sprintf to pad etc
	//   What chars are OK here? handle conversions here not in NewDirectoryEntry
	//dirEntry = append(dirEntry, make([]byte, 16)...) // unused bytes

	//name := fmt.Sprintf("%-8s", strings.ToUpper(de.Name))

	name := PadOrTrim([]byte(de.Name), 8)
	ext := PadOrTrim([]byte(de.Ext), 3)

	dirEntry = append(dirEntry, name...)                                    // 8 bytes
	dirEntry = append(dirEntry, ext...)                                     // 3 bytes
	dirEntry = append(dirEntry, []byte{0x00, 0x00}...)                      // 2 bytes
	dirEntry = append(dirEntry, IntToBytes(de.StartTrack, 1, false)...)     // 1 byte
	dirEntry = append(dirEntry, IntToBytes(de.StartSector, 1, false)...)    // 1 byte
	dirEntry = append(dirEntry, IntToBytes(de.EndTrack, 1, false)...)       // 1 byte
	dirEntry = append(dirEntry, IntToBytes(de.EndSector, 1, false)...)      // 1 byte
	dirEntry = append(dirEntry, IntToBytes(de.TotalSectors, 2, true)...)    // 2 bytes
	dirEntry = append(dirEntry, IntToBytes(de.RandomFileFlag, 1, false)...) // 1 byte
	dirEntry = append(dirEntry, 0)                                          // 1 byte (unused)
	dirEntry = append(dirEntry, IntToBytes(de.Month, 1, false)...)          // 1 byte
	dirEntry = append(dirEntry, IntToBytes(de.Day, 1, false)...)            // 1 byte
	dirEntry = append(dirEntry, IntToBytes(de.Year, 1, false)...)           // 1 byte

	return dirEntry
}

func (de *DirectoryEntry) FullName() string {
	return fmt.Sprintf("%s.%s", CleanString(de.Name), CleanString(de.Ext))
}
