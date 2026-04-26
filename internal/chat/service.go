package chat

type MessagingService struct {
	chatService *ChatService
	roomService *ChatroomService
}

func NewMessagingService(cht *ChatService, room *ChatroomService) *MessagingService {

	return &MessagingService{
		chatService: cht,
		roomService: room,
	}
}
