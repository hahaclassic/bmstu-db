package service

import (
	"context"

	"github.com/hahaclassic/databases/01_init/internal/storage"
)

type MusicService struct {
	storage *storage.MusicServiceStorage
}

func New(storage *storage.MusicServiceStorage) *MusicService {
	return &MusicService{storage: storage}
}

func (m *MusicService) Generate(ctx context.Context, recordsPerTable int) {

}

func (m *MusicService) GenerateCSV(ctx context.Context, pathToFolder string, recordsPerTable int) {

}

func (m *MusicService) generate() {

}

func (m *MusicService) DeleteAll(ctx context.Context) {

}
