package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/Syipmong/acoustic-sensor/internal/domain"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type AlertUsecaseInterface interface {
	ProcessIncomingFrame(ctx context.Context, rawPayload []byte, alert *domain.Alert) error
}

type AlertHandler struct {
	usecase       AlertUsecaseInterface
	activeClients map[*websocket.Conn]bool
	clientsMu     sync.Mutex
}

func NewAlertHandler(u AlertUsecaseInterface) *AlertHandler {
	return &AlertHandler{
		usecase:       u,
		activeClients: make(map[*websocket.Conn]bool),
	}
}

func (h *AlertHandler) IngestLoraPacket(w http.ResponseWriter, r *http.Request) {
	var alert domain.Alert
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		http.Error(w, "Invalid Payload format", http.StatusBadRequest)
		return
	}

	err := h.usecase.ProcessIncomingFrame(r.Context(), nil, &alert)
	if err != nil {
		log.Printf("Security Alert: %v", err)
		http.Error(w, "Unauthorised Frame", http.StatusUnauthorized)
		return
	}
	h.broadcastToMobileClients(alert)
	w.WriteHeader(http.StatusCreated)
}

func (h *AlertHandler) MobileWebsocketEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade Failed: %v", err)
		return
	}
	h.clientsMu.Lock()
	h.activeClients[ws] = true
	h.clientsMu.Unlock()

	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			h.clientsMu.Lock()
			delete(h.activeClients, ws)
			h.clientsMu.Unlock()
			ws.Close()
			return
		}
	}
}

func (h *AlertHandler) broadcastToMobileClients(alert domain.Alert) {
	h.clientsMu.Lock()
	defer h.clientsMu.Unlock()

	for client := range h.activeClients {
		err := client.WriteJSON(alert)
		if err != nil {
			client.Close()
			delete(h.activeClients, client)
		}
	}
}
