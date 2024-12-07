package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/hahaclassic/databases/09_redis/internal/cache"
	"github.com/hahaclassic/databases/09_redis/internal/models"
	"github.com/hahaclassic/databases/09_redis/internal/storage"
	"github.com/hahaclassic/databases/09_redis/pkg/rndgenre"
	tableoutput "github.com/hahaclassic/databases/09_redis/pkg/table"
	"github.com/jedib0t/go-pretty/table"
)

const (
	getInterval    = 500 * time.Millisecond
	changeInterval = 1000 * time.Millisecond
	limit          = 200
)

type Option int

const (
	DB = iota
	Cache
	CachedDB
)

type Service struct {
	cached cache.Cache
	db     storage.Storage
}

func New(cache cache.Cache, storage storage.Storage) *Service {
	return &Service{
		cached: cache,
		db:     storage,
	}
}

func (s *Service) Top10MostStreamedTracks(option Option, output bool) func(ctx context.Context) error {
	getTracks := map[Option]func(context.Context) ([]*models.Track, error){
		DB: s.db.Top10MostStreamedTracks,

		Cache: func(ctx context.Context) ([]*models.Track, error) {
			return s.cached.Get(ctx, cache.Key)
		},

		CachedDB: s.Top10TracksCached,
	}

	return func(ctx context.Context) error {
		tracks, err := getTracks[option](ctx)
		if err != nil {
			return err
		}

		if output {
			headers := []string{"Name", "Album", "...", "Genre", "Streams"}
			rows := [][]interface{}{}
			for _, track := range tracks {
				rows = append(rows, []interface{}{track.Name, track.Album, "...", track.Genre, track.StreamCount})
			}
			tableoutput.PrintTable(table.StyleColoredDark, headers, rows)
		}

		return nil
	}
}

func (s *Service) Top10TracksCached(ctx context.Context) ([]*models.Track, error) {
	if tracks, err := s.cached.Get(ctx, cache.Key); err == nil {
		return tracks, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		return nil, err
	}

	tracks, err := s.db.Top10MostStreamedTracks(ctx)
	if err != nil {
		return nil, err
	}

	if err = s.cached.Set(ctx, cache.Key, tracks); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (s *Service) AddTrack(ctx context.Context) error {
	track := &models.Track{
		Genre: rndgenre.RandomGenre(),
	}

	var err error
	if track.ID, err = uuid.NewRandom(); err != nil {
		return err
	}

	if err = fake.Struct(track); err != nil {
		return err
	}

	track.Name = track.Name[:len(track.Name)-1]
	track.Album = track.Album[:len(track.Album)-1]

	if err = s.db.AddTrack(ctx, track); err != nil {
		return err
	}

	return s.cached.Delete(ctx, cache.Key)
}

func (s *Service) DeleteRandomTrack(ctx context.Context) error {
	if err := s.db.DeleteRandomTrack(ctx); err != nil {
		return err
	}

	return s.cached.Delete(ctx, cache.Key)
}

func (s *Service) UpdateRandomTrackStreams(ctx context.Context) error {
	if err := s.db.UpdateRandomTrackStreams(ctx); err != nil {
		return err
	}

	return s.cached.Delete(ctx, cache.Key)
}

func (s *Service) Bench(change func(context.Context) error, changeType string) func(context.Context) error {
	return func(ctx context.Context) error {
		tickerDB := time.NewTicker(getInterval)
		tickerCache := time.NewTicker(getInterval)
		tickerChange := time.NewTicker(changeInterval)

		fileCache, err := os.Create(fmt.Sprintf("data/%s_cache.txt", changeType))
		if err != nil {
			return err
		}
		defer fileCache.Close()

		fileDB, err := os.Create(fmt.Sprintf("data/%s_db.txt", changeType))
		if err != nil {
			return err
		}
		defer fileDB.Close()

		i := 0
		begin := time.Now()

		for i < limit*2 {
			select {
			case <-tickerDB.C:
				i++
				log.Println("bench db")

				start := time.Now()
				err := s.Top10MostStreamedTracks(DB, false)(ctx)
				end := time.Now()
				if err != nil {
					slog.Error("top10 db", "err", err)
				}

				fmt.Fprintln(fileDB, end.Sub(begin).Milliseconds(), end.Sub(start).Microseconds())

			case <-tickerCache.C:
				i++
				log.Println("bench cache")

				start := time.Now()
				err := s.Top10MostStreamedTracks(CachedDB, false)(ctx)
				end := time.Now()
				if err != nil {
					slog.Error("top10 cached", "err", err)
				}

				fmt.Fprintln(fileCache, end.Sub(begin).Milliseconds(), end.Sub(start).Microseconds())

			case <-tickerChange.C:
				if change == nil {
					continue
				}

				if err := change(ctx); err != nil {
					slog.Error("change", "err", err)
				} else {
					slog.Info("changed data")
				}
			}
		}

		return nil
	}
}
