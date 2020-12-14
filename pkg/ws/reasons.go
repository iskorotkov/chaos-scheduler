package ws

type CloseReason string

var (
	ErrorOccurred    = CloseReason("error occurred while waiting for connection closing")
	DeadlineExceeded = CloseReason("websocket connection deadline exceeded")
	ClosedOnServer   = CloseReason("websocket connection was closed on the server")
	ClosedOnClient   = CloseReason("websocket connection was closed on the client")
)
