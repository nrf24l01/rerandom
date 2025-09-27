package modes

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	echoMw "github.com/labstack/echo/v4/middleware"
	mainBackendSchemas "github.com/nrf24l01/rerandom/backend/schemas"
	"github.com/nrf24l01/rerandom/gs_sync/config"
	"github.com/nrf24l01/rerandom/gs_sync/redis"
	"github.com/nrf24l01/rerandom/gs_sync/schemas"
	"github.com/nrf24l01/rerandom/gs_sync/task"
)

type Handler struct {
	cfg *config.Config
	redis *redis.RedisClient
} 

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // allow all origins
    },
}

func (h *Handler) wsHandler(c echo.Context) error {
    // Upgrade initial GET request to a websocket
    conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return err
    }
    defer conn.Close()

	var rows []task.SheetRow
	if err := h.redis.LoadStruct(&rows); err != nil {
		log.Printf("Could not load from Redis: %v", err)
	}
	sheet := task.Sheet{
		ClearUsers:    rows,
		ModifiedUsers: rows,
		Actions:       []task.Action{},
	}
		
	// Echo message back to client
	data, _ := json.Marshal(sheet.ModifiedUsers)
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Println("Write error:", err)
		return err
	}

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	// start goroutine to read messages into channel
	msgChan := make(chan []byte)
	errChan := make(chan error)
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			msgChan <- msg
		}
	}()

	// handle both incoming and periodic updates
	for {
		select {
		case <-ticker.C:
			var rows []task.SheetRow
			if err := h.redis.LoadStruct(&rows); err != nil {
				log.Printf("Could not load from Redis: %v", err)
			}
			sheet.ClearUsers, sheet.ModifiedUsers, sheet.Actions = rows, rows, []task.Action{}
			sheet.Rebuild()

			answer, err := json.Marshal(sheet.ModifiedUsers)
			if err != nil {
				log.Println("Marshal error:", err)
				return err
			}

			if err := conn.WriteMessage(websocket.TextMessage, answer); err != nil {
				log.Println("Write error:", err)
				return err
			}

		case msg := <-msgChan:
			log.Printf("Received: %s\n", msg)

			var req schemas.TaskRequest
			if err := json.Unmarshal(msg, &req); err != nil {
				log.Println("Unmarshal error:", err)
				return err
			}

			switch req.Type {
			case 1:
				sheet.ChangeFraction(req.RowId, req.Param)
			case 2:
				sheet.ChangeAlive(req.RowId, req.Param != 0)
			default:
				log.Printf("Unknown request type: %d", req.Type)
				continue
			}

			sheet.Rebuild()

			answer, err := json.Marshal(sheet.ModifiedUsers)
			if err != nil {
				log.Println("Marshal error:", err)
				return err
			}

			if err := conn.WriteMessage(websocket.TextMessage, answer); err != nil {
				log.Println("Write error:", err)
				return err
			}

		case err := <-errChan:
			log.Println("Read error:", err)
			return err
		}
	}
}

func RunWebserver() {
	cfg := config.BuildConfigFromEnv()

	h := &Handler{
		cfg: cfg,
		redis: redis.InitRedisFromCFG(cfg),
	}
    e := echo.New()

	// Configure logging and recovery middleware
	if os.Getenv("RUNTIME_PRODUCTION") != "true" {
		e.Use(echoMw.Logger())
	}
    e.Use(echoMw.Recover())

	e.Use(echoMw.CORSWithConfig(echoMw.CORSConfig{
		AllowOrigins: []string{os.Getenv("ALLOWED_ORIGINS")},
		AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Configure routes
    e.GET("/ws", h.wsHandler)

	e.GET("/ping", func(c echo.Context) error {
	return c.JSON(200, mainBackendSchemas.Message{Status: "Rerandom client predict backend is ok"})
	})

	// !! RUN !!
    log.Println("Server started at :8080")
    e.Logger.Fatal(e.Start(cfg.APP_HOST))
}
