CREATE TABLE IF NOT EXISTS file_chunks(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id INTEGER NOT NULL,
    chunk_size INTEGER NOT NULL,
    chunk_number INTEGER NOT NULL,
    download_completed BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_file_chunk_file_id ON file_chunks(file_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_file_chunk_number_file_id ON file_chunks(file_id, chunk_number);
CREATE INDEX IF NOT EXISTS idx_file_chunk_download_completed ON file_chunks(download_completed);