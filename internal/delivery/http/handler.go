package http

import (
	"acoustic-sensor-backend/internal/domain"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Syipmong/acoustic-sensor-backend/internal/domain"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool{ return true},
}

type AlertUsecaseInterface interface {
	ProcessIncomingFrame(ctx context.Context, rawPayload []byte, alert *domain.Alert) error
}

type AlertHandler struct {
	usecase AlertUsecaseInterface
	activeClients  map[*websocket.Conn]bool
}

func NewAlertHAndler(u AlertUsecaseInterface) *AlertHandler{
	return &AlertHandler{
		usecase: u,
		activeClients: make(map[*websocket.Conn]bool),
	}
}