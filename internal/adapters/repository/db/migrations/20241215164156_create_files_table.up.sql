CREATE TABLE IF NOT EXISTS files(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  file_path TEXT,
  upload_id TEXT NOT NULL UNIQUE,
  offset INTEGER,
  filetype TEXT
);

CREATE INDEX IF NOT EXISTS idx_file_name ON files(name);
CREATE UNIQUE INDEX IF NOT EXISTS idx_file_upload_id ON files(upload_id);