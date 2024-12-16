package model

type File struct {
	ID   int64
	Info FileInfo
}

type FileInfo struct {
	ID             string
	Size           int64
	SizeIsDeferred bool
	Offset         int64
	MetaData       map[string]string
	IsPartial      bool
	IsFinal        bool
	PartialUploads []string
	Storage        map[string]string
}
