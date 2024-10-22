package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/hahaclassic/databases/01_init/internal/models"
	"github.com/hahaclassic/databases/01_init/internal/storage"
	"github.com/hahaclassic/databases/01_init/pkg/mutex"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrExceededContextTime error = errors.New("exceeded context time")
	ErrGenerateData        error = errors.New("failed to generate data")
	ErrDeleteAll           error = errors.New("failed to delete all data")
	ErrDuplicate           error = errors.New("duplicate key")
)

const (
	DuplicateSQLCode = "23505"
)

type UniqueController struct {
	artistNames *mutex.Collection[string]
	tracks      *mutex.Slice[uuid.UUID]
	users       atomic.Int64
	playlists   atomic.Int64
	albums      atomic.Int64
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

func (m *MusicService) Generate(ctx context.Context, recordsPerTable int) error {
	_ = fake.Seed(0)

	// Generates data about artists, albums, tracks
	err := m.generate(ctx, recordsPerTable, m.generateArtistWithAlbumsAndTracks)
	if err != nil {
		return err
	}

	// Generates data about users, playlists and added track into playlists
	err = m.generate(ctx, recordsPerTable, m.generateUserWithPlaylists)
	if err != nil {
		return err
	}

	fmt.Println("\nRESULT",
		"\n- artists:", m.uniq.artistNames.Len(),
		"\n- albums:", m.uniq.albums.Load(),
		"\n- tracks:", m.uniq.tracks.Len(),
		"\n- users:", m.uniq.users.Load(),
		"\n- playlists:", m.uniq.playlists.Load())

	return nil
}

func (m *MusicService) generate(ctx context.Context, num int, generator func(context.Context) error) error {
	workers := runtime.GOMAXPROCS(0)
	numPerWorker := int(float64(num) / float64(workers))

	if workers >= num {
		workers = num
		numPerWorker = 1
	}

	wg := &sync.WaitGroup{}
	errChan := make(chan error)

	for i := range workers {
		wg.Add(1)
		go func(ctx context.Context, idx int) {
			for range numPerWorker {
				errChan <- generator(ctx)
			}
			wg.Done()
		}(ctx, i)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if errors.Is(err, ErrDuplicate) {
			slog.Error("DUPLICATE", "err", err)
		} else if err != nil {
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

	var perr *pgconn.PgError
	err = m.storage.CreateArtist(ctx, artist)
	if ok := errors.As(err, &perr); ok && perr.Code == DuplicateSQLCode {
		return fmt.Errorf("%w: artist '%s' already exists", ErrDuplicate, artist.Name)
	}
	if err != nil {
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
	words := strings.Split(name, " ")
	for (len(words) == 1 && len(words[0]) < 5) || m.uniq.artistNames.Contains(name) {
		name = fake.Sentence(int(numOfWords))
		name = name[:len(name)-1]
		words = strings.Split(name, " ")
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

		var perr *pgconn.PgError
		err = m.storage.CreateAlbum(ctx, album)
		if ok := errors.As(err, &perr); ok && perr.Code == DuplicateSQLCode {
			continue
		}
		if err != nil {
			return fmt.Errorf("%w: err with album %v", err, album)
		}

		m.uniq.albums.Add(1)

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
		numOfTracks int32 = rand.Int32N(maxTracksPerAlbum)
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

		var perr *pgconn.PgError
		err = m.storage.CreateTrack(ctx, track)
		if ok := errors.As(err, &perr); ok && perr.Code == DuplicateSQLCode {
			continue
		}
		if err != nil {
			return fmt.Errorf("%w: err with track %v", err, track)
		}

		m.uniq.tracks.Add(track.ID)

		if err := m.storage.AddArtistTrack(ctx, track.ID, artistID); err != nil {
			return fmt.Errorf("%w: err with artistTrack %s, %s", err, track.ID, artistID)
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

		var perr *pgconn.PgError
		err = m.storage.CreateUser(ctx, user)
		if ok := errors.As(err, &perr); ok && perr.Code == DuplicateSQLCode {
			errChan <- fmt.Errorf("%w: user '%s' already exists", ErrDuplicate, user.ID)
			return
		}
		if err != nil {
			errChan <- fmt.Errorf("%w: err with user %v", err, user)
			return
		}

		m.uniq.users.Add(1)

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
		numOfPlaylists int32 = rand.Int32N(maxPlaylistsPerUser)
	)

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

		var perr *pgconn.PgError
		err = m.storage.CreatePlaylist(ctx, playlist)
		if ok := errors.As(err, &perr); ok && perr.Code == DuplicateSQLCode {
			continue
		}
		if err != nil {
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

		err = m.storage.AddPlaylist(ctx, userPlaylist)
		if ok := errors.As(err, &perr); ok && perr.Code == DuplicateSQLCode {
			slog.Error("DUPLICATE", "err", err)
			continue
		}
		if err != nil {
			return fmt.Errorf("%w: err with userPlaylist: %v", err, userPlaylist)
		}

		m.uniq.playlists.Add(1)

		if err = m.fillPlaylist(ctx, playlist.ID); err != nil {
			return fmt.Errorf("%w: err with fill playlist: %v", err, playlist)
		}
	}

	return nil
}

func (m *MusicService) fillPlaylist(ctx context.Context, playlistID uuid.UUID) error {
	numOfTracks := rand.Int32N(maxTracksPerPlaylist)
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

		var perr *pgconn.PgError
		err := m.storage.AddTrackToPlaylist(ctx, playlistTrack)
		if ok := errors.As(err, &perr); ok && perr.Code == DuplicateSQLCode {
			continue
		}
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

	slog.Info("[OK]: All data deleted.")

	return nil
}
