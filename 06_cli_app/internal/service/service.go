package service

import "github.com/hahaclassic/databases/06_cli_app/internal/storage"

type MusicService struct {
	storage storage.MusicServiceStorage
}

func NewMusicService(storage storage.MusicServiceStorage) *MusicService {
	return &MusicService{storage: storage}
}
