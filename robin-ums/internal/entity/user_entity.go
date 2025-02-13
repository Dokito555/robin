package entity

import "time"

type User struct {
	ID           int    `gorm:"column:id;type:primaryKey"`
	UserName     string `gorm:"column:username;type:varchar(255)"`
	Email        string `gorm:"column:email;unique"`
	Password     string `gorm:"column:password"`
	Role         string `gorm:"column:role"`
	Token        string `gorm:"column:token"`
	RefreshToken string `gorm:"column:refresh_token"`
	// SubscriptionType string `gorm:"column:subscription_type"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (u *User) TableName() string {
	return "users"
}
