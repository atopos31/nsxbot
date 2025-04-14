package driver

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/atopos31/nsxbot/nlog"
	"github.com/atopos31/nsxbot/types"
	"github.com/gorilla/websocket"
)

type WServer struct {
	conn *websocket.Conn
	url  url.URL
	log  *slog.Logger
}

func NewWSverver(host string, path string) *WServer {
	return &WServer{
		url: url.URL{
			Scheme: "ws",
			Host:   host,
			Path:   path,
		},
		log: nlog.Logger(),
	}
}

func (ws *WServer) Listen(ctx context.Context, eventChan chan<- types.Event) error {
	mux := http.NewServeMux()
	mux.HandleFunc(ws.url.Path, func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			ws.log.Error("Upgrade", "err", err)
			return
		}
		defer c.Close()
		for {
			_, content, err := c.ReadMessage()
			if err != nil {
				ws.log.Error("Read", "err", err)
				break
			}
			event, err := contentToEvent(content)
			if err != nil {
				ws.log.Error("Invalid event", "err", err)
				continue
			}
			eventChan <- event
		}
	})
	ws.log.Info("WS listener start... ", "addr", ws.url.Host)
	server := &http.Server{Addr: ws.url.Host, Handler: mux}
	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			ws.log.Error("WS server shutdown error", "err", err)
			return
		}
	}()
	return server.ListenAndServe()
}
