package model

import (
	"time"
)

// FriendType 友链类型模型
type FriendType struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"size:50;not null" json:"name"`   // 类型名称
	Sort      int       `gorm:"default:0" json:"sort"`          // 排序权重
	IsVisible bool      `gorm:"default:true" json:"is_visible"` // 是否展示
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Friend 友链模型
type Friend struct {
	ID          uint        `gorm:"primarykey" json:"id"`
	Name        string      `gorm:"size:50;not null" json:"name"`            // 网站名称
	URL         string      `gorm:"size:255;not null" json:"url"`            // 网站链接
	Description string      `gorm:"type:text" json:"description"`            // 网站描述
	Avatar      string      `gorm:"size:255" json:"avatar"`                  // 网站头像/logo
	Screenshot  string      `gorm:"size:255" json:"screenshot"`              // 网站截图
	Sort        int         `gorm:"default:5" json:"sort"`                   // 排序（数字越大越靠前）
	IsInvalid   bool        `gorm:"default:false" json:"is_invalid"`         // 是否失效
	IsPending   bool        `gorm:"default:false" json:"is_pending"`         // 是否为待审核申请
	TypeID      *uint       `gorm:"column:type" json:"type_id"`              // 类型 ID（外键）
	Type        *FriendType `gorm:"foreignKey:TypeID" json:"type,omitempty"` // 关联的类型
	RSSUrl      string      `gorm:"size:500;default:''" json:"rss_url"`      // RSS订阅地址
	RSSLatime   *time.Time  `json:"rss_latime,omitempty"`                    // RSS订阅最后更新时间
	Accessible  int         `gorm:"default:0" json:"accessible"`             // 可访问性状态: 0=正常, -1=忽略检查, >0=连续失败次数
	BlockCircle bool        `gorm:"default:false" json:"block_circle"`       // 友圈屏蔽
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
