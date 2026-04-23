package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"server/internal/websocket"
)

type App struct {
	router *gin.Engine
	server *http.Server
	hub    *websocket.Hub
}

func New() (*App, error) {
	container, err := NewContainer()
	if err != nil {
		return nil, err
	}

	hub := websocket.NewHub(container.ChatService)
	websocket.SetRedisClient(container.RedisClient)
	go hub.Run()

	router := SetupRouter(container, hub)

	addr := fmt.Sprintf(":%s", strings.TrimSpace(container.Config.AppPort))
	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return &App{
		router: router,
		server: server,
		hub:    hub,
	}, nil
}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

func init() {
	gin.SetMode(getGinMode())
}

func getGinMode() string {
	mode := strings.TrimSpace(os.Getenv("GIN_MODE"))
	if mode == "" {
		return gin.DebugMode
	}
	return mode
}
