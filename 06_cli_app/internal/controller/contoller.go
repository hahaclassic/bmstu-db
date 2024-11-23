package controller

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"log/slog"

	"github.com/google/uuid"
	"github.com/hahaclassic/databases/06_cli_app/internal/models"
	"github.com/hahaclassic/databases/06_cli_app/internal/storage"
	"github.com/jedib0t/go-pretty/table"
)

type Contoller struct {
	// c.storage c.storage.MusicService
	storage storage.MusicServiceStorage
}

func NewController(storage storage.MusicServiceStorage) *Contoller {
	return &Contoller{storage: storage}
}

// Start обрабатывает выбор операции и выполняет соответствующий метод
func (c *Contoller) Start(ctx context.Context) {
	methods := map[Operation]func(context.Context){
		CountArtists:             c.countArtists,
		TracksAdditionalInfo:     c.tracksAdditionalInfo,
		BestTracks:               c.bestTracks,
		TablesNames:              c.tablesNames,
		CountArtistTracks:        c.countArtistTracks,
		ArtistAlbums:             c.artistAlbums,
		FreePremiumForPensioners: c.freePremiumForPensioners,
		CurrentUser:              c.currentUser,
		CreateReviewsTable:       c.createReviewsTable,
		AddReview:                c.addReview,
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			menu()
			fmt.Print("Enter operation number (or 0 for exit): ")
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				fmt.Println("Invalid operation number. Try again.")
			}

			operation, err := strconv.Atoi(input)
			if err != nil || operation < 0 || operation > int(AddReview) {
				fmt.Println("Invalid data. Try again.")
				continue
			}

			if operation == 0 {
				fmt.Println("Exit")
				return
			}

			methods[Operation(operation)](ctx)
		}
	}
}

func menu() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Operation"})
	for i := CountArtists; i <= AddReview; i++ {
		t.AppendRow(table.Row{int(i), i.String()})
	}
	t.AppendRow(table.Row{0, "Exit"})
	t.Render()
}

func (app *Contoller) countArtists(ctx context.Context) {
	count, err := app.storage.CountArtists(ctx)
	if err != nil {
		log.Printf("Error counting artists: %v", err)
		return
	}
	fmt.Printf("Total number of artists: %d\n", count)
}

func (app *Contoller) tracksAdditionalInfo(ctx context.Context) {
	var limit int
	fmt.Print("Enter limit for tracks: ")
	_, err := fmt.Scan(&limit)
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}

	tracks, err := app.storage.TracksAdditionalInfo(ctx, limit)
	if err != nil {
		log.Printf("Error getting tracks info: %v", err)
		return
	}

	if len(tracks) == 0 {
		fmt.Println("No additional track info found.")
		return
	}

	headers := []string{"Track Name", "Artist", "Album", "ReleaseDate", "StreamCount"}
	rows := [][]interface{}{}
	for _, track := range tracks {
		rows = append(rows, []interface{}{track.Name, track.Artist, track.Album,
			track.ReleaseDate, track.StreamCount})
	}

	printTable(headers, rows)
}

func (app *Contoller) bestTracks(ctx context.Context) {
	var limit int
	fmt.Print("Enter limit for best tracks: ")
	_, err := fmt.Scan(&limit)
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}

	tracks, err := app.storage.BestTracks(ctx, limit)
	if err != nil {
		log.Printf("Error getting best tracks: %v", err)
		return
	}

	if len(tracks) == 0 {
		fmt.Println("No tracks found.")
		return
	}

	headers := []string{"Track Name", "Stream Count", "Rank"}
	rows := [][]interface{}{}
	for _, track := range tracks {
		rows = append(rows, []interface{}{track.Name, track.StreamCount, track.Rank})
	}

	printTable(headers, rows)
}

func (app *Contoller) tablesNames(ctx context.Context) {
	names, err := app.storage.TablesNames(ctx)
	if err != nil {
		log.Printf("Error getting table names: %v", err)
		return
	}

	if len(names) == 0 {
		fmt.Println("No tables found in the database.")
		return
	}

	headers := []string{"Table Names"}
	rows := [][]interface{}{}
	for _, name := range names {
		rows = append(rows, []interface{}{name})
	}

	printTable(headers, rows)
}

func (app *Contoller) countArtistTracks(ctx context.Context) {
	var artistID string
	fmt.Print("Enter artist ID (UUID): ")
	_, err := fmt.Scan(&artistID)
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}

	parsedID, err := uuid.Parse(artistID)
	if err != nil {
		log.Printf("Error parsing artist ID: %v", err)
		return
	}

	count, err := app.storage.CountArtistTracks(ctx, parsedID)
	if err != nil {
		log.Printf("Error getting artist tracks count: %v", err)
		return
	}
	fmt.Printf("Total tracks for artist %s: %d\n", artistID, count)
}

func (app *Contoller) artistAlbums(ctx context.Context) {
	var artistID string
	fmt.Print("Enter artist ID (UUID): ")
	_, err := fmt.Scan(&artistID)
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}

	parsedID, err := uuid.Parse(artistID)
	if err != nil {
		log.Printf("Error parsing artist ID: %v", err)
		return
	}

	albums, err := app.storage.ArtistAlbums(ctx, parsedID)
	if err != nil {
		log.Printf("Error getting albums for artist: %v", err)
		return
	}

	if len(albums) == 0 {
		fmt.Println("No albums found for this artist.")
		return
	}

	headers := []string{"Album ID", "Title", "Release Date"}
	rows := [][]interface{}{}
	for _, album := range albums {
		rows = append(rows, []interface{}{album.ID, album.Title, album.ReleaseDate})
	}

	printTable(headers, rows)
}

func (app *Contoller) freePremiumForPensioners(ctx context.Context) {
	err := app.storage.FreePremiumForPensioners(ctx)
	if err != nil {
		log.Printf("Error granting free premium: %v", err)
		return
	}
	fmt.Println("Free premium granted to pensioners.")
}

func (app *Contoller) currentUser(ctx context.Context) {
	user, err := app.storage.CurrentUser(ctx)
	if err != nil {
		log.Printf("Error getting current user: %v", err)
		return
	}
	fmt.Printf("Current user: %s\n", user)
}

func (app *Contoller) createReviewsTable(ctx context.Context) {
	err := app.storage.CreateReviewsTable(ctx)
	if errors.Is(err, storage.ErrTableAlreadyExists) {
		log.Print("Table already exists.")
		return
	} else if err != nil {
		log.Printf("Error creating reviews table: %v", err)
		return
	}

	fmt.Println("Reviews table created.")
}

func (app *Contoller) addReview(ctx context.Context) {
	id, err := uuid.NewRandom()
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}
	review := &models.Review{ID: id}

	var strUUID string
	fmt.Print("Enter review user ID (UUID): ")
	_, err = fmt.Scan(&strUUID)
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}
	parsedUUID, err := uuid.Parse(strUUID)
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}
	review.UserID = parsedUUID

	fmt.Print("Enter album ID (UUID): ")
	_, err = fmt.Scan(&strUUID)
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}
	parsedUUID, err = uuid.Parse(strUUID)
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}
	review.AlbumID = parsedUUID

	fmt.Print("Enter rating (1-10): ")
	_, err = fmt.Scan(&review.Rating)
	if err != nil {
		slog.Error("[ERR]", "err", err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter comment: ")
	if scanner.Scan() {
		review.Comment = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}

	err = app.storage.AddReview(ctx, review)
	if err != nil {
		log.Printf("Error adding review: %v", err)
		return
	}
	fmt.Println("Review added successfully!")
}

func printTable(headers []string, rows [][]interface{}) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	tableHeader := make(table.Row, 0, len(headers))
	for i := range headers {
		tableHeader = append(tableHeader, headers[i])
	}
	t.AppendHeader(tableHeader)

	for _, row := range rows {
		t.AppendRow(row)
	}
	t.SetStyle(table.StyleColoredDark)
	t.Render()
}
