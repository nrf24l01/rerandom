package main

import (
	"html/template"

	"github.com/go-playground/validator"
	"github.com/nrf24l01/go-web-utils/echokit"
	"github.com/nrf24l01/rerandom/backend/core"
	"github.com/nrf24l01/rerandom/backend/handlers"
	"github.com/nrf24l01/rerandom/backend/models"
	"github.com/nrf24l01/rerandom/backend/routes"
	"github.com/nrf24l01/rerandom/backend/schemas"
	"github.com/nrf24l01/rerandom/backend/templater"

	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
)
func main() {
	if os.Getenv("PRODUCTION_ENV") != "true" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("failed to load .env: %v", err)
		}
	}
	
	config, err := core.BuildConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to build config: %v", err)
	}

	e := echo.New()

	db := models.RegisterPostgres(config)

	// Register custom validator
	e.Validator = &echokit.CustomValidator{Validator: validator.New()}

	// Инициализация шаблонов
    renderer := &templater.TemplateRenderer{
        Templates: template.Must(template.ParseGlob("templates/*.html")),
    }
    e.Renderer = renderer

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

	e.GET("/ping", func(c echo.Context) error {
	return c.JSON(200, schemas.Message{Status: "Rerandom backend is ok"})
	})

	handler := &handlers.Handler{DB: db, Config: config}
	routes.RegisterRoutes(e, handler)
	
	e.Logger.Fatal(e.Start(config.APPHost))
}