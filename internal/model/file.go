package model

type File struct {
	ID              int64
	Size            int64
	Offset          int64
	UploadID        string
	Name            string
	FilePath        string
	FileType        string
	ReadyToDownload bool
}

type FileInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
