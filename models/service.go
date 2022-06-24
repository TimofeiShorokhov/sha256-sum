package models

type HashDataUtils struct {
	FileName  string
	Checksum  string
	FilePath  string
	Algorithm string
}

type ChangedHashes struct {
	FileName    string
	OldChecksum string
	NewChecksum string
	FilePath    string
	Algorithm   string
}

type PodInfo struct {
	PodName       string
	CreationTime  string
	ImageName     string
	ContainerName string
}
