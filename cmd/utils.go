package cmd

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

// getFlexFilename Validates and returns the correctly formatted flex file name.
func getFlexFilename(localFilename string) (string, error) {

	// get filename from path
	filename := strings.ToUpper(filepath.Base(localFilename))

	// Single regex to match 8.3 filename:
	// Filename part: starts with A-Z, 1-8 chars total, allowed A-Z, 0-9, _, -
	// Extension part: starts with A-Z, exactly 3 chars, allowed A-Z, 0-9, _, -
	re := regexp.MustCompile(`^([A-Z][A-Z0-9_-]{0,7})\.([A-Z][A-Z0-9_-]{2})$`)
	matches := re.FindStringSubmatch(filename)

	if matches == nil {
		return filename, errors.New("invalid 8.3 filename (must be NAME.EXT, start with a letter, and use only A-Z, 0-9, _ or -)")
	}
	return filename, nil
}

// getVolumeLabel validates and returns the correctly formatted volume label.
func getVolumeLabel(volumeLabel string) (string, error) {

	// get filename from path
	label := strings.ToUpper(filepath.Base(volumeLabel))

	// Single regex to match volume label:
	// Label part: starts with A-Z, 1-8 chars total, allowed A-Z, 0-9, _, -
	re := regexp.MustCompile(`^([A-Z][A-Z0-9_-]{0,10})$`)
	matches := re.FindStringSubmatch(label)

	if matches == nil {
		return label, errors.New("invalid volume label (must start with a letter, 1-8 characters, and use only A-Z, 0-9, _ or -)")
	}
	return label, nil
}
