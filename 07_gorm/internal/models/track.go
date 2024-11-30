package models

import "github.com/google/uuid"

type Track struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	Name         string
	OrderInAlbum int
	AlbumID      uuid.UUID
	Explicit     bool
	Duration     int
	Genre        string
	StreamCount  int64
}
