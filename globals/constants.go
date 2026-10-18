package globals

const (
	Version              = "0.10"
	BytesPerSector       = 256
	SectorsPerTrack      = 10
	BytesPerDirEntry     = 24
	FilenameLen          = 7
	DirectorySectors     = 6
	DirectoryStartSector = 5
	DirEntriesPerSector  = 10
	RandomFlag           = 2
	DeleteChar           = "\xff"
	BootSectorLink       = "FLEX9.SYS"
)
