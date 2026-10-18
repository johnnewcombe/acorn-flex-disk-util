package globals

const (
	Version              = "1.00"
	BytesPerSector       = 256
	SectorsPerTrack      = 10
	BytesPerDirEntry     = 24
	DirectorySectors     = 6
	DirectoryStartSector = 5
	DirEntriesPerSector  = 1
	DeleteChar           = "\xff"
	BootSectorLink       = "FLEX9.SYS"
)
