package cmd

import (
	"testing"
)

func TestGetFlexFilename(t *testing.T) {
	tests := []struct {
		name          string
		localFilename string
		want          string
		wantErr       bool
	}{
		{"valid 8.3", "./dir/FILE.EXT", "FILE.EXT", false},
		{"valid 8.3", "FILE.EXT", "FILE.EXT", false},
		{"valid max name", "FILENAME.EXT", "FILENAME.EXT", false},
		{"valid with underscore/hyphen", "A_1-B.E_1", "A_1-B.E_1", false},
		{"missing extension", "FILENAME", "FILENAME", true},
		{"too long extension", "FILE.EXTN", "FILE.EXTN", true},
		{"too short extension", "FILE.EX", "FILE.EX", true},
		{"filename starts with digit", "1FILE.EXT", "1FILE.EXT", true},
		{"extension starts with digit", "FILE.1XT", "FILE.1XT", true},
		{"too long filename", "VERYLONGNAME.EXT", "VERYLONGNAME.EXT", true},
		{"invalid char in name", "FILE@.EXT", "FILE@.EXT", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFlexFilename(tt.localFilename)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFlexFilename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("getFlexFilename() got = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetVolumeLabel(t *testing.T) {
	tests := []struct {
		name        string
		volumeLabel string
		want        string
		wantErr     bool
	}{
		{"valid volume label", "FLEXDISK", "FLEXDISK", false},
		{"valid 1 char label", "A", "A", false},
		{"valid max label", "VOLNAME1234", "VOLNAME1234", false},
		{"label with hyphen", "FLEX-DSK", "FLEX-DSK", false},
		{"label starts with digit", "1VOLNAME", "1VOLNAME", true},
		{"too long label", "VERYLONGNAME", "VERYLONGNAME", true},
		{"invalid char", "LABEL@", "LABEL@", true},
		{"with path", "/path/to/LABEL", "LABEL", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getVolumeLabel(tt.volumeLabel)
			if (err != nil) != tt.wantErr {
				t.Errorf("getVolumeLabel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("getVolumeLabel() got = %q, want %q", got, tt.want)
			}
		})
	}
}
