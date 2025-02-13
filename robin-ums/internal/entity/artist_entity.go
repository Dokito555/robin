package entity

import "time"

type Artist struct {
	ID           int    `gorm:"column:id;type:primaryKey"`
	Username     string `gorm:"column:username;type:varchar(255)"`
	Email        string `gorm:"column:email;unique"`
	Genre        string `gorm:"column:genre"`
	Bio          string `gorm:"column:bio;type:text"`
	Role         string `gorm:"column:role"`
	Token        string `gorm:"column:token"`
	RefreshToken string `gorm:"column:refresh_token"`
	// SubscriptionType string `gorm:"column:subscription_type"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (u *Artist) TableName() string {
	return "artists"
}