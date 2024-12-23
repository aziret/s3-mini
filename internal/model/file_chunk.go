package model

type FileChunk struct {
	UUID        string
	ChunkSize   int64
	ChunkNumber int64
	FilePath    string
	ServerID    string
}

type FileChunkUpload struct {
	UUID string
	Data []byte
}

type FileChunkDownload struct {
	UUID        string
	Data        []byte
	ChunkSize   int64
	ChunkNumber int64
}
