package ws

import "errors"

var (
	ConnectionError       = errors.New("couldn't upgrade websocket connection")
	DeadlineSettingError  = errors.New("couldn't set deadline")
	DeadlineExceededError = errors.New("connection deadline exceeded")
	WaitError             = errors.New("error occurred while waiting for client closing connection")
	ReadError             = errors.New("couldn't read next message")
	EOF                   = errors.New("read all messages")
	DecodeError           = errors.New("couldn't decode json message")
	EncodeError           = errors.New("couldn't encode json message")
	FlushError            = errors.New("couldn't flush encoded message")
	CloseError            = errors.New("couldn't close websocket connection")
)
