package file

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/aziret/s3-mini/internal/model"
)

func (s *Service) UploadFileChunks(ctx context.Context) {
	const op = "service.file.UploadFileChunks"

	log := s.logger.With(
		slog.String("op", op),
	)

	s.serversMap.mu.RLock()
	defer s.serversMap.mu.RUnlock()
	if len(s.serversMap.servers) == 0 {
		log.Info("no servers found")
		return
	}

	fileChunksToSend, err := s.fileRepo.GetNotSentFileChunks(ctx)
	if err != nil {
		log.Error("failed to query file chunks", sl.Err(err))
		return
	}

	err = s.enqueueJobsToUploadFileChunks(ctx, fileChunksToSend)
	if err != nil {
		log.Error("failed to enqueue file chunks", sl.Err(err))
		return
	}
	return
}

func (s *Service) enqueueJobsToUploadFileChunks(ctx context.Context, fileChunks *[]model.FileChunk) error {
	const op = "service.file.enqueueJobsToUploadFileChunks"

	log := s.logger.With(
		slog.String("op", op),
	)

	s.serversMap.mu.RLock()
	defer s.serversMap.mu.RUnlock()
	if len(s.serversMap.servers) == 0 {
		log.Info("no servers found")
		return nil
	}

	serverIDs := make([]string, 0, len(s.serversMap.servers))
	for id := range s.serversMap.servers {
		serverIDs = append(serverIDs, id)
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	currentIdx := rand.Intn(len(serverIDs))

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

	wg := sync.WaitGroup{}

	workersChannel := make(chan struct{}, workersNumber)

	i := 0

	for i < len(serverIDs) {
		workersChannel <- struct{}{}
		workersChannel <- struct{}{}

		currentIdx += i
		currentIdx %= len(serverIDs)

		fileChunksChan := make(chan model.FileChunkUpload)
		serverID := serverIDs[currentIdx]

		wg.Add(2)
		go func(idx int, serverID string) {
			defer wg.Done()
			defer func() { <-workersChannel }()

			s.transferFileChunks(ctx, serverID, fileChunksChan)
		}(currentIdx, serverID)
		go func(startIdx int, serverID string) {
			defer wg.Done()
			defer func() { <-workersChannel }()

			s.sendFileChunksToChannel(startIdx, serverID, len(serverIDs), fileChunks, fileChunksChan)
		}(i, serverID)
		i++
	}

	wg.Wait()

	return nil
}

func (s *Service) transferFileChunks(ctx context.Context, serverID string, fileChunksChan chan model.FileChunkUpload) {
	const op = "service.file.transferFileChunks"

	log := s.logger.With(
		slog.String("op", op),
		slog.String("serverID", serverID),
	)

	s.serversMap.mu.RLock()
	defer s.serversMap.mu.RUnlock()
	grpcClient := s.serversMap.servers[serverID]
	err := grpcClient.UploadFileChunks(ctx, fileChunksChan)
	if err != nil {
		log.Error("failed to transfer file chunks", sl.Err(err))
	}

}

func (s *Service) sendFileChunksToChannel(startIdx int, serverID string, step int, fileChunks *[]model.FileChunk, fileChunksChan chan model.FileChunkUpload) {
	const op = "service.file.getFileChunkUpload"
	log := s.logger.With(slog.String("op", op))
	defer close(fileChunksChan)

	for startIdx < len(*fileChunks) {
		chunk := (*fileChunks)[startIdx]

		fileChunkUpload, err := s.getFileChunkUpload(&chunk)
		if err != nil {
			log.Error("failed to get file chunk upload", sl.Err(err))
			break
		}
		fileChunksChan <- *fileChunkUpload

		s.markFileChunkUploadedSuccessfully(chunk.UUID, serverID)
		startIdx += step
	}
}

func (s *Service) getFileChunkUpload(fileChunk *model.FileChunk) (*model.FileChunkUpload, error) {
	const op = "service.file.getFileChunkUpload"

	log := s.logger.With(
		slog.String("op", op),
	)

	file, err := os.Open(fileChunk.FilePath)
	if err != nil {
		log.Error("failed to open file", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error("failed to close file", sl.Err(err))
		}
	}(file)

	byteStart := fileChunk.ChunkNumber * fileChunk.ChunkSize
	byteEnd := byteStart + fileChunk.ChunkSize

	_, err = file.Seek(byteStart, 0)
	if err != nil {
		log.Error("failed to seek file", sl.Err(err), slog.Int64("byte start", byteStart), slog.Int64("byte end", byteEnd))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	buffer := make([]byte, fileChunk.ChunkSize)
	n, err := file.Read(buffer)
	if err != nil {
		log.Error("failed to read file", sl.Err(err), slog.Int64("byte start", byteStart), slog.Int64("byte end", byteEnd))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &model.FileChunkUpload{
		UUID: fileChunk.UUID,
		Data: buffer[:n],
	}, nil
}

func (s *Service) markFileChunkUploadedSuccessfully(UUID string, serverID string) {
	s.fileRepo.MarkFileChunkSuccessfullyUploaded(UUID, serverID)
}
