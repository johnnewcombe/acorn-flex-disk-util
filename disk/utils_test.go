package disk

import "testing"

func TestGetSectorAddress(t *testing.T) {
	tests := []struct {
		name      string
		sectorId  int
		wantTrack int
		wantSect  int
	}{
		{"First sector of first track", 0, 0, 0}, // special case
		{"Last sector of first track", 9, 0, 10},
		{"First sector of second track", 10, 1, 1},
		{"Last sector of second track", 19, 1, 10},
		{"A middle sector", 24, 2, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTrack, gotSect := getSectorAddress(tt.sectorId)
			if gotTrack != tt.wantTrack || gotSect != tt.wantSect {
				t.Errorf("getSectorAddress(%d) gotTrack = %v, wantTrack %v, gotSect = %v, wantSect %v", tt.sectorId, gotTrack, tt.wantTrack, gotSect, tt.wantSect)
			}
		})
	}
}

func TestGetSectorId(t *testing.T) {
	tests := []struct {
		name     string
		track    int
		sector   int
		wantSect int
	}{
		{"First sector of first track", 0, 1, 0},
		{"Last sector of first track", 0, 10, 9},
		{"First sector of second track", 1, 1, 10},
		{"Last sector of second track", 1, 10, 19},
		{"A middle sector", 23, 4, 233},
		{"A middle sector", 23, 5, 234},
		{"A middle sector", 0, 0, 0}, // speial case
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSect := getSectorId(tt.track, tt.sector)
			if gotSect != tt.wantSect {
				t.Errorf("getSectorId(%d, %d) = %v, want %v", tt.track, tt.sector, gotSect, tt.wantSect)
			}
		})
	}
}
