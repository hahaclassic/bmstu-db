package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"sync"
	"time"

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
	artist := &models.Artist{
		Genre: randGenre(),
	}
	err := faker.FakeData(artist)
	if err != nil {
		slog.Error("", "error", err)
		return
	}

	if artist.ID, err = uuid.NewRandom(); err != nil {
		slog.Error("", "error", err)
		return
	}

	if err := m.storage.CreateArtist(ctx, artist); err != nil {
		slog.Error("", "error", err)
		return
	}

	m.generateAlbums(ctx, artist, tracks)
}

func (m *MusicService) generateAlbums(ctx context.Context, artist *models.Artist, tracks *mutexslice.Slice[uuid.UUID]) {
	var (
		err         error
		numOfAlbums int32
	)
	for numOfAlbums == 0 {
		numOfAlbums = rand.Int32N(maxAlbumsPerArtist)
	}

	for range numOfAlbums {
		album := &models.Album{
			Genre:       artist.Genre,
			ReleaseDate: randomDateAfter(time.Date(artist.DebutYear, 1, 1, 0, 0, 0, 0, time.UTC)),
		}

		if album.ID, err = uuid.NewRandom(); err != nil {
			slog.Error("", "error", err)
			continue
		}

		if err = faker.FakeData(album); err != nil {
			slog.Error("", "error", err)
			continue
		}

		if err := m.storage.CreateAlbum(ctx, album); err != nil {
			slog.Error("", "error", err)
			continue
		}

		if err := m.storage.AddArtistAlbum(ctx, album.ID, artist.ID); err != nil {
			slog.Error("", "error", err)
			continue
		}

		m.generateTracks(ctx, album, artist.ID, tracks)
	}
}

func (m *MusicService) generateTracks(ctx context.Context, album *models.Album, artistID uuid.UUID, tracks *mutexslice.Slice[uuid.UUID]) {
	var (
		err         error
		numOfTracks int32
	)

	for numOfTracks == 0 {
		numOfTracks = rand.Int32N(maxTracksPerAlbum)
	}

	for i := range int(numOfTracks) {
		track := &models.Track{
			AlbumID:      album.ID,
			OrderInAlbum: i + 1,
			Genre:        album.Genre,
		}

		if track.ID, err = uuid.NewRandom(); err != nil {
			slog.Error("", "error", err)
			continue
		}

		if err = faker.FakeData(track); err != nil {
			slog.Error("", "error", err)
			continue
		}

		if err = m.storage.CreateTrack(ctx, track); err != nil {
			slog.Error("", "error", err)
			continue
		}

		tracks.Add(track.ID)

		if err := m.storage.AddArtistTrack(ctx, track.ID, artistID); err != nil {
			slog.Error("", "error", err)
			continue
		}
	}
}

// Generates data about users, playlists and added track into playlists
func (m *MusicService) generateСonsumersData(ctx context.Context, tracks *mutexslice.Slice[uuid.UUID], numOfUsers int) {
	wg := &sync.WaitGroup{}

	for range numOfUsers {
		wg.Add(1)
		go func(ctx context.Context, sliceOfTracks *mutexslice.Slice[uuid.UUID]) {
			m.generateUserWithPlaylists(ctx, sliceOfTracks)
			wg.Done()
		}(ctx, tracks)
	}

	wg.Wait()
}

func (m *MusicService) generateUserWithPlaylists(ctx context.Context, tracks *mutexslice.Slice[uuid.UUID]) {
	user := &models.User{}
	err := faker.FakeData(user)
	if err != nil {
		slog.Error("", "error", err)
		return
	}

	if user.ID, err = uuid.NewRandom(); err != nil {
		slog.Error("", "error", err)
		return
	}

	user.BirthDate, user.RegistrationDate, user.PremiumExpiration = randomDates()

	if err = m.storage.CreateUser(ctx, user); err != nil {
		slog.Error("", "error", err)
		return
	}

	m.generatePlaylists(ctx, user, tracks)
}

func (m *MusicService) generatePlaylists(ctx context.Context, user *models.User, tracks *mutexslice.Slice[uuid.UUID]) {
	var (
		err            error
		numOfPlaylists int32
	)
	for numOfPlaylists == 0 {
		numOfPlaylists = rand.Int32N(maxPlaylistsPerUser)
	}

	for i := range numOfPlaylists {
		playlist := &models.Playlist{
			LastUpdated: randomDateAfter(user.RegistrationDate),
		}

		if playlist.ID, err = uuid.NewRandom(); err != nil {
			slog.Error("", "error", err)
			continue
		}

		if err := faker.FakeData(&playlist); err != nil {
			slog.Error("", "error", err)
			continue
		}

		if err = m.storage.CreatePlaylist(ctx, playlist); err != nil {
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

		if err = m.storage.AddPlaylist(ctx, userPlaylist); err != nil {
			slog.Error("", "error", err)
			continue
		}

		m.fillPlaylist(ctx, playlist.ID, tracks)
	}
}

func (m *MusicService) fillPlaylist(ctx context.Context, playlistID uuid.UUID, tracks *mutexslice.Slice[uuid.UUID]) {
	var numOfTracks int32
	for numOfTracks == 0 {
		numOfTracks = rand.Int32N(maxTracksPerPlaylist)
	}

	for j := range int(numOfTracks) {
		trackIdx := rand.Int32N(int32(tracks.Len()))

		listTrack := &models.PlaylistTrack{
			ID:         tracks.Get(int(trackIdx)),
			PlaylistID: playlistID,
			TrackOrder: j + 1,
		}

		err := m.storage.AddTrackToPlaylist(ctx, listTrack)
		if err != nil {
			slog.Error("", "error", err)
			continue
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
