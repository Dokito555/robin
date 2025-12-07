package entity

import "time"

type Song struct {
	ID        int       `gorm:"primaryKey"`
	ArtistID  int       `gorm:"column:artist_id"`
	AlbumID   int       `gorm:"column:album_id"`
	Name      string    `gorm:"column:name"`
	Link      string    `gorm:"column:link"`
	Duration  int       `gorm:"column:duration"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Song) TableName() string {
	return "songs"
}

