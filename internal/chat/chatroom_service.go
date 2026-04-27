package chat

import (
	"context"
	"errors"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
)

type ChatroomService struct {
	repo *ChatroomRepository
}

func NewChatroomService(rp *ChatroomRepository) *ChatroomService {
	return &ChatroomService{
		repo: rp,
	}
}

func (css *ChatroomService) GetChatroomByID(ctx context.Context, id uuid.UUID) (model.Chatroom, *shared.ErrorResponse) {
	data, err := css.repo.GetChatroomByID(ctx, id)
	if err != nil {
		if errors.Is(err, errChatroomNotFound) {
			return model.Chatroom{}, shared.NewErrorResponse(404, "chatroom not found")
		}
		return model.Chatroom{}, shared.NewErrorResponse(500, "something wrong while trying to get chatroom")
	}

	return data, nil
}

func (css *ChatroomService) GetChatroomByProjectID(ctx context.Context, projectId uuid.UUID) (model.Chatroom, *shared.ErrorResponse) {
	data, err := css.repo.GetChatroomByProjectID(ctx, projectId)
	if err != nil {
		if errors.Is(err, errChatroomNotFound) {
			return model.Chatroom{}, shared.NewErrorResponse(404, "chatroom for this project not found")
		}
		return model.Chatroom{}, shared.NewErrorResponse(500, "failed to get project chatroom")
	}

	return data, nil
}

func (css *ChatroomService) GetChatroomByParticipantKey(ctx context.Context, key string) (model.Chatroom, *shared.ErrorResponse) {
	data, err := css.repo.GetChatroomByParticipantKey(ctx, key)
	if err != nil {
		if errors.Is(err, errChatroomNotFound) {
			return model.Chatroom{}, shared.NewErrorResponse(404, "chatroom not found for this participants")
		}
		return model.Chatroom{}, shared.NewErrorResponse(500, "failed to get personal chatroom")
	}

	return data, nil
}

func (css *ChatroomService) CreateProjectChatroom(ctx context.Context, projectId uuid.UUID) (model.Chatroom, *shared.ErrorResponse) {
	data, err := css.repo.CreateProjectChatroom(ctx, projectId)
	if err != nil {
		if errors.Is(err, errProjectIdNotValid) {
			return model.Chatroom{}, shared.NewErrorResponse(409, "invalid project id")
		}
		return model.Chatroom{}, shared.NewErrorResponse(500, "failed to create project chatroom")
	}

	return data, nil
}

func (css *ChatroomService) CreatePersonalChatroom(ctx context.Context, key string) (model.Chatroom, *shared.ErrorResponse) {
	if key == "" {
		return model.Chatroom{}, shared.NewErrorResponse(400, "participant key is required")
	}

	data, err := css.repo.CreatePersonalChatroom(ctx, key)
	if err != nil {
		return model.Chatroom{}, shared.NewErrorResponse(500, "failed to create personal chatroom")
	}

	return data, nil
}

func (css *ChatroomService) GetProjectChatroomMember(ctx context.Context, projectId uuid.UUID) ([]model.ProjectMember, *shared.ErrorResponse) {

	data, err := css.repo.GetProjectChatroomMember(ctx, projectId)
	if err != nil {
		return []model.ProjectMember{}, shared.NewErrorResponse(500, "something wrong while trying to get member data")
	}

	return data, nil
}
