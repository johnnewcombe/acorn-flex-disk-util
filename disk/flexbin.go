// Package flexbin converts raw binary images into FLEX binary (.CMD/.BIN) format.
//
// A FLEX binary file is a stream of records:
//
//	Data record:     $02 addrHi addrLo count data...   (count = 1..255 bytes)
//	Transfer record: $16 addrHi addrLo                 (execution address, optional)
//
// The FLEX loader ignores $00 bytes between records, so the stream may be
// padded with zeros to fill the last disk sector.
package disk

import (
	"errors"
	"fmt"
)

// BinToFlex turns a raw binary image into a FLEX binary file.
// LoadAddress is where the first byte of the raw image is loaded.
// TransferAddress is the execution address used when HasTransfer is true.
// RecordSize is the number of data bytes per record (1..255). Zero means MaxRecordData.
func BinToFlex(raw []byte, loadAddress uint16, transferAddress uint16, hasTransfer bool, recordSize int) ([]byte, error) {
	const (
		recordData     = 0x02
		recordTransfer = 0x16

		// MaxRecordData is the most data bytes a single record can hold (the count is one byte).
		MaxRecordData = 255
	)

	//func Convert(raw []byte, opts Options) ([]byte, error) {
	size := recordSize
	if size == 0 {
		size = MaxRecordData
	}
	if size < 1 || size > MaxRecordData {
		return nil, fmt.Errorf("record size %d out of range 1..%d", size, MaxRecordData)
	}
	if int(loadAddress)+len(raw) > 0x10000 {
		return nil, fmt.Errorf("image of %d bytes at $%04X runs past $FFFF", len(raw), loadAddress)
	}
	if len(raw) == 0 && !hasTransfer {
		return nil, errors.New("empty image and no transfer address: nothing to write")
	}

	records := (len(raw) + size - 1) / size
	out := make([]byte, 0, len(raw)+records*4+3)

	addr := int(loadAddress)
	for off := 0; off < len(raw); off += size {
		end := off + size
		if end > len(raw) {
			end = len(raw)
		}
		chunk := raw[off:end]
		out = append(out, recordData, byte(addr>>8), byte(addr), byte(len(chunk)))
		out = append(out, chunk...)
		addr += len(chunk)
	}

	if hasTransfer {
		out = append(out, recordTransfer, byte(transferAddress>>8), byte(transferAddress))
	}

	return out, nil
}
