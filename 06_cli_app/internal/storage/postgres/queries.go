package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hahaclassic/databases/06_cli_app/internal/models"
	"github.com/hahaclassic/databases/06_cli_app/internal/storage"
	"github.com/jackc/pgx/v5/pgconn"
)

// 1. Скалярный запрос
// Получение количества исполнителей
func (s *MusicServiceStorage) CountArtists(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRow(ctx, "SELECT COUNT(*) FROM artists").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", storage.ErrCountArtists, err)
	}

	return count, nil
}

// 2. Запрос с несколькими соединениями (JOIN)
// Получение треков с альбомами, артистами, датой релиза, количеством прослушиваний
func (s *MusicServiceStorage) TracksAdditionalInfo(ctx context.Context, limit int) (tracks []*models.TracksAdditionalInfo, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("%w: %w", storage.ErrTracksAdditionalInfo, err)
		}
	}()

	query := `
		SELECT t.name AS track_name,
			al.title AS album_title,
			a.name AS artist_name,
			al.release_date,
			t.stream_count
		FROM tracks t
		JOIN albums al ON t.album_id = al.id
		JOIN albums_by_artists ab ON al.id = ab.album_id
		JOIN artists a ON ab.artist_id = a.id
		ORDER BY t.stream_count DESC LIMIT $1`

	rows, err := s.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		curr := &models.TracksAdditionalInfo{}
		if err := rows.Scan(&curr.Name, &curr.Album, &curr.Artist,
			&curr.ReleaseDate, &curr.StreamCount); err != nil {
			return nil, err
		}
		tracks = append(tracks, curr)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return tracks, nil
}

// 3. Запрос с ОТВ (CTE) и оконными функциями
// Получение 5 самых популярных треков (название, кол-во прослушиваний, ранг)
func (s *MusicServiceStorage) BestTracks(ctx context.Context, limit int) (tracks []*models.RankedTrack, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("%w: %w", storage.ErrRankedTracks, err)
		}
	}()

	query := `
        WITH RankedTracks AS (
            SELECT name, stream_count, ROW_NUMBER() OVER (ORDER BY stream_count DESC) AS rank
            FROM tracks
        )
        SELECT name, stream_count, rank FROM RankedTracks WHERE rank <= $1`

	rows, err := s.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		curr := &models.RankedTrack{}
		if err := rows.Scan(&curr.Name, &curr.StreamCount, &curr.Rank); err != nil {
			return nil, err
		}
		tracks = append(tracks, curr)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return tracks, nil
}

// 4. Запрос к метаданным
// Получение списка названий таблиц
func (s *MusicServiceStorage) TablesNames(ctx context.Context) (names []string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("%w: %w", storage.ErrTablesNames, err)
		}
	}()

	rows, err := s.db.Query(context.Background(), "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		names = append(names, tableName)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return names, nil
}

// 5. Вызов скалярной функции
// Получение количества треков исполнителя
func (s *MusicServiceStorage) CountArtistTracks(ctx context.Context, artistID uuid.UUID) (int, error) {
	var count int
	err := s.db.QueryRow(ctx, "SELECT count_artist_tracks($1)", artistID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", storage.ErrCountArtistTracks, err)
	}

	return count, nil
}

// 6. Вызов табличной функции
// Получение информации о альбомах исполнителя
func (s *MusicServiceStorage) ArtistAlbums(ctx context.Context, artistID uuid.UUID) (albums []*models.Album, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("%w: %w", storage.ErrAlbumsInfo, err)
		}
	}()

	rows, err := s.db.Query(ctx,
		`SELECT album_id, title, release_date, label, genre
    	FROM get_albums_by_artist($1)`, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		curr := &models.Album{}
		if err := rows.Scan(&curr.ID, &curr.Title, &curr.ReleaseDate,
			&curr.Label, &curr.Genre); err != nil {
			return nil, err
		}
		albums = append(albums, curr)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return albums, nil
}

// 7. Вызов хранимой процедуры без параметров
// Процедура для выдачи премиума всем пользователям от 65 лет.
func (s *MusicServiceStorage) FreePremiumForPensioners(ctx context.Context) error {
	_, err := s.db.Exec(ctx, "call free_premium_to_pensioners()")
	if err != nil {
		return fmt.Errorf("%w: %w", storage.ErrFreePremiumForPensioners, err)
	}

	return nil
}

// 8. Вызов системную функцию или процедуру
// Получение текущего пользователя
func (s *MusicServiceStorage) CurrentUser(ctx context.Context) (string, error) {
	var currentUser string
	err := s.db.QueryRow(ctx, "SELECT current_user").Scan(&currentUser)
	if err != nil {
		return "", fmt.Errorf("%w: %w", storage.ErrCurrentUser, err)
	}

	return currentUser, nil
}

// 9. Создание таблицы в базе данных
// Создание таблицы ревью на альбомы
func (s *MusicServiceStorage) CreateReviewsTable(ctx context.Context) error {
	query := `CREATE TABLE reviews (
		id UUID PRIMARY KEY,                  
		user_id UUID,                         
		album_id UUID,                         
		rating INT CHECK (rating >= 1 AND rating <= 10),
		comment TEXT,                          
		review_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE, 
		FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE)`

	_, err := s.db.Exec(ctx, query)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "42P07" {
			return fmt.Errorf("%w: %w", storage.ErrTableAlreadyExists, err)
		}
		return fmt.Errorf("%w: %w", storage.ErrCreateTable, err)
	}

	return nil
}

// 10. Вставка данных в созданную таблицу с помощью INSERT
// Добавление ревью
func (s *MusicServiceStorage) AddReview(ctx context.Context, review *models.Review) error {
	query := `INSERT INTO reviews (id, user_id, album_id, rating, comment) VALUES 
		($1, $2, $3, $4, $5)`

	_, err := s.db.Exec(ctx, query, review.ID, review.UserID, review.AlbumID,
		review.Rating, review.Comment)
	if err != nil {
		return fmt.Errorf("%w: %w", storage.ErrAddReview, err)
	}

	return nil
}
