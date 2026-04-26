package chat

import (
	"context"
	"errors"
	"fmt"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/Agmer17/srd_lab_creative/internal/storage"
	"github.com/google/uuid"
)

const chatMediaAtt = "chat_attachment"

type ChatService struct {
	repo          *ChatRepository
	mediaRepo     *ChatMediaRepository
	serverStorage *storage.FileStorage
}

func NewChatService(rp *ChatRepository, mr *ChatMediaRepository, sv *storage.FileStorage) *ChatService {
	return &ChatService{
		repo:          rp,
		mediaRepo:     mr,
		serverStorage: sv,
	}
}

// =======================
// CREATE CHAT
// =======================
func (cs *ChatService) CreateChat(ctx context.Context, senderId uuid.UUID, dto createChatDto) (model.Chat, *shared.ErrorResponse) {
	roomId, err := uuid.Parse(dto.RoomId)
	if err != nil {
		return model.Chat{}, shared.NewErrorResponse(400, "invalid room id please provide a valid uuid")
	}

	data, err := cs.repo.CreateChat(ctx, roomId, senderId, &dto.Text)
	if err != nil {
		if errors.Is(err, errChatroomNotFound) {
			return model.Chat{}, shared.NewErrorResponse(409, "chatroom not found! you can't send chat to invalid room id")
		}
		return model.Chat{}, shared.NewErrorResponse(500, "failed to create chat")
	}

	suc, tp, err := cs.serverStorage.SaveAllPrivateFile(ctx, dto.Attachment, chatMediaAtt)
	if err != nil {
		cs.repo.DeleteChat(ctx, data.ID)
		fmt.Println(err)
		return model.Chat{}, shared.NewErrorResponse(500, "something wrong while trying to save the file. Rolling back operation")
	}

	cmm := make([]model.ChatMedia, len(suc))
	for i, v := range suc {
		tempMed := tp[i]
		cmm[i] = model.ChatMedia{
			ChatID:    data.ID,
			FileName:  v,
			MediaType: tempMed,
		}
	}

	mediaData, err := cs.mediaRepo.CreateChatMedia(ctx, cmm)
	if err != nil {
		cs.repo.DeleteChat(ctx, data.ID)
		fmt.Println(err)
		return model.Chat{}, shared.NewErrorResponse(500, "something wrong while trying to save the file. Rolling back operation")
	}

	data.Medias = mediaData
	return data, nil
}

// =======================
// GET CHATS BY ROOM
// =======================
func (cs *ChatService) GetChatsByRoomID(ctx context.Context, roomID uuid.UUID) ([]model.Chat, *shared.ErrorResponse) {
	data, err := cs.repo.GetChatsByRoomID(ctx, roomID)
	if err != nil {
		return []model.Chat{}, shared.NewErrorResponse(500, "failed to get chats")
	}

	return data, nil
}

// =======================
// DELETE CHAT
// =======================
func (cs *ChatService) DeleteChat(ctx context.Context, id uuid.UUID, curr uuid.UUID) *shared.ErrorResponse {
	oldData, err := cs.repo.GetChatById(ctx, id)
	if err != nil {
		return shared.NewErrorResponse(500, "something wrng while trying delete chat data")
	}

	if oldData.Sender.ID != curr {
		return shared.NewErrorResponse(403, "you can't delete this message")
	}

	err = cs.repo.DeleteChat(ctx, id)
	if err != nil {
		if errors.Is(err, errChatNotFound) {
			return shared.NewErrorResponse(404, "chat not found")
		}
		return shared.NewErrorResponse(500, "failed to delete chat")
	}

	fileToDelete := make([]string, len(oldData.Medias))

	for i, v := range oldData.Medias {
		fileToDelete[i] = v.FileName
	}

	cs.serverStorage.DeleteAllPrivateFiles(fileToDelete, chatMediaAtt)
	return nil
}

// =======================
// GET LATEST CHAT PREVIEW
// =======================
func (cs *ChatService) GetLatestChatPreview(ctx context.Context, userID uuid.UUID) ([]LatestChatDto, *shared.ErrorResponse) {
	data, err := cs.repo.GetLatestChatPreview(ctx, userID)
	if err != nil {
		return []LatestChatDto{}, shared.NewErrorResponse(500, "failed to get latest chats")
	}

	return data, nil
}
