CREATE TABLE IF NOT EXISTS servers(
    id UUID PRIMARY KEY
);

ALTER TABLE file_chunks
ADD COLUMN server_id UUID;

ALTER TABLE file_chunks
ADD CONSTRAINT fk_server
FOREIGN KEY (server_id)
REFERENCES servers (id);

CREATE INDEX IF NOT EXISTS idx_file_chunk_server_id ON file_chunks(server_id);