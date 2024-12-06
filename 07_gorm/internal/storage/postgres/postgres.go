package postgres

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/databases/07_gorm/config"
	"github.com/hahaclassic/databases/07_gorm/internal/models"
	"github.com/hahaclassic/databases/07_gorm/internal/storage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func New(ctx context.Context, config *config.PostgresConfig) (*Storage, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		config.User, config.Password, net.JoinHostPort(config.Host, config.Port), config.DB, config.SSLMode)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrStorageConnection, err)
	}

	return &Storage{db}, nil
}

// Самые прослушиваемые limit треков
func (s *Storage) BestExplicitTracks(ctx context.Context, limit int) ([]*models.Track, error) {
	tracks := []*models.Track{}

	if res := s.db.WithContext(ctx).Where("explicit = ?", true).
		Order("stream_count DESC").Limit(limit).Find(&tracks); res.Error != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrExplicitTracks, res.Error)
	}

	return tracks, nil
}

// Подсчет количества треков по жанрам
func (s *Storage) CountTracksByGenre(ctx context.Context) ([]*models.GenreCount, error) {
	genres := []*models.GenreCount{}

	if res := s.db.WithContext(ctx).Model(&models.GenreCount{}).Table("tracks").
		Select("genre, count(*) as count").Group("genre").Order("count DESC").
		Find(&genres); res.Error != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrCountTracksByGenre, res.Error)
	}

	return genres, nil
}

// Получение списка альбомов, которые содержат не более maxNumOfTracks треков
func (s *Storage) AlbumsWithMaxTracks(ctx context.Context, minNumOfTracks int) ([]*models.Album, error) {
	albums := []*models.Album{}

	if err := s.db.WithContext(ctx).
		Joins("JOIN tracks t ON t.album_id = albums.id").
		Group("albums.id").
		Having("COUNT(t.id) <= ?", minNumOfTracks).
		Find(&albums).Error; err != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrAlbumsWithMaxTracks, err)
	}

	return albums, nil
}

// Получение списка исполнителей,
func (s *Storage) ArtistsWithReleasedAlbumYear(ctx context.Context, year int) ([]*models.Artist, error) {
	start, end := fmt.Sprintf("%d-01-01", year), fmt.Sprintf("%d-12-31", year)

	var artists []*models.Artist
	if err := s.db.WithContext(ctx).
		Joins("JOIN albums_by_artists a ON a.artist_id = artists.id").
		Joins("JOIN albums al ON al.id = a.album_id").
		Where("al.release_date BETWEEN ? AND ?", start, end).
		Order("artists.genre ASC, artists.debut_year DESC").
		Find(&artists).Error; err != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrArtistsWithReleasedAlbumYear, err)
	}

	return artists, nil
}

// Получение списка всех пользователей, старше указанного возраста
func (s *Storage) UsersOlderThan(ctx context.Context, age int) ([]*models.User, error) {
	var users []*models.User

	if err := s.db.WithContext(ctx).
		Model(&models.User{}).
		Where("DATE_PART('year', AGE(birth_date)) > ?", age).
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrUsersOlderThan, err)
	}

	return users, nil
}

// Однотабличный запрос
// Получение всех треков одного жанра
func (s *Storage) TracksByGenre(ctx context.Context, genre string) ([]*models.Track, error) {
	var tracks []*models.Track
	if err := s.db.WithContext(ctx).Where("genre = ?", genre).Find(&tracks).Error; err != nil {
		return nil, err
	}

	return tracks, nil
}

// Многотабличный запрос
func (s *Storage) AlbumsWithTrackCounts(ctx context.Context, genre string) ([]*models.AlbumTrackCount, error) {
	var results []*models.AlbumTrackCount

	err := s.db.WithContext(ctx).Table("albums").
		Select("albums.id as album_id, albums.title, COUNT(tracks.id) as track_count").
		Joins("LEFT JOIN tracks ON albums.id = tracks.album_id").
		Where("albums.genre = ?", genre).
		Group("albums.id").
		Scan(&results).Error

	return results, err
}

// Добавление пользователя
func (s *Storage) AddUser(ctx context.Context, user *models.User) error {
	return s.db.WithContext(ctx).Create(&user).Error
}

// Обновление username
func (s *Storage) UpdateUserName(ctx context.Context, userID uuid.UUID, newName string) error {
	return s.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).Update("name", newName).Error
}

// Удаление пользователя
func (s *Storage) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.db.WithContext(ctx).Delete(&models.User{}, userID).Error
}

// Получение альбомов исполнителя
func (s *Storage) AlbumsByArtist(ctx context.Context, artistID uuid.UUID) ([]*models.Album, error) {
	var albums []*models.Album

	if err := s.db.Raw(`SELECT album_id as id, title, release_date, label, 
						genre FROM get_albums_by_artist(?)`, artistID).
		Scan(&albums).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch albums by artist: %w", err)
	}

	return albums, nil
}

func (s *Storage) ExportUsersToJSON(ctx context.Context) ([]byte, error) {
	var result string
	query := `SELECT json_agg(row_to_json(users)) FROM users`

	if err := s.db.WithContext(ctx).Raw(query).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to export users to JSON: %w", err)
	}

	return []byte(result), nil
}

func (s *Storage) ImportUsers(ctx context.Context, users []*models.User) error {
	defaultUsers := make([]struct {
		ID                uuid.UUID `json:"id"`
		Name              string    `json:"name"`
		RegistrationDate  time.Time `json:"registration_date"`
		BirthDate         time.Time `json:"birth_date"`
		Premium           bool      `json:"premium"`
		PremiumExpiration time.Time `json:"premium_expiration"`
	}, len(users))

	for i := range users {
		defaultUsers[i].ID = users[i].ID
		defaultUsers[i].Name = users[i].Name
		defaultUsers[i].RegistrationDate = users[i].RegistrationDate
		defaultUsers[i].BirthDate = time.Time(users[i].BirthDate)
		defaultUsers[i].Premium = users[i].Premium
		defaultUsers[i].PremiumExpiration = users[i].PremiumExpiration
	}

	if err := s.db.WithContext(ctx).Table("temp_users").Create(&defaultUsers).Error; err != nil {
		return fmt.Errorf("failed to insert users into the database: %w", err)
	}

	return nil
}
