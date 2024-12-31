package file

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/aziret/s3-mini/internal/model"
)

func (s *Service) DownloadFile(ctx context.Context, id int64) (string, error) {
	const op = "service.file.DownloadFile"
	log := s.logger.With(
		slog.String("op", op),
		slog.Int64("ID", id),
	)

	serverIDs, err := s.fileRepo.GetFileChunksServerIDs(ctx, id)

	if err != nil {
		log.Error("failed to get server IDs", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	workersPool := os.Getenv("WORKERS_POOL")
	workersNumber := 3
	if workersPool != "" {
		number, err := strconv.Atoi(workersPool)
		if err != nil {
			workersNumber = 3
		}
		workersNumber = number
	}

	workersNumber *= 2

	workersChannel := make(chan struct{}, workersNumber)

	downloadedFileChunks := make(chan model.FileChunkDownload)

	streamWG := sync.WaitGroup{}
	downloadWG := sync.WaitGroup{}
	downloadWG.Add(1)
	go func() {
		defer downloadWG.Done()
		for _, serverID := range *serverIDs {
			workersChannel <- struct{}{}
			workersChannel <- struct{}{}

			fileChunksChan := make(chan model.FileChunk)

			streamWG.Add(2)
			go func(ctx context.Context, serverID string, fileChunksChan chan model.FileChunk) {
				defer streamWG.Done()
				defer func() { <-workersChannel }()

				s.downloadFileChunks(ctx, serverID, fileChunksChan, downloadedFileChunks)
			}(ctx, serverID, fileChunksChan)
			go func(ctx context.Context, fileID int64, serverID string, fileChunksChan chan model.FileChunk) {
				defer streamWG.Done()
				defer func() { <-workersChannel }()

				s.sendFileChunksToDownloadToChannel(ctx, fileID, serverID, fileChunksChan)
			}(ctx, id, serverID, fileChunksChan)
		}
	}()

	filepath := fmt.Sprintf("./downloads/%d", id)

	file, err := os.Create(filepath)
	if err != nil {
		log.Error("failed to create file", slog.String("filepath", filepath), sl.Err(err))
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error("error closing file", slog.Int64("file ID", id), sl.Err(err))
		}
	}(file)

	composeWG := sync.WaitGroup{}
	composeWG.Add(1)
	go func() {
		defer composeWG.Done()
		for downloadedFileChunk := range downloadedFileChunks {
			byteBegin := downloadedFileChunk.ChunkSize * downloadedFileChunk.ChunkNumber
			_, err = file.Seek(byteBegin, 0)
			if err != nil {
				log.Error(
					"error seeking to begin of file",
					slog.Int64("Chunk Number", downloadedFileChunk.ChunkNumber),
					slog.Int64("Chunk Size", downloadedFileChunk.ChunkSize),
					sl.Err(err),
				)
				return
			}
			_, err = file.Write(downloadedFileChunk.Data)
			if err != nil {
				log.Error(
					"error writing to file",
					slog.String("UUID", downloadedFileChunk.UUID),
					slog.Int64("Chunk Number", downloadedFileChunk.ChunkNumber),
					slog.Int64("Chunk Size", downloadedFileChunk.ChunkSize),
					sl.Err(err),
				)
				return
			}
		}
	}()

	streamWG.Wait()
	downloadWG.Wait()
	close(downloadedFileChunks)
	composeWG.Wait()

	return filepath, nil
}

func (s *Service) downloadFileChunks(ctx context.Context, serverID string, fileChunksChan <-chan model.FileChunk, downloadedFileChunks chan<- model.FileChunkDownload) {
	s.serversMap.mu.RLock()
	defer s.serversMap.mu.RUnlock()
	grpcClient := s.serversMap.servers[serverID]

	grpcClient.DownloadFileChunks(ctx, fileChunksChan, downloadedFileChunks)
}

func (s *Service) sendFileChunksToDownloadToChannel(ctx context.Context, fileID int64, serverID string, fileChunksChan chan<- model.FileChunk) {
	defer close(fileChunksChan)
	const op = "service.file.sendFileChunksToDownloadToChannel"

	log := s.logger.With(
		slog.String("op", op),
		slog.Int64("fileID", fileID),
	)

	fileChunks, err := s.fileRepo.GetFileChunksByFileIDAndServerID(ctx, fileID, serverID)
	if err != nil {
		log.Error("failed to get file chunks", sl.Err(err))
		return
	}

	for _, fileChunk := range *fileChunks {
		fileChunksChan <- fileChunk
	}
}
