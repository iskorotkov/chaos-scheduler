package ws

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

type Websocket struct {
	conn   net.Conn
	Closed <-chan CloseReason
	logger *zap.SugaredLogger
}

func NewWebsocket(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger) (Websocket, error) {
	logger.Info("opening websocket connection")

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		logger.Error(err.Error())
		return Websocket{}, ConnectionError
	}

	ch := make(chan CloseReason, 1)

	socket := Websocket{conn: conn, Closed: ch, logger: logger}

	go func() {
		defer close(ch)
		ch <- socket.waitForClosing()
	}()

	return socket, nil
}

func (w Websocket) Read(ctx context.Context, data interface{}) error {
	if ctx.Err() != nil {
		if ctx.Err() == context.Canceled {
			return ContextCancelledError
		} else {
			return DeadlineExceededError
		}
	}

	if err := w.setDeadline(ctx); err != nil {
		return err
	}

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
		return EOFError
	}

	if err := decoder.Decode(&data); err != nil {
		w.logger.Error(err)
		return DecodeError
	}

	return nil
}

func (w Websocket) Write(ctx context.Context, data interface{}) error {
	if ctx.Err() != nil {
		if ctx.Err() == context.Canceled {
			return ContextCancelledError
		} else {
			return DeadlineExceededError
		}
	}

	if err := w.setDeadline(ctx); err != nil {
		return err
	}

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

func (w Websocket) setDeadline(ctx context.Context) error {
	t, ok := ctx.Deadline()
	if !ok {
		t = time.Time{}
	}

	if err := w.conn.SetDeadline(t); err != nil {
		w.logger.Error(err)
		return DeadlineSettingError
	}

	return nil
}

func (w Websocket) waitForClosing() CloseReason {
	for {
		header, err := ws.ReadHeader(w.conn)
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				return ReasonDeadlineExceeded
			}

			if _, ok := err.(*net.OpError); ok {
				return ReasonClosedOnServer
			}

			if err == io.EOF {
				return ReasonEOF
			}

			w.logger.Warn(err.Error())
			return ReasonErrorOccurred
		}

		if header.OpCode == ws.OpClose {
			return ReasonClosedOnClient
		}
	}
}
