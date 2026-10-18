package disk

/*  Notes:

 */

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/johnnewcombe/System3-disk-util/globals"
)

/*
Flex Disk Initialisation

The first track on the diskette, track 0, is special.
None of the sectors on track 0 are available for data files, they are
reserved for use by the FLEX system. The first two sectors contain a
"boot" program which is loaded by the "D" command of the SBUG monitor
or by whatever comparable ROM based bootstrap is in use.

The boot program, once loaded, then loads FLEX from the diskette.
Another sector on track 0 is the System Information Record. This
sector contains the track and sector addresses of the beginning and
ending sectors of the chain of free sectors, those available for data
files. The rest of track 0 is used for the directory of file names.

After initialisation, the free tracks on the diskette have a common
format. The first two bytes of each sector contain the track and sector
number of the next sector in the chain. The next two bytes are used to
store the logical record number of the sector in the file. The
remaining 252 bytes are zero. Initially, all record number bytes are
zero. When data is stored in a file, the two linkage bytes at the
beginning of each sector are modified to point to the next sector in the
file, not the next the free chain.
*/

const (
	flexDataPerSector = globals.BytesPerSector - 4 // 252 data bytes after the link and record number

	utilitySpaceStart = 0xC100 // FLEX utility command space
	utilitySpaceEnd   = 0xC6FF
)

// supported disk sizes: 35, 40 and 80 tracks of 10 single density 256 byte sectors
var supportedTrackCounts = []int{35, 40, 80}

// diskSizeForTracks returns the image size in bytes for a disk with the given number of tracks.
func diskSizeForTracks(tracks int) int {
	return tracks * globals.SectorsPerTrack * globals.BytesPerSector
}

// tracksForDiskSize returns the number of tracks for an image of the given size,
// or an error if the size is not one of the supported disk sizes.
func tracksForDiskSize(size int) (int, error) {
	for _, tracks := range supportedTrackCounts {
		if size == diskSizeForTracks(tracks) {
			return tracks, nil
		}
	}
	return 0, fmt.Errorf("disk image is %d bytes; must be %d, %d or %d bytes (35, 40 or 80 tracks)",
		size, diskSizeForTracks(35), diskSizeForTracks(40), diskSizeForTracks(80))
}

type Disk struct {
	BootSector1 []byte     // sect 1
	BootSector2 []byte     // sect 2
	SystemInfo  *SIR       //sect 3
	Directory   *Directory // sect 5-20 (16 sects)
	DataSectors []*DataSector
}

func (d *Disk) Load(diskImage []byte) error {

	tracks, err := tracksForDiskSize(len(diskImage))
	if err != nil {
		return err
	}

	d.Directory = new(Directory)

	// directory sectors
	index := 0
	for sector := range slices.Chunk(diskImage, globals.BytesPerSector) {

		switch index {
		case 0:
			d.BootSector1 = sector
		case 1:
			d.BootSector2 = sector
		case 2:
			d.SystemInfo = new(SIR)
			if err = d.SystemInfo.Load(sector); err != nil {
				return fmt.Errorf("unable to read system information record: %w", err)
			}
		case 3: //unused sector
		case 4, 5, 6, 7, 8, 9:
			// directory sectors
			directorySector := new(DirectorySector)
			directorySector.Load(sector)
			d.Directory.DirectorySectors = append(d.Directory.DirectorySectors, directorySector)
		default:
			//data sectors
			// get linked sector id
			datasector := new(DataSector)
			if err = datasector.Load(sector); err != nil {
				return fmt.Errorf("unable to read data sector: %w", err)
			}
			d.DataSectors = append(d.DataSectors, datasector)
		}

		index++
	}

	// the geometry recorded in the SIR must agree with the size of the image
	if d.SystemInfo.EndTrack != tracks-1 || d.SystemInfo.EndSector != globals.SectorsPerTrack {
		return fmt.Errorf("system information record says the disk ends at T%02d-S%02d, "+
			"but a %d byte image is %d tracks of %d sectors (ends at T%02d-S%02d)",
			d.SystemInfo.EndTrack, d.SystemInfo.EndSector, len(diskImage),
			tracks, globals.SectorsPerTrack, tracks-1, globals.SectorsPerTrack)
	}

	return nil

}

func (d *Disk) LinkBootSector(filename string) error {

	// find directory entry for file
	// put sector address into sector 0
	for _, ds := range d.Directory.DirectorySectors {
		for _, de := range ds.DirectoryEntries {
			if strings.EqualFold(de.FullName(), filename) {
				d.BootSector1[5] = byte(de.StartTrack)
				d.BootSector1[6] = byte(de.StartSector)
			}
		}
	}
	return nil
}

// AddBinaryFile adds a FLEX binary (command) file to the disk image.
//
// fileData may be a raw binary image, which is converted to FLEX binary format
// using loadAddress and (if hasTransfer) transferAddress, or data that is
// already in FLEX binary format, which is stored unchanged (its own records
// supply the addresses, so loadAddress and transferAddress are ignored).
//
// All checks are made before anything is modified, so if an error is returned
// the disk image is left exactly as it was.
func (d *Disk) AddBinaryFile(filename string, fileData []byte, loadAddress uint16,
	transferAddress uint16, hasTransfer bool) error {

	const recordSize = 255

	var (
		flexFileData []byte
		err          error
	)

	// validate the name
	name, ext, err := splitFlexName(filename)
	if err != nil {
		return err
	}
	if d.FileExists(name + "." + ext) {
		return fmt.Errorf("file %s.%s already exists, delete old file first", name, ext)
	}

	// produce FLEX binary data
	if len(fileData) == 0 {
		return errors.New("file is empty")
	}

	if IsFlexBinary(fileData) {
		// TODO perhaps add an addCommandFile command without the load and transfer stuff and raise an error here?
		//  we could use a -force-convert flag in cases where parsing detected a flex file but it wasn't (impossible/rare)
		// already a flex file so don't convert and ignore the load and transfer stuff
		flexFileData = fileData
	} else {
		// convert to flex binary format
		if flexFileData, err = BinToFlex(fileData, loadAddress, transferAddress, hasTransfer, recordSize); err != nil {
			return err
		}
	}

	// make sure nothing overruns the utility command space ($C100-$C6FF)
	segments, _, err := ParseFlexBinary(flexFileData)
	if err != nil {
		return fmt.Errorf("invalid FLEX binary: %w", err)
	}
	for _, seg := range segments {
		start := int(seg.Address)
		end := start + len(seg.Data) - 1
		if start >= utilitySpaceStart && start <= utilitySpaceEnd && end > utilitySpaceEnd {
			return fmt.Errorf("segment at $%04X-$%04X overruns the utility command space ($%04X-$%04X), consider a different start address",
				start, end, utilitySpaceStart, utilitySpaceEnd)
		}
	}

	// work out the space needed and check it is available
	// round up: a partly filled last sector still needs a whole sector
	sectorsNeeded := (len(flexFileData) + flexDataPerSector - 1) / flexDataPerSector

	if sectorsNeeded > d.SystemInfo.NumberFreeSectors {
		return fmt.Errorf("disk is full: %d sectors needed, %d free", sectorsNeeded, d.SystemInfo.NumberFreeSectors)
	}

	dirSectorId, entryId, ok := d.Directory.getNextFreeDirectoryEntry()
	if !ok {
		return errors.New("directory is full")
	}
	dirIndex := dirSectorId - globals.DirectoryStartSector

	// walk the free chain and reserve the sectors (nothing is changed yet)
	allocated := make([]*DataSector, 0, sectorsNeeded)
	allocatedIds := make([]int, 0, sectorsNeeded)

	t, s := d.SystemInfo.FirstFreeTrack, d.SystemInfo.FirstFreeSector
	for i := 0; i < sectorsNeeded; i++ {
		if t == 0 && s == 0 {
			return fmt.Errorf("free chain ends after %d sectors but the system information record says %d are free",
				i, d.SystemInfo.NumberFreeSectors)
		}
		id := getSectorId(t, s)
		if id < globals.SectorsPerTrack || id-globals.SectorsPerTrack >= len(d.DataSectors) {
			return fmt.Errorf("free chain points to invalid sector T%02d-S%02d, corrupt or non-flex disk", t, s)
		}
		sector := d.DataSectors[id-globals.SectorsPerTrack] // data starts at track 1
		allocated = append(allocated, sector)
		allocatedIds = append(allocatedIds, id)
		t, s = sector.NextTrack, sector.NextSector
	}
	// t,s is now the link out of the last reserved sector, i.e. the new start of the
	// free chain. It must be captured here, before that link is set to 0,0 below.
	nextFreeTrack, nextFreeSector := t, s

	// write the data into the reserved sectors
	for i, sector := range allocated {
		start := i * flexDataPerSector
		end := min(start+flexDataPerSector, len(flexFileData))

		buf := make([]byte, flexDataPerSector) // zero padded, independent of the caller's slice
		copy(buf, flexFileData[start:end])

		sector.Data = buf
		sector.SectorSequenceNumber = i + 1
		// the links between the file's sectors are already correct: they came from the free chain
	}
	last := allocated[len(allocated)-1]
	last.NextTrack, last.NextSector = 0, 0 // end of this file's chain

	// directory entry
	firstTrack, firstSector := getSectorAddress(allocatedIds[0])
	lastTrack, lastSector := getSectorAddress(allocatedIds[len(allocatedIds)-1])
	day, month, year := GetFlexDate()

	entry := d.Directory.DirectorySectors[dirIndex].DirectoryEntries[entryId]
	entry.Name = name
	entry.Ext = ext
	entry.StartTrack, entry.StartSector = firstTrack, firstSector
	entry.EndTrack, entry.EndSector = lastTrack, lastSector
	entry.TotalSectors = sectorsNeeded
	entry.RandomFileFlag = 0 // binary files are always sequential
	entry.Month = month
	entry.Day = day
	entry.Year = year

	// system information record
	d.SystemInfo.NumberFreeSectors -= sectorsNeeded
	d.SystemInfo.FirstFreeTrack, d.SystemInfo.FirstFreeSector = nextFreeTrack, nextFreeSector
	if nextFreeTrack == 0 && nextFreeSector == 0 {
		// the last free sector has been used, so the free chain is now empty
		d.SystemInfo.LastFreeTrack, d.SystemInfo.LastFreeSector = 0, 0
	}

	return nil
}

// splitFlexName validates a host filename and returns the FLEX name (1-8
// characters) and extension (1-3 characters), upper-cased. The first character
// must be a letter, the rest letters, digits, '-' or '_'.
func splitFlexName(filename string) (string, string, error) {
	base := strings.ToUpper(filepath.Base(filename))
	ext := filepath.Ext(base)
	if ext == "" {
		return "", "", fmt.Errorf("%q has no extension (e.g. .CMD or .BIN)", filename)
	}
	name := strings.TrimSuffix(base, ext)
	ext = ext[1:]

	if len(name) < 1 || len(name) > 8 {
		return "", "", fmt.Errorf("flex filename %q must be 1-8 characters", name)
	}
	if len(ext) < 1 || len(ext) > 3 {
		return "", "", fmt.Errorf("flex extension %q must be 1-3 characters", ext)
	}
	for i, part := range []string{name, ext} {
		for j, c := range part {
			letter := c >= 'A' && c <= 'Z'
			valid := letter || (c >= '0' && c <= '9') || c == '-' || c == '_'
			if !valid || (i == 0 && j == 0 && !letter) {
				return "", "", fmt.Errorf("invalid flex filename %q", filename)
			}
		}
	}
	return name, ext, nil
}

func (d *Disk) DeleteFile(filename string) error {

	var (
		err error
	)

	if !d.FileExists(filename) {
		err = fmt.Errorf("file %s not found", filename)
		return err
	}

	// find and delete the named element
	for _, ds := range d.Directory.DirectorySectors {
		for _, de := range ds.DirectoryEntries {
			if strings.EqualFold(de.FullName(), filename) {

				// deleted files have the 7th bit set, typically 0xFF is used,
				de.Name = fmt.Sprintf("%s%s", globals.DeleteChar, de.Name[1:])

				// add the released sectors back to the system
				if d.SystemInfo.LastFreeTrack == 0 && d.SystemInfo.LastFreeSector == 0 {
					// the free chain is empty (disk was full), so the deleted file becomes the whole chain
					d.SystemInfo.FirstFreeTrack = de.StartTrack
					d.SystemInfo.FirstFreeSector = de.StartSector
				} else {
					// make the current last free track/sector point to the start track/sector of the deleted file
					lastFreeSector := getSectorId(d.SystemInfo.LastFreeTrack, d.SystemInfo.LastFreeSector) - globals.SectorsPerTrack
					d.DataSectors[lastFreeSector].NextTrack = de.StartTrack
					d.DataSectors[lastFreeSector].NextSector = de.StartSector
				}

				// then make the last free track/sector equal to the deleted files end track/sector, this sector
				d.SystemInfo.LastFreeTrack = de.EndTrack
				d.SystemInfo.LastFreeSector = de.EndSector

				// increase the number of the free sectors
				d.SystemInfo.NumberFreeSectors += de.TotalSectors

			}
		}
	}

	return nil
}

func (d *Disk) GetFile(filename string) ([]byte, error) {

	var (
		data []byte
		err  error
	)

	if !d.FileExists(filename) {
		err = fmt.Errorf("file %s not found", filename)
		return nil, err
	}

	// find and get the named element
	for _, ds := range d.Directory.DirectorySectors {
		for _, de := range ds.DirectoryEntries {
			if strings.EqualFold(de.FullName(), filename) {

				// get the start sector
				nextSectorId := getSectorId(de.StartTrack, de.StartSector)

				for {

					// get the data sector
					dataSector := d.DataSectors[nextSectorId-globals.SectorsPerTrack]

					// get the data
					data = append(data, dataSector.Data...)

					// all done? a link of 0,0 marks the last sector of the file
					if dataSector.NextTrack == 0 && dataSector.NextSector == 0 {
						break
					}
					nextSectorId = getSectorId(dataSector.NextTrack, dataSector.NextSector)
				}
			}
		}

	}

	return data, nil
}

func (d *Disk) Catalogue() (string, error) {

	var (
		info    strings.Builder
		files   int
		sectors int
		random  string
	)
	info.WriteString("\n--------------------------------------------------\n")
	info.WriteString(fmt.Sprintf("     Volume Label: %s\n", d.SystemInfo.VolumeLabel))
	info.WriteString(fmt.Sprintf("    Volume Number: %d\n", d.SystemInfo.VolumeNumber))
	info.WriteString(fmt.Sprintf(" First Free Track: %d\n", d.SystemInfo.FirstFreeTrack))
	info.WriteString(fmt.Sprintf("First Free Sector: %d\n", d.SystemInfo.FirstFreeSector))
	info.WriteString(fmt.Sprintf("  Last Free Track: %d\n", d.SystemInfo.LastFreeTrack))
	info.WriteString(fmt.Sprintf(" Last Free Sector: %d\n", d.SystemInfo.LastFreeSector))
	info.WriteString("--------------------------------------------------\n")
	info.WriteString("NAME          START    END      SIZE  DATE      R\n")
	info.WriteString("--------------------------------------------------\n")

	for _, f := range d.Directory.DirectorySectors {
		for _, j := range f.DirectoryEntries {

			if j.Name == "" || j.Name[0] == 0x00 || j.Name[0]&0x80 != 0 { // never used or deleted
				continue
			}
			files++
			sectors += j.TotalSectors

			if j.RandomFileFlag != 0 {
				random = "R"
			} else {
				random = ""
			}

			info.WriteString(fmt.Sprintf("%-12s  T%02d-S%02d  T%02d-S%02d  %4d  %s\n",
				j.FullName(),
				j.StartTrack,
				j.StartSector,
				j.EndTrack,
				j.EndSector,
				j.TotalSectors,
				fmt.Sprintf("%02d-%02d-%02d  %s", j.Day, j.Month, j.Year, random)),
			)
		}
	}
	info.WriteString("--------------------------------------------------\n")
	info.WriteString(fmt.Sprintf("       FILES: %d\n", files))
	info.WriteString(fmt.Sprintf("FREE SECTORS: %d\n", d.SystemInfo.NumberFreeSectors))

	return info.String(), nil
}

func (d *Disk) FileExists(filename string) bool {
	// does the file already exist in the directory
	for _, f := range d.Directory.DirectorySectors {
		for i := range f.DirectoryEntries {
			if strings.EqualFold(f.DirectoryEntries[i].FullName(), filename) {
				return true //errors.New("file already exists, delete old file first")
			}
		}
	}
	return false
}

func (d *Disk) ToByte() ([]byte, error) {

	var (
		disk   []byte
		sector []byte
		err    error
	)

	// add the first boot sector
	disk = append(disk, d.BootSector1...)
	disk = append(disk, d.BootSector2...)

	// add system information record (SIR) sector
	if sector, err = d.SystemInfo.ToByte(); err != nil {
		return nil, err
	}
	disk = append(disk, PadOrTrim(sector, globals.BytesPerSector)...)

	// append an empty unused sector
	// on Acorn disks the sector links to track 0, sector so
	// this is 'fudged' here.
	unusedSector := createEmptySector(0x00)
	unusedSector[1] = 5
	disk = append(disk, unusedSector...)

	// add the directory sectors
	disk = append(disk, d.Directory.ToByte()...)

	// add the data sectors (track 1 onwards)
	for i, ds := range d.DataSectors {
		if sector, err = ds.ToByte(); err != nil {
			return nil, err
		}
		if len(sector) != globals.BytesPerSector {
			return nil, fmt.Errorf("data sector %d size is incorrect", i)
		}
		disk = append(disk, sector...)
	}
	// track 0 plus the data sectors must make up one of the supported disk sizes
	if _, err = tracksForDiskSize(len(disk)); err != nil {
		return nil, fmt.Errorf("disk image built incorrectly: %w", err)
	}
	return disk, nil
}

func NewDisk(tracks int, volumeLabel string, volumeNumber int) (*Disk, error) {

	var (
		disk *Disk
		ds   *DataSector
		err  error
	)

	if !slices.Contains(supportedTrackCounts, tracks) {
		err = errors.New("only 35, 40 or 80 track disks are supported")
		return nil, err
	}

	disk = &Disk{}
	disk.BootSector1 = bootSector1
	disk.BootSector2 = createEmptySector(0x00) // boot sector 2 is blank in the Acorn system
	disk.BootSector2[1] = 0x03                 // add linked sector to boot sector 2

	if disk.SystemInfo, err = NewSIR(volumeLabel, volumeNumber, tracks); err != nil {
		return nil, errors.New("unable to initialise new system information record")
	}

	if disk.Directory, err = NewDirectory(); err != nil {
		return nil, errors.New("unable to initialise new directory")
	}

	// track 0 is the system track, so add track 1 onwards...
	sectorCount := (tracks - 1) * globals.SectorsPerTrack

	for s := 0; s < sectorCount; s++ {

		if s == sectorCount-1 {
			// last sector so no link address i.e. 0
			if ds, err = NewDataSector(0); err != nil {
				return nil, err
			}
		} else {
			// data starts at sectorId 10 but we are needing the link address to the
			//next sector so add 11
			if ds, err = NewDataSector(s + 11); err != nil {
				return nil, err
			}
		}
		disk.DataSectors = append(disk.DataSectors, ds)
	}

	return disk, nil
}
