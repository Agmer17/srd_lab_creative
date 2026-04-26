package chat

import (
	"context"
	"errors"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type ChatMediaRepository struct {
	db *sqlcgen.Queries
}

func NewChatMediaRepository(q *sqlcgen.Queries) *ChatMediaRepository {
	return &ChatMediaRepository{
		db: q,
	}
}

// =======================
// CREATE
// =======================
func (cmr *ChatMediaRepository) CreateChatMedia(ctx context.Context, m []model.ChatMedia) ([]model.ChatMedia, error) {
	chatIDs := []uuid.UUID{}
	fileNames := []string{}
	mediaTypes := []string{}

	for _, m := range m {
		chatIDs = append(chatIDs, m.ChatID)
		fileNames = append(fileNames, m.FileName)
		mediaTypes = append(mediaTypes, m.MediaType)
	}

	data, err := cmr.db.CreateChatMedia(ctx, sqlcgen.CreateChatMediaParams{
		ChatID:    chatIDs,
		Filename:  fileNames,
		MediaType: mediaTypes,
	})
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return []model.ChatMedia{}, errChatNotFound
			}
		}
		return []model.ChatMedia{}, err
	}

	return model.MapListMediaModel(data), nil
}

// =======================
// GET BY ROOM ID
// =======================
func (cmr *ChatMediaRepository) GetChatMediasByRoomID(
	ctx context.Context,
	roomID uuid.UUID,
) ([]model.ChatMedia, error) {

	rows, err := cmr.db.GetChatMediasByRoomID(ctx, roomID)
	if err != nil {
		return []model.ChatMedia{}, err
	}

	return model.MapListMediaModel(rows), nil
}

// =======================
// GET BY CHAT ID
// =======================
func (cmr *ChatMediaRepository) GetChatMediasByChatID(
	ctx context.Context,
	chatID uuid.UUID,
) ([]model.ChatMedia, error) {

	rows, err := cmr.db.GetChatMediasByChatID(ctx, chatID)
	if err != nil {
		return []model.ChatMedia{}, err
	}

	return model.MapListMediaModel(rows), nil
}
