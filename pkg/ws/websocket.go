package ws

import (
	"encoding/json"
	"errors"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"time"
)

type Websocket struct {
	conn    net.Conn
	timeout time.Duration
	logger  *zap.SugaredLogger
}

func NewWebsocket(w http.ResponseWriter, r *http.Request, timeout time.Duration, logger *zap.SugaredLogger) (Websocket, error) {
	logger.Info("opening websocket connection")

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		logger.Error(err.Error())
		return Websocket{}, ConnectionError
	}

	return Websocket{conn: conn, timeout: timeout, logger: logger}, nil
}

func (w Websocket) Read(request interface{}) error {
	reader := wsutil.NewReader(w.conn, ws.StateServerSide)
	decoder := json.NewDecoder(reader)

	header, err := reader.NextFrame()
	if err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			w.logger.Info("websocket connection deadline exceeded")
			return DeadlineExceededError
		}

		w.logger.Error(err)
		return ReadError
	}

	if header.OpCode == ws.OpClose {
		w.logger.Error("couldn't read message from websocket due to EOF")
		return EOF
	}

	if err := decoder.Decode(&request); err != nil {
		w.logger.Error(err)
		return DecodeError
	}

	if err := w.setDeadline(time.Now().Add(w.timeout)); err != nil {
		return err
	}

	return nil
}

func (w Websocket) Write(data interface{}) error {
	writer := wsutil.NewWriter(w.conn, ws.StateServerSide, ws.OpText)
	encoder := json.NewEncoder(writer)

	if err := encoder.Encode(&data); err != nil {
		w.logger.Error(err.Error())
		return EncodeError
	}

	if err := writer.Flush(); err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			w.logger.Info("websocket connection deadline exceeded")
			return DeadlineExceededError
		}

		w.logger.Error(err.Error())
		return FlushError
	}

	if err := w.setDeadline(time.Now().Add(w.timeout)); err != nil {
		return err
	}

	return nil
}

func (w Websocket) Closed() <-chan CloseReason {
	ch := make(chan CloseReason, 1)

	go func() {
		defer close(ch)
		ch <- w.waitForClosing()
	}()

	return ch
}

func (w Websocket) waitForClosing() CloseReason {
	for {
		header, err := ws.ReadHeader(w.conn)
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				w.logger.Info("websocket connection deadline exceeded")
				return DeadlineExceeded
			}

			if _, ok := err.(*net.OpError); ok {
				w.logger.Info("websocket was closed on the server")
				return ClosedOnServer
			}

			w.logger.Error(err.Error())
			return ErrorOccurred
		}

		if header.OpCode == ws.OpClose {
			return ClosedOnClient
		}
	}
}

func (w Websocket) setDeadline(t time.Time) error {
	if err := w.conn.SetDeadline(t); err != nil {
		w.logger.Error(err.Error())
		return DeadlineSettingError
	}

	return nil
}

func (w Websocket) Close() error {
	w.logger.Infow("closing websocket connection")

	err := w.conn.Close()
	if err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			w.logger.Info("websocket connection deadline exceeded")
			return DeadlineExceededError
		}

		w.logger.Error(err.Error())
		return CloseError
	}

	return nil
}
