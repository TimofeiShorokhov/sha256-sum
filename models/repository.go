package models

type HashData struct {
	Id        int
	FileName  string
	CheckSum  string
	FilePath  string
	Algorithm string
}

type PodData struct {
	PodName       string
	CreationTime  string
	ImageName     string
	ContainerName string
}
