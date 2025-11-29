package entity

import "time"

type Album struct {
	ID        int       `gorm:"column:id;type:primaryKey"`
	Name      string    `gorm:"column:name"`
	ArtistID  int       `gorm:"column:artist_id"`
	Type      string    `gorm:"column:type"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (u *Album) TableName() string {
	return "albums"
}