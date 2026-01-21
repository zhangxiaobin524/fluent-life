package models

import (
	"time"

	"github.com/google/uuid"
)

type RandomMatchRecord struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null;index:idx_random_match_user" json:"user_id"`
	MatchedUserID *uuid.UUID `gorm:"type:uuid;index:idx_random_match_matched_user" json:"matched_user_id,omitempty"`
	Status        string     `gorm:"type:varchar(20);not null;index:idx_random_match_status" json:"status"` // pending/matched/cancelled/timeout
	WaitSeconds   *int       `json:"wait_seconds,omitempty"`
	MatchedAt     *time.Time `json:"matched_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	User         User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	MatchedUser  *User `gorm:"foreignKey:MatchedUserID" json:"matched_user,omitempty"`
}
