package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HelpCategory 幫助中心文章的分類
type HelpCategory struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"` // 分類名稱，例如 "常見問題"
	Order     int       `gorm:"default:0" json:"order"`                 // 排序值，數字越小越靠前
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Articles []HelpArticle `gorm:"foreignKey:CategoryID" json:"articles,omitempty"` // 關聯的文章列表
}

func (c *HelpCategory) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}
