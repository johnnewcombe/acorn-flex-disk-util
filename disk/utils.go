package disk

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/johnnewcombe/System3-disk-util/globals"
)

// IntToBytes Returns the integer value as a byte slice of the specified length
func IntToBytes(val int, length int, bigEndian bool) []byte {

	const maxLen = 8
	if length > maxLen {
		length = maxLen
	}

	// convert to bytes from int
	result := make([]byte, maxLen)

	if bigEndian {
		binary.BigEndian.PutUint64(result, uint64(val))
		return result[maxLen-length:]
	} else {
		binary.LittleEndian.PutUint64(result, uint64(val))
		return result[:length]
	}

}

func BytesToInt(val []byte, bigEndian bool) int {

	if bigEndian {
		return int(binary.BigEndian.Uint16(val))
	} else {
		return int(binary.LittleEndian.Uint16(val))
	}
}

func HexToInt(value string) int {
	n := new(big.Int)
	n.SetString(value, 16)
	return int(n.Int64())
}

func IntToHex(value int) string {
	result := fmt.Sprintf("%04X", value)
	return result
}

func FormatVolumeLabel(volumeName string) (string, error) {

	vl := len(volumeName)

	if vl > 12 {
		volumeName = volumeName[:12]
	}

	if vl > 0 {
		// restrict filenames to A-F, a-z and 0-9 as its hex
		r, _ := regexp.Compile("^[A-Z0-9]+$")
		if !r.MatchString(volumeName) {
			return volumeName, errors.New("contains invalid characters")
		}
	}
	return volumeName, nil
	//return fmt.Sprintf("%-12s", strings.ToUpper(volumeName)), nil
}

// PadOrTrim returns (size) bytes from input (slice)
// Short input slice gets zeros suffixed, a long slice gets right trimmed.
func PadOrTrim(slice []byte, size int) []byte {

	l := len(slice)

	//correct size so return
	if l == size {
		return slice
	}

	// slice is too big so trim
	if l > size {
		return slice[:size]
	}

	// return new correctly sized slce and copy bba to it
	tmp := make([]byte, size)
	copy(tmp[:size], slice)
	return tmp
}

func CleanString(value string) string {

	clean := strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, value)
	return strings.TrimSpace(clean)

}

func createEmptySector(fill byte) []byte {
	var sector []byte
	var i int

	for i = 0; i < globals.BytesPerSector; i++ {
		sector = append(sector, fill)
	}
	return sector
}

func createEmptySectors(sectors int, fill byte) []byte {

	var sectorData []byte
	var i int

	for i = 0; i < sectors; i++ {

		sectorData = append(sectorData, createEmptySector(fill)...)
	}
	return sectorData
}

// getSectorAddress returns the track and sector representing the address of the specified
// zero based, sector id. See getSectorId().
func getSectorAddress(sectorId int) (int, int) {

	var (
		track int
		sect  int
	)

	// special case if sectorId is 0 then track and sector are both 0
	if sectorId == 0 {
		return 0, 0
	}

	// determine track
	track = sectorId / globals.SectorsPerTrack // data starts at the second track i.e. track 1
	sect = sectorId % globals.SectorsPerTrack

	// set to nect sector
	sect++

	return track, sect
}

// GetFlexDate returns day, month and year as integers for FLEX directory entries
// and the SIR. The year is stored as year-1950 (2026 = 76), a deliberate
// workaround for FLEX; this runs out in 2050, when the result reaches 100.
func GetFlexDate() (int, int, int) {

	year, month, day := time.Now().Date()
	return day, int(month), year - 1950
}

// getSectorId Returns the sector ID. This is the id of the sector on
// disk, where the first sector of the disk is zero.
func getSectorId(track int, sector int) int {

	// special case if sector and track are zeror, id is zero
	if track == 0 && sector == 0 {
		return 0
	}

	return track*globals.SectorsPerTrack + sector - 1
}
