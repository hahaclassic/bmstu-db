package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math/rand/v2"
	"sync"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/hahaclassic/databases/01_init/internal/models"
	"github.com/hahaclassic/databases/01_init/internal/storage"
	"github.com/hahaclassic/databases/01_init/pkg/mutexslice"
)

// *******************************************************************
// NOTE: error handling and transmission to the upper level is omitted
// *******************************************************************

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

	tracks := m.generateProducersData(ctx, recordsPerTable)
	if err != nil {
		return err
	}

	m.generateСonsumersData(ctx, tracks, recordsPerTable)
	if err != nil {
		return err
	}

	return nil
}

// Generates data about artists, albums, tracks
func (m *MusicService) generateProducersData(ctx context.Context, numOfArtists int) *mutexslice.Slice[uuid.UUID] {
	tracks := &mutexslice.Slice[uuid.UUID]{}

	wg := &sync.WaitGroup{}

	for range numOfArtists {
		wg.Add(1)
		go func(ctx context.Context, sliceOfTracks *mutexslice.Slice[uuid.UUID]) {
			m.generateArtistWithAlbumsAndTracks(ctx, sliceOfTracks)
			wg.Done()
		}(ctx, tracks)
	}

	wg.Wait()

	return tracks
}

func (m *MusicService) generateArtistWithAlbumsAndTracks(ctx context.Context, tracks *mutexslice.Slice[uuid.UUID]) {
	artist := &models.Artist{}
	err := faker.FakeData(artist)
	if err != nil {
		slog.Error("", "error", err)
		return
	}

	artist.Genre = randGenre()
	artist.ID, err = uuid.NewRandom()
	if err != nil {
		slog.Error("", "error", err)
		return
	}

	if err := m.storage.CreateArtist(ctx, artist); err != nil {
		slog.Error("", "error", err)
		return
	}

	var numOfAlbums int32
	for numOfAlbums == 0 {
		numOfAlbums = rand.Int32N(maxAlbumsPerArtist)
	}

	for range numOfAlbums {
		album := &models.Album{}

		album.ID, err = uuid.NewRandom()
		if err != nil {
			slog.Error("", "error", err)
			continue
		}

		err = faker.FakeData(album)
		if err != nil {
			slog.Error("", "error", err)
			continue
		}

		album.Genre = artist.Genre

		if err := m.storage.CreateAlbum(ctx, album); err != nil {
			slog.Error("", "error", err)
			continue
		}

		if err := m.storage.AddArtistAlbum(ctx, album.ID, artist.ID); err != nil {
			slog.Error("", "error", err)
			continue
		}

		var numOfTracks int32
		for numOfTracks == 0 {
			numOfTracks = rand.Int32N(maxTracksPerAlbum)
		}

		for i := range int(numOfTracks) {
			track := &models.Track{}

			track.ID, err = uuid.NewRandom()
			if err != nil {
				slog.Error("", "error", err)
				continue
			}

			err = faker.FakeData(album)
			album.Genre = artist.Genre
			if err != nil {
				slog.Error("", "error", err)
				continue
			}

			tracks.Add(track.ID)

			listTrack := &models.ListTrack{
				ID:         track.ID,
				ListID:     album.ID,
				TrackOrder: i + 1,
			}

			if err := m.storage.AddTrackToAlbum(ctx, listTrack); err != nil {
				slog.Error("", "error", err)
				continue
			}

			if err := m.storage.AddArtistTrack(ctx, track.ID, artist.ID); err != nil {
				slog.Error("", "error", err)
				continue
			}
		}
	}
}

// Generates data about users, playlists and added track into playlists
func (m *MusicService) generateСonsumersData(ctx context.Context, tracks *mutexslice.Slice[uuid.UUID], numOfUsers int) {
	wg := &sync.WaitGroup{}

	for range numOfUsers {
		wg.Add(1)
		go func(ctx context.Context, sliceOfTracks *mutexslice.Slice[uuid.UUID]) {
			m.generateUsersWithPlaylists(ctx, sliceOfTracks)
			wg.Done()
		}(ctx, tracks)
	}

	wg.Wait()
}

func (m *MusicService) generateUsersWithPlaylists(ctx context.Context, tracks *mutexslice.Slice[uuid.UUID]) {
	user := &models.User{}
	err := faker.FakeData(user)
	if err != nil {
		slog.Error("", "error", err)
		return
	}
	id, err := uuid.NewRandom()
	user.ID = id
	if err != nil {
		slog.Error("", "error", err)
		return
	}

	user.BirthDate, user.RegistrationDate, user.PremiumExpiration = randomDates()

	err = m.storage.CreateUser(ctx, user)
	if err != nil {
		slog.Error("", "error", err)
		return
	}

	var numOfPlaylists int32
	for numOfPlaylists == 0 {
		numOfPlaylists = rand.Int32N(maxPlaylistsPerUser)
	}

	for i := range numOfPlaylists {
		playlist := &models.Playlist{}

		playlist.ID, err = uuid.NewRandom()
		if err != nil {
			slog.Error("", "error", err)
			continue
		}
		if err := faker.FakeData(&playlist); err != nil {
			log.Fatal(err)
		}

		playlist.LastUpdated = randomLastUpdated()

		err = m.storage.CreatePlaylist(ctx, playlist)
		if err != nil {
			slog.Error("", "error", err)
			continue
		}

		userPlaylist := &models.UserPlaylist{
			ID:          playlist.ID,
			UserID:      user.ID,
			AccessLevel: 1, // user
		}

		if i == 0 {
			userPlaylist.IsFavorite = true
		}

		err = m.storage.AddPlaylist(ctx, userPlaylist)
		if err != nil {
			slog.Error("", "error", err)
			continue
		}

		var numOfTracks int32
		for numOfTracks == 0 {
			numOfTracks = rand.Int32N(maxTracksPerPlaylist)
		}

		for j := range int(numOfTracks) {
			trackIdx := rand.Int32N(int32(tracks.Len()))

			listTrack := &models.ListTrack{
				ID:         tracks.Get(int(trackIdx)),
				ListID:     playlist.ID,
				TrackOrder: j + 1,
			}

			err = m.storage.AddTrackToPlaylist(ctx, listTrack)
			if err != nil {
				slog.Error("", "error", err)
				continue
			}
		}
	}
}

func (m *MusicService) DeleteAll(ctx context.Context) error {
	err := m.storage.DeleteAll(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteAll, err)
	}

	return nil
}
