package ws

import "errors"

var (
	ConnectionError = errors.New("couldn't upgrade websocket connection")
	ReadError       = errors.New("couldn't read next message")
	EOF             = errors.New("read all messages")
	DecodeError     = errors.New("couldn't decode json message")
	EncodeError     = errors.New("couldn't encode json message")
	FlushError      = errors.New("couldn't flush encoded message")
	CloseError      = errors.New("couldn't close websocket connection")
)
