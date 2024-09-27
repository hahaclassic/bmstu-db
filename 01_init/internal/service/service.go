package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/hahaclassic/databases/01_init/internal/models"
	"github.com/hahaclassic/databases/01_init/internal/storage"
)

var (
	ErrGenerateData error = errors.New("failed to generate data")
	ErrDeleteAll    error = errors.New("failed to delete all data")
)

type MusicService struct {
	storage storage.MusicServiceStorage
}

func New(storage storage.MusicServiceStorage) *MusicService {
	return &MusicService{storage: storage}
}

func (m *MusicService) Generate(ctx context.Context, recordsPerTable int) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("%w: %v", ErrGenerateData, err)
		}
	}()

	err = m.generateProducersData(ctx, recordsPerTable)
	if err != nil {
		return err
	}

	err = m.generateСonsumersData(ctx, recordsPerTable)
	if err != nil {
		return err
	}

	return nil
}

func (m *MusicService) generateProducersData(ctx context.Context, recordsPerTable int) error {
	for range recordsPerTable {
		artist := &models.Artist{}
		err := faker.FakeData(artist)
		if err != nil {
			return err
		}

		artist.Genre = randGenre()
		artist.ID = uuid.New()

		if err := m.storage.CreateArtist(ctx, artist); err != nil {
			return err
		}

		for range albumsPerArtist {
			album := &models.Album{}
			err = faker.FakeData(album)
			album.Genre = artist.Genre
			if err != nil {
				return err
			}

			if err := m.storage.CreateAlbum(ctx, album); err != nil {
				log.Printf("Error creating album: %v", err)
			}

			for range tracksPerAlbum {
				track := &models.Track{}
				err = faker.FakeData(album)
				album.Genre = artist.Genre
				if err != nil {
					return err
				}

				if err := m.storage.AddTrackToAlbum(ctx, track.ID, album.ID); err != nil {
					log.Printf("Error adding track to album: %v", err)
				}

				if err := m.storage.AddArtistTrack(ctx, track.ID, artist.ID); err != nil {
					log.Printf("Error adding track to artist: %v", err)
				}
			}

		}
	}

	return nil
}

func (m *MusicService) generateСonsumersData(ctx context.Context, tracks []uuid.UUID, n int) {

}

func (m *MusicService) DeleteAll(ctx context.Context) error {
	err := m.storage.DeleteAll(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteAll, err)
	}

	return nil
}

func (m *MusicService) generateRandomData(storage storage.MusicServiceStorage) {
	ctx := context.Background()

	for i := 0; i < numArtists; i++ {
		artist := &models.Artist{}
		faker.FakeData(artist)
		artist.ID = uuid.New()
		if err := storage.CreateArtist(ctx, artist); err != nil {
			log.Printf("Error creating artist: %v", err)
		}
	}

	// Храним ссылки на созданных артистов
	artistIDs := make([]uuid.UUID, numArtists)

	for i := 0; i < numAlbums; i++ {
		album := &models.Album{}
		faker.FakeData(album)
		album.ID = uuid.New()
		album.ReleaseDate = time.Now()
		if err := storage.CreateAlbum(ctx, album); err != nil {
			log.Printf("Error creating album: %v", err)
		}
		// Добавляем связь между альбомами и артистами
		artistID := artistIDs[faker.Number(0, numArtists-1)] // Случайный артист
		if err := storage.AddTrackToAlbum(ctx, album.ID, artistID); err != nil {
			log.Printf("Error adding album to artist: %v", err)
		}
	}

	for i := 0; i < numTracks; i++ {
		track := &models.Track{}
		faker.FakeData(track)
		track.ID = uuid.New()
		if err := storage.CreateTrack(ctx, track); err != nil {
			log.Printf("Error creating track: %v", err)
		}
	}

	for i := 0; i < numPlaylists; i++ {
		playlist := &models.Playlist{}
		faker.FakeData(playlist)
		playlist.ID = uuid.New()
		playlist.LastUpdated = time.Now()
		if err := storage.CreatePlaylist(ctx, playlist); err != nil {
			log.Printf("Error creating playlist: %v", err)
		}
	}

	// Добавляем треки в плейлисты
	for i := 0; i < numPlaylists; i++ {
		playlistID := playlist.ID
		// Генерируем случайное количество треков для каждого плейлиста
		numTracksInPlaylist := faker.Number(1, 5)
		for j := 0; j < numTracksInPlaylist; j++ {
			trackID := uuid.New() // Случайный трек
			if err := storage.AddTrackToPlaylist(ctx, trackID, playlistID); err != nil {
				log.Printf("Error adding track to playlist: %v", err)
			}
		}
	}

	for i := 0; i < numUsers; i++ {
		user := &models.User{}
		faker.FakeData(user)
		user.ID = uuid.New()
		user.RegistrationDate = time.Now()
		if err := storage.CreateUser(ctx, user); err != nil {
			log.Printf("Error creating user: %v", err)
		}
	}
}
