package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HelpArticle 幫助中心的具體問答文章
type HelpArticle struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null" json:"category_id"` // 所屬分類 ID
	Question   string    `gorm:"type:varchar(255);not null" json:"question"` // 問題標題
	Answer     string    `gorm:"type:text;not null" json:"answer"`         // 答案內容 (支援 Markdown 或 HTML)
	Order      int       `gorm:"default:0" json:"order"`                 // 在分類內的排序值
	IsActive   bool      `gorm:"not null;default:true" json:"is_active"` // 是否啟用
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (a *HelpArticle) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
