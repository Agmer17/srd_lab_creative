package chat

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type ChatDataDto struct {
	Id                   uuid.UUID       `json:"id"`
	ChatRoomId           uuid.UUID       `json:"chatroom_id"`
	SenderId             uuid.UUID       `json:"sender_id"`
	SenderFullName       string          `json:"sender_full_name"`
	SenderProfilePiCture string          `json:"sender_profile_picture"`
	Text                 string          `json:"text"`
	Media                []ChatMediaType `json:"chat_media,omitempty"`
	CreatedAt            time.Time       `json:"created_at"`
}

type ChatMediaType struct {
	Type string `json:"media_type"`
	Url  string `json:"media_access_url"`
}

type LatestChatDto struct {
	ChatroomID    string     `json:"chatroom_id"`
	Type          string     `json:"type"`
	Name          string     `json:"name"`
	Avatar        *string    `json:"avatar"`
	LastMessage   string     `json:"last_message"`
	LastMessageAt *time.Time `json:"last_message_at"`
}

type createChatDto struct {
	Text       string                  `form:"text" binding:"required,min=1"`
	RoomId     string                  `form:"room_id" binding:"required,uuid"`
	Attachment []*multipart.FileHeader `form:"attachment"`
}
