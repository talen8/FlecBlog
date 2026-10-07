package model

import (
	"time"
)

// RssArticle RSS文章模型
type RssArticle struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	FriendID    uint       `gorm:"not null;index" json:"friend_id"`
	Friend      *Friend    `gorm:"foreignKey:FriendID" json:"friend"`
	Title       string     `gorm:"size:500;not null" json:"title"`
	Link        string     `gorm:"size:1000;not null;unique" json:"link"`
	PublishedAt *time.Time `json:"published_at"`
	IsRead      bool       `gorm:"default:false" json:"is_read"`
	BlockCircle bool       `gorm:"default:false" json:"block_circle"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
