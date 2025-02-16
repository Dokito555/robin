package entity

import "time"

type Playlist struct {
	ID          int       `gorm:"column:id;type:primaryKey"`
	UserID      int       `gorm:"column:user_id"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (p *Playlist) TableName() string {
	return "playlists"
}
