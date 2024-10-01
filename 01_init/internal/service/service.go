package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"runtime"
	"sync"
	"time"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/hahaclassic/databases/01_init/internal/models"
	"github.com/hahaclassic/databases/01_init/internal/storage"
	"github.com/hahaclassic/databases/01_init/pkg/mutex"
)

var (
	ErrExceededContextTime error = errors.New("exceeded context time")
	ErrGenerateData        error = errors.New("failed to generate data")
	ErrDeleteAll           error = errors.New("failed to delete all data")
)

type UniqueController struct {
	artistNames *mutex.Collection[string]
	tracks      *mutex.Slice[uuid.UUID]
}

type MusicService struct {
	uniq    *UniqueController
	storage storage.MusicServiceStorage
}

func New(storage storage.MusicServiceStorage) *MusicService {
	return &MusicService{
		uniq: &UniqueController{
			tracks:      mutex.NewSlice[uuid.UUID](),
			artistNames: mutex.NewCollection[string](),
		},
		storage: storage,
	}
}

func (m *MusicService) Generate(ctx context.Context, recordsPerTable int) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("%w: %v", ErrGenerateData, err)
		}
	}()

	_ = fake.Seed(0)

	err = m.generateProducersData(ctx, recordsPerTable)
	if err != nil {
		return err
	}

	err = m.generateСonsumersData(ctx, recordsPerTable)
	if err != nil {
		return err
	}

	slog.Info("RESULT", "artists", m.uniq.artistNames.Len(), "tracks", m.uniq.tracks.Len())

	return nil
}

// Generates data about artists, albums, tracks
func (m *MusicService) generateProducersData(ctx context.Context, numOfArtists int) error {
	workers := runtime.GOMAXPROCS(0)
	artistsPerWorker := int(float64(numOfArtists) / float64(workers))

	if workers >= numOfArtists {
		workers = numOfArtists
		artistsPerWorker = 1
	}

	wg := &sync.WaitGroup{}
	errChan := make(chan error)

	for i := range workers {
		wg.Add(1)
		go func(ctx context.Context, idx int) {
			for range artistsPerWorker {
				errChan <- m.generateArtistWithAlbumsAndTracks(ctx)
			}
			wg.Done()
		}(ctx, i)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MusicService) generateArtistWithAlbumsAndTracks(ctx context.Context) error {
	var err error

	artist := &models.Artist{
		Genre:     randomGenre(),
		Country:   fake.Country(),
		DebutYear: randomDateAfter(time.Date(1920, 1, 1, 0, 0, 0, 0, time.UTC)).Year(),
		Name:      m.randomArtistName(),
	}

	artist.ID, err = uuid.NewRandom()
	if err != nil {
		return err
	}

	if err := m.storage.CreateArtist(ctx, artist); err != nil {
		return fmt.Errorf("%w: err with artist %v", err, artist)
	}

	m.uniq.artistNames.Store(artist.Name)

	return m.generateAlbums(ctx, artist)
}

func (m *MusicService) randomArtistName() string {
	numOfWords := rand.Int32N(5)
	for numOfWords == 0 {
		numOfWords = rand.Int32N(5)
	}

	name := fake.Sentence(int(numOfWords))
	name = name[:len(name)-1]
	for m.uniq.artistNames.Contains(name) {
		name = fake.Sentence(int(numOfWords))
		name = name[:len(name)-1]
	}

	return name
}

func (m *MusicService) generateAlbums(ctx context.Context, artist *models.Artist) error {
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
			return err
		}

		if err = fake.Struct(album); err != nil {
			return err
		}
		album.Title = album.Title[:len(album.Title)-1]
		album.Label = album.Label[:len(album.Label)-1]

		if err := m.storage.CreateAlbum(ctx, album); err != nil {
			return fmt.Errorf("%w: err with album %v", err, album)
		}

		if err := m.storage.AddArtistAlbum(ctx, album.ID, artist.ID); err != nil {
			return fmt.Errorf("%w: err with artistAlbum %s %s", err, album.ID, artist.ID)
		}

		if err = m.generateTracks(ctx, album, artist.ID); err != nil {
			return err
		}
	}

	return nil
}

func (m *MusicService) generateTracks(ctx context.Context, album *models.Album, artistID uuid.UUID) error {
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
			return err
		}

		if err = fake.Struct(track); err != nil {
			return err
		}

		track.Name = track.Name[:len(track.Name)-1]

		if err = m.storage.CreateTrack(ctx, track); err != nil {
			return fmt.Errorf("%w: err with track %v", err, track)
		}

		m.uniq.tracks.Add(track.ID)

		if err := m.storage.AddArtistTrack(ctx, track.ID, artistID); err != nil {
			return fmt.Errorf("%w: err with artistTrack %s, %s", err, track.ID, artistID)
		}
	}

	return nil
}

// Generates data about users, playlists and added track into playlists
func (m *MusicService) generateСonsumersData(ctx context.Context, numOfUsers int) error {
	workers := runtime.GOMAXPROCS(0)
	usersPerWorker := int(float64(numOfUsers) / float64(workers))

	if workers >= numOfUsers {
		workers = numOfUsers
		usersPerWorker = 1
	}

	wg := &sync.WaitGroup{}
	errChan := make(chan error)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	for i := range workers {
		wg.Add(1)
		go func(ctx context.Context, idx int) {
			for range usersPerWorker {
				errChan <- m.generateUserWithPlaylists(ctx)
			}
			wg.Done()
		}(ctx, i)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MusicService) generateUserWithPlaylists(ctx context.Context) error {
	errChan := make(chan error)

	generateUser := func(ctx context.Context, errChan chan error) {
		user := &models.User{}
		err := fake.Struct(user)
		if err != nil {
			errChan <- err
			return
		}

		if user.ID, err = uuid.NewRandom(); err != nil {
			errChan <- err
			return
		}

		user.BirthDate, user.RegistrationDate, user.PremiumExpiration = randomDates()

		if err = m.storage.CreateUser(ctx, user); err != nil {
			errChan <- fmt.Errorf("%w: err with user %v", err, user)
			return
		}

		errChan <- m.generatePlaylists(ctx, user)
	}

	go generateUser(ctx, errChan)

	for {
		select {
		case <-ctx.Done():
			return ErrExceededContextTime

		case err := <-errChan:
			return err
		}
	}
}

func (m *MusicService) generatePlaylists(ctx context.Context, user *models.User) error {
	var (
		err            error
		numOfPlaylists int32
	)
	for numOfPlaylists == 0 {
		numOfPlaylists = rand.Int32N(maxPlaylistsPerUser)
	}

	for i := range numOfPlaylists {
		// playlist := &models.Playlist{
		// 	LastUpdated: randomDateAfter(user.RegistrationDate),
		// }
		playlist := &models.Playlist{}

		if playlist.ID, err = uuid.NewRandom(); err != nil {
			return err
		}

		if err := fake.Struct(&playlist); err != nil {
			return fmt.Errorf("%w: err with playlist %v", err, playlist)
		}

		if err = m.storage.CreatePlaylist(ctx, playlist); err != nil {
			return fmt.Errorf("%w: err with playlist %v", err, playlist)
		}

		userPlaylist := &models.UserPlaylist{
			ID:          playlist.ID,
			UserID:      user.ID,
			AccessLevel: models.Owner,
		}

		if i == 0 {
			userPlaylist.IsFavorite = true
		}

		if err = m.storage.AddPlaylist(ctx, userPlaylist); err != nil {
			return fmt.Errorf("%w: err with userPlaylist: %v", err, userPlaylist)
		}

		if err = m.fillPlaylist(ctx, playlist.ID); err != nil {
			return fmt.Errorf("%w: err with fill playlist: %v", err, playlist)
		}
	}

	return nil
}

func (m *MusicService) fillPlaylist(ctx context.Context, playlistID uuid.UUID) error {
	var numOfTracks int32
	for numOfTracks == 0 {
		numOfTracks = rand.Int32N(maxTracksPerPlaylist)
	}
	previosIdx := make(map[int32]struct{})

	for j := range int(numOfTracks) {
		trackIdx := rand.Int32N(int32(m.uniq.tracks.Len()))

		_, ok := previosIdx[trackIdx]
		for ok {
			trackIdx = rand.Int32N(int32(m.uniq.tracks.Len()))
			_, ok = previosIdx[trackIdx]
		}

		playlistTrack := &models.PlaylistTrack{
			ID:         m.uniq.tracks.Get(int(trackIdx)),
			PlaylistID: playlistID,
			TrackOrder: j + 1,
		}

		err := m.storage.AddTrackToPlaylist(ctx, playlistTrack)
		if err != nil {
			return fmt.Errorf("%w: err with playlistTrack: %v", err, playlistTrack)
		}

		previosIdx[trackIdx] = struct{}{}
	}

	return nil
}

func (m *MusicService) DeleteAll(ctx context.Context) error {
	err := m.storage.DeleteAll(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteAll, err)
	}

	return nil
}
