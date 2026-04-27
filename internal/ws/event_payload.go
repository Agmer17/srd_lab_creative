package ws

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	TypeNotification = "SYSTEM_NOTIFICATION"
	TypeSystem       = "SYSTEM"
	TypeSystemError  = "SYSTEM_ERROR"
	TypeRoomJoin     = "USER_JOIN_ROOM"
	TypeChat         = "CHAT"
)

type WebsocketEventType string

type WebsocketEvent struct {
	Type WebsocketEventType `json:"type"`
	Data json.RawMessage    `json:"data"`
}

type ChatData struct {
	Id                   uuid.UUID       `json:"id"`
	ChatRoomId           string          `json:"chatroom_id"`
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

type JoinRoomData struct {
	RoomId string `json:"room_id"`
}

type SystemNotificationData struct {
	Message string `json:"message"`
}
