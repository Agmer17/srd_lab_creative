package model

import (
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/google/uuid"
)

type Chatroom struct {
	Id             uuid.UUID  `json:"id"`
	Type           string     `json:"type"`
	ProjectId      *uuid.UUID `json:"project_id,omitempty"`
	ParticipantKey *string    `json:"participant_key,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

func MapChatroomModel(gen sqlcgen.Chatroom) Chatroom {
	var projectId *uuid.UUID = nil

	if gen.ProjectID != uuid.Nil {
		projectId = &gen.ProjectID
	}

	return Chatroom{
		Id:             gen.ID,
		Type:           gen.Type,
		ProjectId:      projectId,
		ParticipantKey: gen.ParticipantKey,
		CreatedAt:      gen.CreatedAt,
	}
}
