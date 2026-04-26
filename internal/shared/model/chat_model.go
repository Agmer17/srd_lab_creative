package model

import (
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID  `json:"id"`
	RoomID    uuid.UUID  `json:"room_id"`
	SenderID  *uuid.UUID `json:"sender_id,omitempty"`
	Text      *string    `json:"text,omitempty"`
	CreatedAt time.Time  `json:"created_at"`

	Sender *User       `json:"sender,omitempty"`
	Medias []ChatMedia `json:"attachment,omitempty"`
}

type ChatMedia struct {
	ID        uuid.UUID `json:"id"`
	ChatID    uuid.UUID `json:"chat_id"`
	FileName  string    `json:"file_name"`
	MediaType string    `json:"media_type,omitempty"`
	Size      *int64    `json:"size,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func MapChatModel(c sqlcgen.Chat) Chat {
	var senderID *uuid.UUID
	if c.SenderID != uuid.Nil {
		senderID = &c.SenderID
	}

	return Chat{
		ID:        c.ID,
		RoomID:    c.RoomID,
		SenderID:  senderID,
		Text:      c.Text,
		CreatedAt: c.CreatedAt,
	}
}

func MapChatMediaModel(data sqlcgen.ChatMedia) ChatMedia {
	return ChatMedia{
		ID:        data.ID,
		ChatID:    data.ChatID,
		FileName:  data.FileName,
		MediaType: data.MediaType,
		Size:      data.Size,
		CreatedAt: data.CreatedAt,
	}
}

func MapListMediaModel(gen []sqlcgen.ChatMedia) []ChatMedia {
	var result []ChatMedia = make([]ChatMedia, len(gen))

	for i, v := range gen {
		result[i] = MapChatMediaModel(v)
	}

	return result
}
