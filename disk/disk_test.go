package disk

import (
	"testing"
)

const (
	TEST_ERROR_MESSAGE = "Test Description: \"%s.\""
)

func TestDisk_AddBinaryFile(t *testing.T) {
	tests := []struct {
		name         string
		initialFiles []struct {
			name string
			data []byte
		}
		filename    string
		fileData    []byte
		random      bool
		wantErr     bool
		expectedErr string
	}{
		{
			name:     "empty file data returns error",
			filename: "EMPTY.BIN",
			fileData: []byte{},
			wantErr:  true,
		},
		{
			name: "duplicate filename returns error",
			initialFiles: []struct {
				name string
				data []byte
			}{
				{name: "EXIST.BIN", data: []byte("initial data")},
			},
			filename: "EXIST.BIN",
			fileData: []byte("new data"),
			wantErr:  true,
		},
		{
			name:     "add small file (< 252 bytes)",
			filename: "SMALL.TXT",
			fileData: []byte("Hello Flex Disk!"),
			wantErr:  false,
		},
		{
			name:     "add exact 1 sector file (252 bytes)",
			filename: "EXACT1.BIN",
			fileData: make([]byte, 252),
			wantErr:  false,
		},
		{
			name:     "add 2 sector file (253 bytes)",
			filename: "TWOSECT.BIN",
			fileData: make([]byte, 253),
			wantErr:  false,
		},
		{
			name:     "add 4096 byte binary file (PRESTEL.BIN / 17 sectors)",
			filename: "PRESTEL.BIN",
			fileData: make([]byte, 4096),
			wantErr:  false,
		},
		{
			name: "add multiple files sequentially without corrupting free chain",
			initialFiles: []struct {
				name string
				data []byte
			}{
				{name: "FIRST.BIN", data: make([]byte, 256)},
				{name: "SECOND.BIN", data: make([]byte, 4096)},
			},
			filename: "THIRD.BIN",
			fileData: make([]byte, 500),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDisk(80, "TESTDISK", 1)
			if err != nil {
				t.Fatalf("failed to create disk: %v", err)
			}

			for _, f := range tt.initialFiles {
				if err := d.AddBinaryFile(f.name, f.data, false); err != nil {
					t.Fatalf("failed to setup initial file %s: %v", f.name, err)
				}
			}

			err = d.AddBinaryFile(tt.filename, tt.fileData, tt.random)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddBinaryFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDisk_Create(t *testing.T) {
	t.Error("Test not implemented!")
}

func TestDisk_ToByte(t *testing.T) {
	t.Error("Test not implemented!")
}

func TestDisk_GetNextFreeSector(t *testing.T) {
	//t.Error("Test not implemented!")
}

func TestDisk_HexToInt(t *testing.T) {

	type Test struct {
		description string
		input       string
		output      int
	}

	var tests = []Test{
		{"HexToInt Test 1", "A000", 40960},
		{"HexToInt Test 2", "1200", 4608},
		{"HexToInt Test 3", "AAAA", 43690},
		{"HexToInt Test 4", "0A", 10},
	}

	for _, test := range tests {

		if got := HexToInt(test.input); got != test.output {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func TestDisk_IntToHex(t *testing.T) {

	type Test struct {
		description string
		input       int
		output      string
	}

	var tests = []Test{
		{"HexToInt Test 1", 40960, "A000"},
		{"HexToInt Test 2", 4608, "1200"},
		{"HexToInt Test 3", 43690, "AAAA"},
		{"HexToInt Test 4", 10, "000A"},
	}

	for _, test := range tests {

		if got := IntToHex(test.input); got != test.output {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func TestDisk_FileExists(t *testing.T) {
	d, _ := NewDisk(80, "TEST", 1)

	// Add a file manually
	de := d.Directory.DirectorySectors[0].DirectoryEntries[0]
	de.Name = "TESTFILE"
	de.Ext = "TXT"

	// Add another file with some problematic name
	de2 := d.Directory.DirectorySectors[0].DirectoryEntries[1]
	de2.Name = "FILE    " // padded with spaces
	de2.Ext = "ABC"

	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{"exact match", "TESTFILE.TXT", true},
		{"lowercase match", "testfile.txt", true}, // FLEX is usually uppercase but we want case-insensitive check
		{"different extension", "TESTFILE.DAT", false},
		{"substring", "TESTFILE", false},
		{"with spaces", "TESTFILE.TXT ", false},
		{"padded name match", "FILE.ABC", true},
		{"lowercase match check", "TESTFILE.txt", true},
		{"all lowercase", "testfile.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := d.FileExists(tt.filename); got != tt.want {
				t.Errorf("Disk.FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
