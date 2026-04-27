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

var errProjectIdNotValid = errors.New("project id not found! can't create chatroom project")
var errChatroomNotFound = errors.New("chatroom with this id not found")

type ChatroomRepository struct {
	db *sqlcgen.Queries
}

func NewChatroomRepository(q *sqlcgen.Queries) *ChatroomRepository {

	return &ChatroomRepository{
		db: q,
	}
}

func (crr *ChatroomRepository) CreateProjectChatroom(ctx context.Context, projectId uuid.UUID) (model.Chatroom, error) {
	data, err := crr.db.CreateProjectChatroom(ctx, projectId)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return model.Chatroom{}, errProjectIdNotValid
			}
		}
		return model.Chatroom{}, err
	}

	return model.MapChatroomModel(data), nil
}

func (crr *ChatroomRepository) CreatePersonalChatroom(ctx context.Context, key string) (model.Chatroom, error) {

	data, err := crr.db.CreatePersonalChatroom(ctx, &key)
	if err != nil {
		return model.Chatroom{}, err
	}

	return model.MapChatroomModel(data), nil
}

func (crr *ChatroomRepository) GetChatroomByID(ctx context.Context, id uuid.UUID) (model.Chatroom, error) {
	data, err := crr.db.GetChatroomByID(ctx, id)
	if err != nil {
		return model.Chatroom{}, err
	}

	return model.MapChatroomModel(data), nil
}

func (crr *ChatroomRepository) GetChatroomByProjectID(ctx context.Context, projectId uuid.UUID) (model.Chatroom, error) {
	data, err := crr.db.GetChatroomByProjectID(ctx, projectId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Chatroom{}, errChatroomNotFound
		}
		return model.Chatroom{}, err
	}

	return model.MapChatroomModel(data), nil
}

func (crr *ChatroomRepository) GetChatroomByParticipantKey(ctx context.Context, key string) (model.Chatroom, error) {
	data, err := crr.db.GetChatroomByParticipantKey(ctx, &key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Chatroom{}, errChatroomNotFound
		}
		return model.Chatroom{}, err
	}

	return model.MapChatroomModel(data), nil
}

func (crr *ChatroomRepository) GetProjectChatroomMember(ctx context.Context, projectId uuid.UUID) ([]model.ProjectMember, error) {

	data, err := crr.db.GetActiveProjectMembers(ctx, projectId)
	if err != nil {
		return []model.ProjectMember{}, err
	}

	var listData []model.ProjectMember = make([]model.ProjectMember, len(data))
	for i, v := range data {

		var userData model.User
		err := json.Unmarshal(v.User, &userData)
		if err != nil {
			return []model.ProjectMember{}, err
		}

		var roleData model.ProjectRole
		umsErr := json.Unmarshal(v.Role, &roleData)
		if umsErr != nil {
			return []model.ProjectMember{}, err
		}

		listData[i] = model.ProjectMember{
			ID:        v.ID,
			ProjectID: v.ProjectID,
			IsOwner:   v.IsOwner,
			JoinedAt:  v.JoinedAt,
			User:      userData,
			Role:      roleData,
		}

	}

	return listData, nil
}

func (crr *ChatroomRepository) AddPersonalParticipants(ctx context.Context, chatroomId []uuid.UUID, userId []uuid.UUID) error {

	_, err := crr.db.AddPersonalChatroomParticipant(ctx, sqlcgen.AddPersonalChatroomParticipantParams{
		ChatroomID: chatroomId,
		UserID:     userId,
	})

	if err != nil {
		return err
	}

	return nil

}
