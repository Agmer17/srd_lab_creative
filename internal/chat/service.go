package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/Agmer17/srd_lab_creative/internal/storage"
	"github.com/Agmer17/srd_lab_creative/internal/ws"
	"github.com/Agmer17/srd_lab_creative/pkg"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type MessagingService struct {
	chatService   *ChatService
	roomService   *ChatroomService
	rdb           *redis.Client
	hub           *ws.WebsocketHub
	serverStorage *storage.FileStorage
}

func NewMessagingService(cht *ChatService, room *ChatroomService, hub *ws.WebsocketHub, red *redis.Client, strg *storage.FileStorage) *MessagingService {

	// todo : init cache
	return &MessagingService{
		chatService:   cht,
		roomService:   room,
		rdb:           red,
		hub:           hub,
		serverStorage: strg,
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

	chatPayloadData := ms.ChatmodelToDto(ctx, chat, curr)

	dataPayload, _ := json.Marshal(chatPayloadData)

	wsPayload := ws.WebsocketEvent{
		Type: ws.TypeChat,
		Data: dataPayload,
	}

	payload, _ := json.Marshal(wsPayload)

	ms.hub.SendPayloadTo("room:"+chat.RoomID.String(), payload)

	return chatPayloadData, nil
}

func (ms *MessagingService) GetAllMessageFromProject(ctx context.Context, curr, projectId, roomId uuid.UUID) ([]ChatDataDto, *shared.ErrorResponse) {
	vld, err := ms.validateChatroomMember(ctx, curr, projectId)
	if err != nil {
		return []ChatDataDto{}, shared.NewErrorResponse(500, "something wrong while trying to get message data")
	}

	if !vld {
		return []ChatDataDto{}, shared.NewErrorResponse(403, "you are forbidden to sending message to this room")
	}

	data, getErr := ms.chatService.GetChatsByRoomID(ctx, roomId)
	if getErr != nil {
		return []ChatDataDto{}, shared.NewErrorResponse(500, "something wrong while trying to get message data")
	}

	var listDto []ChatDataDto = make([]ChatDataDto, len(data))
	for i, v := range data {
		listDto[i] = ms.ChatmodelToDto(ctx, v, curr)
	}

	return listDto, nil
}

func (ms *MessagingService) validateChatroomMember(ctx context.Context, curr uuid.UUID, projectId uuid.UUID) (bool, error) {

	setKey := "member:" + projectId.String()
	isMember, err := ms.rdb.SIsMember(ctx, setKey, curr.String()).Result()
	if err != nil {
		return false, err
	}

	return isMember, nil
}

func (ms *MessagingService) ChatmodelToDto(ctx context.Context, chat model.Chat, curr uuid.UUID) ChatDataDto {
	// Convert chat media nya jadi signed url
	return ChatDataDto{
		Id:                   chat.ID,
		ChatRoomId:           chat.RoomID,
		SenderId:             chat.Sender.ID,
		SenderFullName:       chat.Sender.FullName,
		SenderProfilePiCture: *chat.Sender.ProfilePicture,
		Text:                 *chat.Text,
		Media:                ms.ChatmediaToModel(ctx, chat.Medias, curr),
		CreatedAt:            chat.CreatedAt,
	}
}

func (ms *MessagingService) ChatmediaToModel(ctx context.Context, med []model.ChatMedia, curr uuid.UUID) []ChatMediaType {
	if len(med) == 0 {
		return nil
	}

	pipe := ms.rdb.Pipeline()
	media := make([]ChatMediaType, len(med))
	for i, v := range med {

		randToken, _ := pkg.GenerateSecureString(24)
		hashkey := "media_access:private:" + randToken
		media[i] = ChatMediaType{
			Type: v.MediaType,
			Url:  "http://localhost/api/chat/private-media/" + randToken,
		}

		pipe.HSet(ctx, hashkey, map[string]interface{}{
			"type":         v.MediaType,
			"filename":     v.FileName,
			"allowed_user": curr.String(),
		})

		pipe.Expire(ctx, hashkey, time.Hour)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		fmt.Println("ERROR TRYING TO SET ACCESS KEY TO MEDIA : ", err)
	}

	return media
}

func (ms *MessagingService) GetMediaAccessFromToken(ctx context.Context, token string, curr uuid.UUID) (string, *shared.ErrorResponse) {

	hashkey := "media_access:private:" + token
	data, err := ms.rdb.HGetAll(ctx, hashkey).Result()
	if err != nil {
		return "", shared.NewErrorResponse(500, "something wrong while trying to get file")
	}

	allowedUser := data["allowed_user"]
	if allowedUser != curr.String() {
		return "", shared.NewErrorResponse(403, "access to this media is forbidden")
	}

	fileName := path.Join(ms.serverStorage.PrivatePath, chatMediaAtt, data["filename"])
	return fileName, nil
}

func (ms *MessagingService) CreatePersonalChat(
	ctx context.Context,
	target, curr uuid.UUID,
	dto createPersonalChatDto,
) (ChatDataDto, *shared.ErrorResponse) {

	participantKey := CompareUUID(curr.String(), target.String())

	// 1. Cari room
	room, err := ms.roomService.repo.GetChatroomByParticipantKey(ctx, participantKey)

	if err != nil {
		// kalau bukan "not found" → error beneran
		if !errors.Is(err, errChatroomNotFound) {
			fmt.Println(err)
			return ChatDataDto{}, shared.NewErrorResponse(500, "failed to get chatroom")
		}

		// 2. Kalau belum ada → create
		room, err = ms.roomService.repo.CreatePersonalChatroom(ctx, participantKey)
		if err != nil {
			return ChatDataDto{}, shared.NewErrorResponse(500, "failed to create chatroom")
		}

		roomIds := []uuid.UUID{
			room.Id,
			room.Id,
		}
		userIds := []uuid.UUID{
			curr,
			target,
		}
		persErr := ms.roomService.AddPersonalMember(ctx, roomIds, userIds)
		if persErr != nil {
			fmt.Println(persErr)
			return ChatDataDto{}, shared.NewErrorResponse(500, "failed to create member")
		}
	}

	// 3. Pasti udah punya room di sini
	chat, inserr := ms.chatService.CreateChat(ctx, curr, createChatDto{
		Text:       dto.Text,
		RoomId:     room.Id.String(),
		Attachment: dto.Attachment,
	})
	if inserr != nil {
		return ChatDataDto{}, shared.NewErrorResponse(500, "failed to create chat")
	}

	dtoChat := ms.ChatmodelToDto(ctx, chat, curr)

	dataPayload, _ := json.Marshal(dtoChat)

	wsPayload := ws.WebsocketEvent{
		Type: ws.TypeChat,
		Data: dataPayload,
	}

	payload, _ := json.Marshal(wsPayload)

	ms.hub.SendPayloadTo(target.String(), payload)
	return dtoChat, nil
}

func (ms *MessagingService) DeleteChat(ctx context.Context, id uuid.UUID, curr uuid.UUID) *shared.ErrorResponse {

	err := ms.chatService.DeleteChat(ctx, id, curr)
	if err != nil {
		return err
	}

	return nil
}

func (ms *MessagingService) GetPersonalChatData(ctx context.Context, roomId, curr uuid.UUID) ([]ChatDataDto, *shared.ErrorResponse) {

	data, err := ms.chatService.GetChatsByRoomID(ctx, roomId)
	if err != nil {
		return []ChatDataDto{}, err
	}

	var listDto []ChatDataDto = make([]ChatDataDto, len(data))
	for i, v := range data {
		listDto[i] = ms.ChatmodelToDto(ctx, v, curr)
	}

	return listDto, nil
}

func CompareUUID(u1, u2 string) string {
	if u1 == u2 {
		return u1 + ":" + u2
	}

	minLen := len(u1)
	if len(u2) < minLen {
		minLen = len(u2)
	}

	for i := 0; i < minLen; i++ {
		if u1[i] < u2[i] {
			return u1 + ":" + u2
		}
		if u1[i] > u2[i] {
			return u2 + ":" + u1
		}
	}

	if len(u1) < len(u2) {
		return u1 + ":" + u2
	}

	return u2 + ":" + u1
}
