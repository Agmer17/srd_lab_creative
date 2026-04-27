package chat

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var errChatNotFound = errors.New("chat not found")

type ChatRepository struct {
	db *sqlcgen.Queries
}

func NewChatRepository(q *sqlcgen.Queries) *ChatRepository {
	return &ChatRepository{
		db: q,
	}
}

func (cr *ChatRepository) CreateChat(ctx context.Context, roomID, senderID uuid.UUID, text *string) (model.Chat, error) {
	data, err := cr.db.CreateChat(ctx, sqlcgen.CreateChatParams{
		RoomID:   roomID,
		SenderID: senderID,
		Text:     text,
	})
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return model.Chat{}, errChatroomNotFound
			}
		}
		return model.Chat{}, err
	}

	userData, err := cr.db.GetUserById(ctx, senderID)
	if err != nil {
		return model.Chat{}, err
	}

	tempuserData := model.MapToUserModel(userData)

	tempModel := model.MapChatModel(data)
	tempModel.Sender = &tempuserData

	return tempModel, nil
}

func (cr *ChatRepository) GetChatsByRoomID(ctx context.Context, roomID uuid.UUID) ([]model.Chat, error) {
	rows, err := cr.db.GetChatsByRoomID(ctx, roomID)
	if err != nil {
		return []model.Chat{}, err
	}

	result, _ := MapGetChatsByRoomID(rows)

	return result, nil
}

func MapGetChatsByRoomID(rows []sqlcgen.GetChatsByRoomIDRow) ([]model.Chat, error) {
	var result []model.Chat

	for _, r := range rows {
		var sender *model.User
		if r.Sender != nil {
			var u model.User
			if err := json.Unmarshal(r.Sender, &u); err != nil {
				return nil, err
			}
			sender = &u
		}

		var medias []model.ChatMedia
		if r.Medias != nil {
			if err := json.Unmarshal(r.Medias, &medias); err != nil {
				return nil, err
			}
		}

		chat := model.Chat{
			ID:        r.ID,
			RoomID:    r.RoomID,
			Text:      r.Text,
			CreatedAt: r.CreatedAt,
			Sender:    sender,
			Medias:    medias,
		}

		// set SenderID kalau sender ada
		if sender != nil {
			chat.SenderID = &sender.ID
		}

		result = append(result, chat)
	}

	return result, nil
}

func (cr *ChatRepository) DeleteChat(ctx context.Context, id uuid.UUID) error {
	aff, err := cr.db.DeleteChat(ctx, id)

	if err != nil {
		return err
	}

	if aff == 0 {
		return errChatNotFound
	}
	return nil
}

func (cr *ChatRepository) GetLatestChatPreview(ctx context.Context, userID uuid.UUID) ([]LatestChatDto, error) {
	rows, err := cr.db.GetLatestChatPreview(ctx, userID)
	if err != nil {
		return []LatestChatDto{}, err
	}

	result := make([]LatestChatDto, len(rows))

	for i, r := range rows {
		var avatar *string = nil

		if r.Avatar != "" {
			avatar = &r.Avatar
		}

		dto := LatestChatDto{
			ChatroomID:    r.ChatroomID.String(),
			Type:          string(r.Type),
			Name:          r.Name,
			Avatar:        avatar,
			LastMessage:   r.LastMessage,
			LastMessageAt: &r.LastMessageAt,
		}
		result[i] = dto
	}

	return result, nil
}

func (cr *ChatRepository) GetChatById(ctx context.Context, id uuid.UUID) (model.Chat, error) {

	data, err := cr.db.GetChatID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Chat{}, errChatNotFound
		}
		return model.Chat{}, err
	}

	var userData model.User
	err = json.Unmarshal(data.Sender, &userData)
	if err != nil {
		return model.Chat{}, err
	}

	var chatMedia []model.ChatMedia
	err = json.Unmarshal(data.Medias, &chatMedia)
	if err != nil {
		return model.Chat{}, err
	}

	chtModel := model.Chat{
		ID:       data.ID,
		RoomID:   data.RoomID,
		Sender:   &userData,
		SenderID: &userData.ID,
		Medias:   chatMedia,
	}

	return chtModel, nil
}
