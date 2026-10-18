package disk

import (
	"fmt"
	"strings"

	"github.com/johnnewcombe/System3-disk-util/globals"
)

type SIR struct {
	VolumeLabel       string // 11 chars
	VolumeNumber      int    // 2 bytes
	FirstFreeTrack    int    // 1 byte
	FirstFreeSector   int    // 1 byte
	LastFreeTrack     int    // 1 byte
	LastFreeSector    int    // 1 byte
	NumberFreeSectors int    // 2 bytes
	Month             int    // 1 byte
	Day               int    // 1 byte
	Year              int    // 1 byte
	EndTrack          int    // 1 byte
	EndSector         int    // 1 byte
}

func (s *SIR) Load(data []byte) error {

	if len(data) < 40 {
		return fmt.Errorf("system information record is %d bytes, at least 40 needed", len(data))
	}

	s.VolumeLabel = CleanString(string(data[16:27]))
	s.VolumeNumber = BytesToInt(data[27:29], true)
	s.FirstFreeTrack = int(data[29])
	s.FirstFreeSector = int(data[30])
	s.LastFreeTrack = int(data[31])
	s.LastFreeSector = int(data[32])
	s.NumberFreeSectors = BytesToInt(data[33:35], true)
	s.Month = int(data[35])
	s.Day = int(data[36])
	s.Year = int(data[37])
	s.EndTrack = int(data[38])
	s.EndSector = int(data[39])

	return nil
}

func (s *SIR) ToByte() ([]byte, error) {

	var (
		sector  []byte
		padding []byte
	)
	padding = make([]byte, 16)
	sector = append(sector, padding...)

	volumeLabel := PadOrTrim([]byte(strings.ToUpper(s.VolumeLabel)), 11)

	sector = append(sector, volumeLabel...)
	sector = append(sector, IntToBytes(s.VolumeNumber, 2, true)...)
	sector = append(sector, IntToBytes(s.FirstFreeTrack, 1, false)...)
	sector = append(sector, IntToBytes(s.FirstFreeSector, 1, false)...)
	sector = append(sector, IntToBytes(s.LastFreeTrack, 1, false)...)
	sector = append(sector, IntToBytes(s.LastFreeSector, 1, false)...)
	sector = append(sector, IntToBytes(s.NumberFreeSectors, 2, true)...)
	sector = append(sector, IntToBytes(s.Month, 1, false)...)
	sector = append(sector, IntToBytes(s.Day, 1, false)...)
	sector = append(sector, IntToBytes(s.Year, 1, false)...)
	sector = append(sector, IntToBytes(s.EndTrack, 1, false)...)
	sector = append(sector, IntToBytes(s.EndSector, 1, false)...)

	return sector, nil
}

func NewSIR(volumeLabel string, volumeNumber int, tracks int) (*SIR, error) {

	var (
		err error
	)

	if volumeLabel, err = FormatVolumeLabel(volumeLabel); err != nil {
		return nil, fmt.Errorf("volume name contains invalid characters: %w", err)
	}

	day, month, year := GetFlexDate()

	sir := &SIR{
		VolumeLabel:       volumeLabel,
		VolumeNumber:      volumeNumber,
		FirstFreeTrack:    1,
		FirstFreeSector:   1,
		LastFreeTrack:     tracks - 1,              // number of tracks -1 (tracks are zero based)
		LastFreeSector:    globals.SectorsPerTrack, // sectors are 1 based
		NumberFreeSectors: (tracks - 1) * globals.SectorsPerTrack,
		Month:             int(month),
		Day:               day,
		Year:              year,       // GetFlexDate applies the FLEX year workaround (2026 = 76)
		EndTrack:          tracks - 1, // tracks are 0 based
		EndSector:         globals.SectorsPerTrack,
	}

	return sir, nil
}
