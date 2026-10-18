package disk

import (
	"fmt"
	"strings"

	"github.com/johnnewcombe/System3-disk-util/globals"
)

// ValidationReport is the result of checking a FLEX disk image.
type ValidationReport struct {
	Tracks      int
	Files       int
	FreeSectors int
	Errors      []string // structural damage: the disk is not safe to write to
	Warnings    []string // suspicious but not necessarily damaged
}

// OK reports whether no errors were found (warnings are allowed).
func (r *ValidationReport) OK() bool { return len(r.Errors) == 0 }

func (r *ValidationReport) errorf(format string, a ...any) {
	r.Errors = append(r.Errors, fmt.Sprintf(format, a...))
}

func (r *ValidationReport) warnf(format string, a ...any) {
	r.Warnings = append(r.Warnings, fmt.Sprintf(format, a...))
}

// ValidateImage checks the structure of a raw FLEX disk image: the system
// information record, the directory, every file's sector chain, the free chain,
// and that every sector belongs to exactly one of those. The contents of .CMD
// and .BIN files are also checked to be valid FLEX binary.
//
// It works directly on the image bytes rather than through Disk, so that it can
// report on images too damaged for Disk.Load to accept.
func ValidateImage(image []byte) *ValidationReport {
	const bps = globals.BytesPerSector
	spt := globals.SectorsPerTrack
	r := &ValidationReport{}

	tracks, err := tracksForDiskSize(len(image))
	if err != nil {
		r.errorf("%v", err)
		return r
	}
	r.Tracks = tracks

	type ts struct{ t, s int }
	sector := func(a ts) []byte {
		off := (a.t*spt + a.s - 1) * bps
		return image[off : off+bps]
	}
	valid := func(a ts) bool { return a.t >= 0 && a.t < tracks && a.s >= 1 && a.s <= spt }
	name := func(a ts) string { return fmt.Sprintf("T%02d-S%02d", a.t, a.s) }

	// owner records which chain each sector belongs to, to find overlaps and lost sectors
	owner := map[ts]string{}
	claim := func(a ts, who string) {
		if prev, taken := owner[a]; taken {
			r.errorf("sector %s belongs to both %s and %s", name(a), prev, who)
			return
		}
		owner[a] = who
	}

	// walk follows a linked chain from start, claiming each sector for who.
	// It returns the sectors in order; a bad link or a loop ends the walk with an error.
	walk := func(start ts, who string) []ts {
		var chain []ts
		seen := map[ts]bool{}
		for a := start; a != (ts{0, 0}); {
			if !valid(a) {
				r.errorf("%s: link to invalid sector %s after %d sectors", who, name(a), len(chain))
				break
			}
			if seen[a] {
				r.errorf("%s: chain loops back to %s after %d sectors", who, name(a), len(chain))
				break
			}
			seen[a] = true
			chain = append(chain, a)
			claim(a, who)
			b := sector(a)
			a = ts{int(b[0]), int(b[1])}
		}
		return chain
	}

	// ---- system area ----------------------------------------------------------------
	for s := 1; s <= 4; s++ {
		claim(ts{0, s}, "the system area")
	}

	// ---- system information record (track 0, sector 3) ------------------------------
	sir := sector(ts{0, 3})
	firstFree := ts{int(sir[29]), int(sir[30])}
	lastFree := ts{int(sir[31]), int(sir[32])}
	freeCount := int(sir[33])<<8 | int(sir[34])
	end := ts{int(sir[38]), int(sir[39])}
	r.FreeSectors = freeCount

	if end != (ts{tracks - 1, spt}) {
		r.errorf("SIR says the disk ends at %s, but a %d byte image ends at %s",
			name(end), len(image), name(ts{tracks - 1, spt}))
	}
	if freeCount == 0 {
		if firstFree != (ts{0, 0}) || lastFree != (ts{0, 0}) {
			r.errorf("SIR says no sectors are free but the free chain runs %s to %s", name(firstFree), name(lastFree))
		}
	} else if firstFree == (ts{0, 0}) || lastFree == (ts{0, 0}) {
		r.errorf("SIR says %d sectors are free but the free chain is %s to %s",
			freeCount, name(firstFree), name(lastFree))
	}

	// ---- directory (starts at track 0, sector 5) -------------------------------------
	type entry struct {
		name          string
		start, finish ts
		count         int
		random        bool
	}
	var entries []entry
	names := map[string]bool{}

	for _, ds := range walk(ts{0, 5}, "the directory") {
		b := sector(ds)
		for i := 0; i < 10; i++ {
			e := b[16+i*24 : 16+(i+1)*24]
			if e[0] == 0x00 || e[0]&0x80 != 0 { // never used, or deleted
				continue
			}
			n := strings.TrimRight(string(e[0:8]), "\x00") + "." + strings.TrimRight(string(e[8:11]), "\x00")
			if names[n] {
				r.warnf("%s appears more than once in the directory", n)
			}
			names[n] = true
			entries = append(entries, entry{
				name:   n,
				start:  ts{int(e[13]), int(e[14])},
				finish: ts{int(e[15]), int(e[16])},
				count:  int(e[17])<<8 | int(e[18]),
				random: e[19] != 0,
			})
		}
	}
	r.Files = len(entries)

	// ---- every file's chain, and the contents of binary files -------------------------
	for _, e := range entries {
		if e.count == 0 {
			r.errorf("%s: directory says it has no sectors", e.name)
			continue
		}
		chain := walk(e.start, e.name)
		if len(chain) != e.count {
			r.errorf("%s: chain has %d sectors but the directory says %d", e.name, len(chain), e.count)
		}
		if len(chain) > 0 && chain[len(chain)-1] != e.finish {
			r.errorf("%s: chain ends at %s but the directory says %s", e.name, name(chain[len(chain)-1]), name(e.finish))
		}
		if !e.random { // random files number their records differently (sector map first)
			for i, a := range chain {
				b := sector(a)
				if rec := int(b[2])<<8 | int(b[3]); rec != i+1 {
					r.errorf("%s: sector %s has record number %d, expected %d", e.name, name(a), rec, i+1)
					break
				}
			}
		}

		ext := e.name[strings.LastIndex(e.name, ".")+1:]
		if !e.random && (ext == "CMD" || ext == "BIN") {
			var data []byte
			for _, a := range chain {
				data = append(data, sector(a)[4:]...)
			}
			if _, _, err := ParseFlexBinary(data); err != nil {
				r.errorf("%s: contents are not valid FLEX binary: %v", e.name, err)
			}
		}
	}

	// ---- free chain ----------------------------------------------------------------
	if firstFree != (ts{0, 0}) {
		free := walk(firstFree, "the free chain")
		if len(free) != freeCount {
			r.errorf("free chain has %d sectors but the SIR says %d", len(free), freeCount)
		}
		if len(free) > 0 && free[len(free)-1] != lastFree {
			r.errorf("free chain ends at %s but the SIR says %s", name(free[len(free)-1]), name(lastFree))
		}
	}

	// ---- sectors that belong to nothing ---------------------------------------------
	var lost []string
	for t := 0; t < tracks; t++ {
		for s := 1; s <= spt; s++ {
			if _, ok := owner[ts{t, s}]; !ok {
				lost = append(lost, name(ts{t, s}))
			}
		}
	}
	if len(lost) > 0 {
		shown := lost
		if len(shown) > 8 {
			shown = append(shown[:8:8], "...")
		}
		r.warnf("%d sectors belong to no file and are not in the free chain: %s",
			len(lost), strings.Join(shown, " "))
	}

	return r
}
