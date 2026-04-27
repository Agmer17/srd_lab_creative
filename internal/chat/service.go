package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/Agmer17/srd_lab_creative/internal/ws"
	"github.com/Agmer17/srd_lab_creative/pkg"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type MessagingService struct {
	chatService *ChatService
	roomService *ChatroomService
	rdb         *redis.Client
	hub         *ws.WebsocketHub
}

func NewMessagingService(cht *ChatService, room *ChatroomService, hub *ws.WebsocketHub, red *redis.Client) *MessagingService {

	// todo : init cache
	return &MessagingService{
		chatService: cht,
		roomService: room,
		rdb:         red,
		hub:         hub,
	}
}

func (ms *MessagingService) GetLatestChat(ctx context.Context, curr uuid.UUID) ([]LatestChatDto, *shared.ErrorResponse) {
	chatData, err := ms.chatService.GetLatestChatPreview(ctx, curr)
	if err != nil {
		fmt.Println("error : ", err)
		return []LatestChatDto{}, shared.NewErrorResponse(500, "somehting wrong while trying to get chat data")
	}

	return chatData, nil
}

func (ms *MessagingService) CreateProjectMessage(ctx context.Context, curr, projectId uuid.UUID, dto createChatDto) (ChatDataDto, *shared.ErrorResponse) {
	vld, err := ms.validateChatroomMember(ctx, curr, projectId)
	if err != nil {
		return ChatDataDto{}, shared.NewErrorResponse(500, "something wrong while trying to validate message : "+err.Error())
	}

	if !vld {
		return ChatDataDto{}, shared.NewErrorResponse(403, "you are forbidden to sending message to this room")
	}

	chat, svErr := ms.chatService.CreateChat(ctx, curr, dto)
	if svErr != nil {
		return ChatDataDto{}, svErr
	}

	chatPayloadData := ms.ChatmodelToDto(ctx, chat)

	dataPayload, _ := json.Marshal(chatPayloadData)

	wsPayload := ws.WebsocketEvent{
		Type: ws.TypeChat,
		Data: dataPayload,
	}

	payload, _ := json.Marshal(wsPayload)

	ms.hub.SendPayloadTo("room:"+chat.RoomID.String(), payload)

	return ChatDataDto{}, nil
}

func (ms *MessagingService) validateChatroomMember(ctx context.Context, curr uuid.UUID, projectId uuid.UUID) (bool, error) {

	setKey := "member:" + projectId.String()
	isMember, err := ms.rdb.SIsMember(ctx, setKey, curr.String()).Result()
	if err != nil {
		return false, err
	}

	return isMember, nil
}

func (ms *MessagingService) ChatmodelToDto(ctx context.Context, chat model.Chat) ChatDataDto {
	// Convert chat media nya jadi signed url
	return ChatDataDto{
		Id:                   chat.ID,
		ChatRoomId:           chat.RoomID,
		SenderId:             chat.Sender.ID,
		SenderFullName:       chat.Sender.FullName,
		SenderProfilePiCture: *chat.Sender.ProfilePicture,
		Text:                 *chat.Text,
		Media:                ms.ChatmediaToModel(ctx, chat.Medias),
		CreatedAt:            chat.CreatedAt,
	}
}

func (ms *MessagingService) ChatmediaToModel(ctx context.Context, med []model.ChatMedia) []ChatMediaType {
	pipe := ms.rdb.Pipeline()
	media := make([]ChatMediaType, len(med))
	for i, v := range med {

		randToken, _ := pkg.GenerateSecureString(24)
		hashkey := "media_access:private:" + randToken
		media[i] = ChatMediaType{
			Type: v.MediaType,
			Url:  randToken,
		}

		pipe.HSet(ctx, hashkey, map[string]interface{}{
			"type":     v.MediaType,
			"filename": v.FileName,
		})

		pipe.Expire(ctx, hashkey, time.Hour)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		fmt.Println("ERROR TRYING TO SET ACCESS KEY TO MEDIA : ", err)
	}

	return media
}
