package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	UserID         uuid.UUID `json:"user_id"`
	TelegramUserID int64     `json:"telegram_user_id"`
	PrimaryChatID  int64     `json:"primary_chat_id"`
	Timezone       string    `json:"timezone"`
	Locale         string    `json:"locale"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func NewUserProfile(telegramUserID, chatID int64) UserProfile {
	now := time.Now()
	return UserProfile{
		UserID:         uuid.New(),
		TelegramUserID: telegramUserID,
		PrimaryChatID:  chatID,
		Timezone:       "America/Sao_Paulo",
		Locale:         "pt-BR",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (u *UserProfile) Location() *time.Location {
	loc, err := time.LoadLocation(u.Timezone)
	if err != nil {
		loc = time.UTC
	}
	return loc
}
