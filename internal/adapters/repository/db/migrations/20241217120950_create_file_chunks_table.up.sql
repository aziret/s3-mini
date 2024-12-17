CREATE TABLE IF NOT EXISTS file_chunks(
    uuid TEXT PRIMARY KEY,
    file_id INTEGER NOT NULL,
    chunk_size INTEGER NOT NULL,
    chunk_number INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_file_chunk_file_id ON file_chunks(file_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_file_chunk_number_file_id ON file_chunks(file_id, chunk_number);