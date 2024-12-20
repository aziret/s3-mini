package model

type FileChunk struct {
	UUID        string
	ChunkSize   int64
	ChunkNumber int64
	FilePath    string
}

type FileChunkUpload struct {
	UUID string
	Data []byte
}
