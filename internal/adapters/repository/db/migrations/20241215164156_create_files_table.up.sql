CREATE TABLE IF NOT EXISTS files(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  file_path TEXT,
  checksum TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_file_name ON files(name);
CREATE UNIQUE INDEX IF NOT EXISTS idx_file_checksum ON files(checksum);