package ws

const (
	TypeNotification = "SYSTEM_NOTIFICATION"
	TypeSystem       = "SYSTEM"
	TypeUser         = "USER"
)

type WebsocketEventType string

type WebsocketEvent struct {
	Type WebsocketEventType
	Data any
}
