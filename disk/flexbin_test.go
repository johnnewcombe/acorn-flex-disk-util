package disk

import (
	"bytes"
	"testing"
)

func TestConvertSmallImageWithTransfer(t *testing.T) {
	got, err := BinToFlex([]byte{0x86, 0x41, 0x39}, 0x0100, 0x0100, true, 0)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x02, 0x01, 0x00, 0x03, 0x86, 0x41, 0x39, 0x16, 0x01, 0x00}
	if !bytes.Equal(got, want) {
		t.Fatalf("got  % X\nwant % X", got, want)
	}
}

func TestConvertNoTransfer(t *testing.T) {
	got, err := BinToFlex([]byte{0xAA, 0xBB}, 0x2000, 0, false, 0)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x02, 0x20, 0x00, 0x02, 0xAA, 0xBB}
	if !bytes.Equal(got, want) {
		t.Fatalf("got  % X\nwant % X", got, want)
	}
}

func TestConvertTransferOnly(t *testing.T) {
	got, err := BinToFlex(nil, 0, 0xC100, true, 0)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x16, 0xC1, 0x00}
	if !bytes.Equal(got, want) {
		t.Fatalf("got  % X\nwant % X", got, want)
	}
}

func TestConvertSplitsIntoRecords(t *testing.T) {
	raw := make([]byte, 600)
	got, err := BinToFlex(raw, 0xC100, 0, false, 0)
	if err != nil {
		t.Fatal(err)
	}
	// 255 + 255 + 90 bytes: three records, each with a 4-byte header
	if len(got) != 600+3*4 {
		t.Fatalf("length %d, want %d", len(got), 600+3*4)
	}
	// second record starts at $C100 + 255 = $C1FF
	h := got[4+255:]
	if h[0] != 0x02 || h[1] != 0xC1 || h[2] != 0xFF || h[3] != 0xFF {
		t.Fatalf("second header % X", h[:4])
	}
	// third record: $C2FE, 90 bytes
	h = got[2*(4+255):]
	if h[0] != 0x02 || h[1] != 0xC2 || h[2] != 0xFE || h[3] != 90 {
		t.Fatalf("third header % X", h[:4])
	}
}

func TestConvertCustomRecordSize(t *testing.T) {
	got, err := BinToFlex([]byte{1, 2, 3, 4, 5}, 0x0100, 0, false, 2)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{
		0x02, 0x01, 0x00, 0x02, 1, 2,
		0x02, 0x01, 0x02, 0x02, 3, 4,
		0x02, 0x01, 0x04, 0x01, 5,
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("got  % X\nwant % X", got, want)
	}
}

func TestConvertFillsToTopOfMemory(t *testing.T) {
	// 16 bytes at $FFF0 ends exactly at $FFFF, which is allowed
	if _, err := BinToFlex(make([]byte, 0x10), 0xFFF0, 0, false, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConvertErrors(t *testing.T) {
	tests := []struct {
		name        string
		raw         []byte
		load        uint16
		hasTransfer bool
		recordSize  int
	}{
		{"runs past $FFFF", make([]byte, 0x11), 0xFFF0, false, 0},
		{"record size too large", []byte{1}, 0, false, 256},
		{"record size negative", []byte{1}, 0, false, -1},
		{"nothing to write", nil, 0, false, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := BinToFlex(tc.raw, tc.load, 0, tc.hasTransfer, tc.recordSize); err == nil {
				t.Error("expected an error")
			}
		})
	}
}
