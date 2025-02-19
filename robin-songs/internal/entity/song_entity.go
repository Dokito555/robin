package entity

import "time"

type Song struct {
	ID        int       `gorm:"column:id;type:primaryKey"`
	ArtistID  int       `gorm:"column:artist_id"`
	AlbumID   int       `gorm:"column:album_id"`
	Name      string    `gorm:"column:name"`
	Link      string    `gorm:"column:link"`
	Duration  int       `gorm:"column:duration"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:update_at"`
}

func (*Song) TableName() string {
	return "songs"
}
