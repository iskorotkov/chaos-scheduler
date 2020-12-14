package ws

type CloseReason string

var (
	ReasonErrorOccurred    = CloseReason("error occurred while waiting for connection closing")
	ReasonDeadlineExceeded = CloseReason("websocket connection deadline exceeded")
	ReasonClosedOnServer   = CloseReason("websocket connection was closed on the server")
	ReasonClosedOnClient   = CloseReason("websocket connection was closed on the client")
	ReasonEOF              = CloseReason("websocket was closed due to EOF caused by disconnected client")
)
