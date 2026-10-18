package disk

import (
	"errors"
	"fmt"
)

// Segment is one block of bytes a FLEX binary loads at Address.
type Segment struct {
	Address uint16
	Data    []byte
}

// ParseFlexBinary checks that b is a valid FLEX binary record stream and
// returns its data segments and transfer address (nil if there is none).
//
// Rules: $00 bytes between records are padding and are skipped; $02 starts a
// data record (address, count 1..255, data); $16 is a transfer record
// (address). Any other byte, a truncated record, a zero count or a segment
// running past $FFFF makes the stream invalid.
func ParseFlexBinary(b []byte) ([]Segment, *uint16, error) {
	var segs []Segment
	var transfer *uint16

	for i := 0; i < len(b); {
		switch b[i] {
		case 0x00:
			i++

		case 0x02:
			if i+4 > len(b) {
				return nil, nil, fmt.Errorf("offset %d: truncated data record header", i)
			}
			addr := uint16(b[i+1])<<8 | uint16(b[i+2])
			n := int(b[i+3])
			if n == 0 {
				return nil, nil, fmt.Errorf("offset %d: data record with zero length", i)
			}
			if i+4+n > len(b) {
				return nil, nil, fmt.Errorf("offset %d: data record runs past end of file", i)
			}
			if int(addr)+n > 0x10000 {
				return nil, nil, fmt.Errorf("offset %d: data record at $%04X runs past $FFFF", i, addr)
			}
			segs = append(segs, Segment{Address: addr, Data: b[i+4 : i+4+n]})
			i += 4 + n

		case 0x16:
			if i+3 > len(b) {
				return nil, nil, fmt.Errorf("offset %d: truncated transfer record", i)
			}
			t := uint16(b[i+1])<<8 | uint16(b[i+2])
			transfer = &t
			i += 3

		default:
			return nil, nil, fmt.Errorf("offset %d: unexpected byte $%02X", i, b[i])
		}
	}

	if len(segs) == 0 && transfer == nil {
		return nil, nil, errors.New("no records found")
	}
	return segs, transfer, nil
}

// IsFlexBinary reports whether b parses completely as a FLEX binary record stream.
func IsFlexBinary(b []byte) bool {
	_, _, err := ParseFlexBinary(b)
	return err == nil
}
