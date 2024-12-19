DROP INDEX IF EXISTS idx_file_chunk_server_id;

ALTER TABLE file_chunks
DROP CONSTRAINT fk_server;

ALTER TABLE file_chunks
DROP COLUMN server_id;

DROP TABLE IF EXISTS servers;