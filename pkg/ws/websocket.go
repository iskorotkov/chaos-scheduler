package ws

import (
	"encoding/json"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"go.uber.org/zap"
	"net"
	"net/http"
)

type Websocket struct {
	conn   net.Conn
	logger *zap.SugaredLogger
}

func NewWebsocket(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger) (Websocket, error) {
	logger.Info("opening websocket connection")

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		logger.Error(err.Error())
		return Websocket{}, ConnectionError
	}

	return Websocket{conn: conn, logger: logger}, nil
}

func (w Websocket) Read(request interface{}) error {
	reader := wsutil.NewReader(w.conn, ws.StateServerSide)
	decoder := json.NewDecoder(reader)

	header, err := reader.NextFrame()
	if err != nil {
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
		w.logger.Error(err.Error())
		return FlushError
	}

	return nil
}

func (w Websocket) Close() error {
	w.logger.Infow("closing websocket connection")

	err := w.conn.Close()
	if err != nil {
		w.logger.Error(err.Error())
		return CloseError
	}

	return nil
}
