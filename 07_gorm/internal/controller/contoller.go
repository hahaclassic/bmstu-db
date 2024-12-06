package controller

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"log/slog"

	"github.com/google/uuid"
	"github.com/hahaclassic/databases/07_gorm/internal/models"
	"github.com/hahaclassic/databases/07_gorm/internal/storage"
	tableoutput "github.com/hahaclassic/databases/07_gorm/pkg/table"
	"github.com/jedib0t/go-pretty/table"
)

type Controller struct {
	storage storage.Storage
}

func NewController(storage storage.Storage) *Controller {
	return &Controller{storage: storage}
}

// Start обрабатывает выбор операции и выполняет соответствующий метод
func (c *Controller) Start(ctx context.Context) {
	methods := map[Operation]func(context.Context){
		BestExplicitTracks:           c.BestExplicitTracks,
		CountTracksByGenre:           c.CountTracksByGenre,
		AlbumsWithMaxTracks:          c.AlbumsWithMaxTracks,
		ArtistsWithReleasedAlbumYear: c.ArtistsWithReleasedAlbumYear,
		UsersOlderThan:               c.UsersOlderThan,
		TracksByGenre:                c.TracksByGenre,
		AlbumsWithTrackCounts:        c.AlbumsWithTrackCounts,
		AddUser:                      c.AddUser,
		UpdateUserName:               c.UpdateUserName,
		DeleteUser:                   c.DeleteUser,
		AlbumsByArtist:               c.AlbumsByArtist,
		ExportUsersToJSON:            c.ExportUsersToJSON,
		ImportUsersFromJSON:          c.ImportUsersFromJSON,
	}

	for {
		menu()
		fmt.Print("Enter operation number (or 0 for exit): ")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println("Invalid operation number. Try again.")
		}

		operation, err := strconv.Atoi(input)
		if err != nil || operation < int(Exit) || operation >= int(operationsEnd) {
			fmt.Println("Invalid data. Try again.")
			continue
		}
		if operation == int(Exit) {
			fmt.Println("Exit")
			return
		}

		methods[Operation(operation)](ctx)
	}
}

func (c *Controller) BestExplicitTracks(ctx context.Context) {
	var limit int
	fmt.Print("Enter limit for best tracks: ")
	_, err := fmt.Scan(&limit)
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}

	tracks, err := c.storage.BestExplicitTracks(ctx, limit)
	if err != nil {
		log.Printf("Error getting best tracks: %v", err)
		return
	}

	if len(tracks) == 0 {
		fmt.Println("No tracks found.")
		return
	}

	headers := []string{"ID", "Name", "...", "Genre", "Stream Count"}
	rows := [][]interface{}{}
	for _, track := range tracks {
		rows = append(rows, []interface{}{track.ID, track.Name, "...", track.Genre, track.StreamCount})
	}

	tableoutput.PrintTable(table.StyleColoredDark, headers, rows)
}

func (c *Controller) CountTracksByGenre(ctx context.Context) {
	genres, err := c.storage.CountTracksByGenre(ctx)
	if err != nil {
		log.Printf("Error getting best tracks: %v", err)
		return
	}
	if len(genres) == 0 {
		fmt.Println("No tracks found.")
		return
	}

	headers := []string{"Genre", "Tracks count"}
	rows := [][]interface{}{}
	for _, genre := range genres {
		rows = append(rows, []interface{}{genre.Genre, genre.Count})
	}

	tableoutput.PrintTable(table.StyleColoredDark, headers, rows)
}

func (c *Controller) AlbumsWithMaxTracks(ctx context.Context) {
	var minNumOfTracks int
	fmt.Print("Enter max number for tracks: ")
	_, err := fmt.Scan(&minNumOfTracks)
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}

	albums, err := c.storage.AlbumsWithMaxTracks(ctx, minNumOfTracks)
	if err != nil {
		log.Printf("Error getting albums: %v", err)
		return
	}

	if len(albums) == 0 {
		fmt.Println("No albums found.")
		return
	}

	headers := []string{"ID", "Title", "Release Date", "Label", "Genre"}
	rows := [][]interface{}{}
	for _, album := range albums {
		rows = append(rows, []interface{}{album.ID, album.Title, album.ReleaseDate, album.Label, album.Genre})
	}

	tableoutput.PrintTable(table.StyleColoredDark, headers, rows)
}

func (c *Controller) ArtistsWithReleasedAlbumYear(ctx context.Context) {
	var year int
	fmt.Print("Enter the release year: ")
	_, err := fmt.Scan(&year)
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}

	artists, err := c.storage.ArtistsWithReleasedAlbumYear(ctx, year)
	if err != nil {
		log.Printf("Error getting artists with albums released in %d: %v", year, err)
		return
	}

	if len(artists) == 0 {
		fmt.Println("No artists found.")
		return
	}

	headers := []string{"ID", "Name", "Genre", "Debut Year"}
	rows := [][]interface{}{}
	for _, artist := range artists {
		rows = append(rows, []interface{}{artist.ID, artist.Name, artist.Genre, artist.DebutYear})
	}

	tableoutput.PrintTable(table.StyleColoredDark, headers, rows)
}

func (c *Controller) UsersOlderThan(ctx context.Context) {
	var age int
	fmt.Print("Enter the minimum age: ")
	_, err := fmt.Scan(&age)
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}

	users, err := c.storage.UsersOlderThan(ctx, age)
	if err != nil {
		log.Printf("Error getting users older than %d: %v", age, err)
		return
	}

	if len(users) == 0 {
		fmt.Println("No users found.")
		return
	}

	headers := []string{"ID", "Name", "Registration", "Birth Date", "Premium", "Prem. Exp"}
	rows := [][]interface{}{}
	for _, user := range users {
		rows = append(rows, []interface{}{user.ID, user.Name, user.RegistrationDate,
			user.BirthDate, user.Premium, user.PremiumExpiration})
	}

	tableoutput.PrintTable(table.StyleColoredDark, headers, rows)
}

// Получение всех треков одного жанра
func (c *Controller) TracksByGenre(ctx context.Context) {
	var genre string
	fmt.Print("Enter genre: ")
	_, err := fmt.Scan(&genre)
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}

	tracks, err := c.storage.TracksByGenre(ctx, genre)
	if err != nil {
		log.Printf("Error getting tracks by genre: %v", err)
		return
	}

	if len(tracks) == 0 {
		fmt.Println("No tracks found.")
		return
	}

	headers := []string{"ID", "Name", "...", "Genre", "Stream Count"}
	rows := [][]interface{}{}
	for _, track := range tracks {
		rows = append(rows, []interface{}{track.ID, track.Name, "...", track.Genre, track.StreamCount})
	}

	tableoutput.PrintTable(table.StyleColoredDark, headers, rows)
}

// Многотабличный запрос
func (c *Controller) AlbumsWithTrackCounts(ctx context.Context) {
	var genre string
	fmt.Print("Enter genre: ")
	_, err := fmt.Scan(&genre)
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}

	albums, err := c.storage.AlbumsWithTrackCounts(ctx, genre)
	if err != nil {
		log.Printf("Error getting albums with track counts: %v", err)
		return
	}

	if len(albums) == 0 {
		fmt.Println("No albums found.")
		return
	}

	headers := []string{"Album ID", "Title", "Track Count"}
	rows := [][]interface{}{}
	for _, album := range albums {
		rows = append(rows, []interface{}{album.AlbumID, album.Title, album.TrackCount})
	}

	tableoutput.PrintTable(table.StyleColoredDark, headers, rows)
}

// Добавление пользователя
func (c *Controller) AddUser(ctx context.Context) {
	user := models.User{}

	fmt.Print("Enter user name: ")
	_, err := fmt.Scan(&user.Name)
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}

	user.BirthDate, err = c.userBirthDay()
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}

	user.Premium, user.PremiumExpiration, err = c.setPremium()
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}

	user.ID = uuid.New()
	user.RegistrationDate = time.Now()

	if err := c.storage.AddUser(ctx, &user); err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}

	fmt.Println("User added successfully.")
}

func (Controller) setPremium() (bool, time.Time, error) {
	var (
		premium    bool
		premiumExp time.Time
	)
	fmt.Print("Is the user premium? (true/false): ")
	_, err := fmt.Scan(&premium)
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return false, time.Time{}, err
	}

	if premium {
		fmt.Print("Enter premium expiration date (YYYY-MM-DD): ")
		var premiumExpiration string
		_, err = fmt.Scan(&premiumExpiration)
		if err != nil {
			slog.Error("[ERR]", "err", err)
			return false, time.Time{}, err
		}
		premiumExp, err = time.Parse("2006-01-02", premiumExpiration)
		if err != nil {
			log.Printf("Invalid premium expiration date format: %v", err)
			return false, time.Time{}, err
		}
	} else {
		premiumExp = time.Time{}
	}

	return premium, premiumExp, nil
}

func (Controller) userBirthDay() (models.JsonBirthDate, error) {
	fmt.Print("Enter birth date (YYYY-MM-DD): ")
	var birthDateStr string
	_, err := fmt.Scan(&birthDateStr)
	if err != nil {
		return models.JsonBirthDate(time.Now()), err
	}

	birthDate, err := time.Parse("2006-01-02", birthDateStr)
	if err != nil {
		return models.JsonBirthDate(time.Now()), err
	}

	return models.JsonBirthDate(birthDate), nil
}

// Обновление имени пользователя
func (c *Controller) UpdateUserName(ctx context.Context) {
	var userIDStr string
	var newName string

	fmt.Print("Enter user ID: ")
	_, err := fmt.Scan(&userIDStr)
	if err != nil {
		slog.Error("Failed to read user ID", "err", err)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		slog.Error("Invalid UUID format", "err", err)
		fmt.Println("Please enter a valid UUID.")
		return
	}

	fmt.Print("Enter new user name: ")
	_, err = fmt.Scan(&newName)
	if err != nil {
		slog.Error("Failed to read new name", "err", err)
		return
	}

	if err := c.storage.UpdateUserName(ctx, userID, newName); err != nil {
		log.Printf("Error updating user name: %v", err)
		return
	}

	fmt.Println("User name updated successfully.")
}

// Удаление пользователя
func (c *Controller) DeleteUser(ctx context.Context) {
	var userIDStr string

	fmt.Print("Enter user ID to delete: ")
	_, err := fmt.Scan(&userIDStr)
	if err != nil {
		slog.Error("Failed to read user ID", "err", err)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		slog.Error("Invalid UUID format", "err", err)
		fmt.Println("Please enter a valid UUID.")
		return
	}

	if err := c.storage.DeleteUser(ctx, userID); err != nil {
		log.Printf("Error deleting user: %v", err)
		return
	}

	fmt.Println("User deleted successfully.")
}

// Получение альбомов исполнителя
func (c *Controller) AlbumsByArtist(ctx context.Context) {
	var artistIDStr string

	fmt.Print("Enter artist ID: ")
	_, err := fmt.Scan(&artistIDStr)
	if err != nil {
		slog.Error("Failed to read artist ID", "err", err)
		return
	}

	artistID, err := uuid.Parse(artistIDStr)
	if err != nil {
		slog.Error("Invalid UUID format", "err", err)
		fmt.Println("Please enter a valid UUID.")
		return
	}

	albums, err := c.storage.AlbumsByArtist(ctx, artistID)
	if err != nil {
		log.Printf("Error getting albums by artist: %v", err)
		return
	}

	if len(albums) == 0 {
		fmt.Println("No albums found for the artist.")
		return
	}

	headers := []string{"ID", "Title", "Release Date", "Label", "Genre"}
	rows := [][]interface{}{}
	for _, album := range albums {
		rows = append(rows, []interface{}{album.ID, album.Title, album.ReleaseDate, album.Label, album.Genre})
	}

	tableoutput.PrintTable(table.StyleColoredDark, headers, rows)
}
