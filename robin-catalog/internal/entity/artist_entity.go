package entity

import "time"

type Artist struct {
	ID           int       `gorm:"primaryKey"`
	UserName     string    `gorm:"column:username;type:varchar(255)"`
	Email        string    `gorm:"column:email;unique"`
	Password     string    `gorm:"column:password"`
	Genre        string    `gorm:"column:genre"`
	Bio          string    `gorm:"column:bio;type:text"`
	Role         string    `gorm:"column:role"`
	Token        string    `gorm:"column:token"`
	RefreshToken string    `gorm:"column:refresh_token"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (Artist) TableName() string {
	return "artists"
}

