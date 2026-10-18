package disk

import (
	"bytes"
	"testing"
)

func TestParseRoundTrip(t *testing.T) {
	raw := make([]byte, 600)
	for i := range raw {
		raw[i] = byte(i)
	}
	bin, err := Convert(raw, 0x0100, 0x0123, true, 0)
	if err != nil {
		t.Fatal(err)
	}
	segs, tr, err := ParseFlexBinary(bin)
	if err != nil {
		t.Fatal(err)
	}
	if tr == nil || *tr != 0x0123 {
		t.Fatalf("transfer %v, want $0123", tr)
	}
	var joined []byte
	next := uint16(0x0100)
	for _, s := range segs {
		if s.Address != next {
			t.Fatalf("segment at $%04X, want $%04X", s.Address, next)
		}
		joined = append(joined, s.Data...)
		next += uint16(len(s.Data))
	}
	if !bytes.Equal(joined, raw) {
		t.Fatal("round-trip data mismatch")
	}
}

func TestParseSkipsPadding(t *testing.T) {
	b := []byte{0x00, 0x00, 0x02, 0x10, 0x00, 0x01, 0xAA, 0x00, 0x16, 0x10, 0x00, 0x00, 0x00}
	segs, tr, err := ParseFlexBinary(b)
	if err != nil {
		t.Fatal(err)
	}
	if len(segs) != 1 || segs[0].Address != 0x1000 || tr == nil || *tr != 0x1000 {
		t.Fatalf("got %+v transfer %v", segs, tr)
	}
}

func TestParseNoTransfer(t *testing.T) {
	_, tr, err := ParseFlexBinary([]byte{0x02, 0x20, 0x00, 0x01, 0x55})
	if err != nil {
		t.Fatal(err)
	}
	if tr != nil {
		t.Fatalf("transfer %v, want nil", *tr)
	}
}

func TestIsFlexBinaryRejectsRaw(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
	}{
		{"6809 code", []byte{0x86, 0x41, 0xBD, 0xCD, 0x18, 0x39}},
		{"all zeros", make([]byte, 252)},
		{"empty", nil},
		{"truncated header", []byte{0x02, 0x01, 0x00}},
		{"truncated data", []byte{0x02, 0x01, 0x00, 0x05, 0x01, 0x02}},
		{"zero count", []byte{0x02, 0x01, 0x00, 0x00}},
		{"runs past $FFFF", []byte{0x02, 0xFF, 0xFF, 0x02, 0x01, 0x02}},
		{"truncated transfer", []byte{0x16, 0x01}},
		{"stray byte after record", []byte{0x02, 0x01, 0x00, 0x01, 0xAA, 0x39}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if IsFlexBinary(tc.b) {
				t.Error("accepted as FLEX binary")
			}
		})
	}
}
