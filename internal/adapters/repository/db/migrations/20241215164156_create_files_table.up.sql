CREATE TABLE IF NOT EXISTS files(
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  file_path TEXT,
  upload_id TEXT NOT NULL UNIQUE,
    size INTEGER NOT NULL,
  chunk_size INTEGER NOT NULL,
  "offset" BIGINT,
  filetype TEXT,
    ready_to_download BOOLEAN NOT NULL DEFAULT FALSE,
    download_completed BOOLEAN NOT NULL DEFAULT FALSE,
    file_chunks_created BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_file_name ON files(name);
CREATE UNIQUE INDEX IF NOT EXISTS idx_file_upload_id ON files(upload_id);
CREATE INDEX IF NOT EXISTS idx_file_ready_to_download ON files(ready_to_download);
CREATE INDEX IF NOT EXISTS idx_file_download_completed ON files(download_completed);
CREATE INDEX IF NOT EXISTS idx_file_file_chunks_created ON files(file_chunks_created);