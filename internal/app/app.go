package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cecepsprd/starworks-test/config"
	"github.com/cecepsprd/starworks-test/internal/handler"
	"github.com/cecepsprd/starworks-test/internal/repository"
	"github.com/cecepsprd/starworks-test/internal/service"
	"github.com/cecepsprd/starworks-test/utils/logger"
	"github.com/cecepsprd/starworks-test/utils/validate"
	en "github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func RunServer() {
	cfg := config.NewConfig()

	db, err := cfg.MysqlConnect()
	if err != nil {
		log.Fatal("error connecting to database: ", err.Error())
	}

	if err = logger.Init(cfg.App.LogLevel, cfg.App.LogTimeFormat); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
		},
	}))

	timeoutContext := time.Duration(cfg.App.ContextTimeout) * time.Second

	customValidator := validate.NewValidator()
	en.RegisterDefaultTranslations(customValidator.Validator, customValidator.Translator)
	e.Validator = customValidator

	userRepository := repository.NewUserRepository(db)
	walletRepository := repository.NewWalletRepository(db)

	userService := service.NewUserService(userRepository, walletRepository, cfg.App.JWTSecret, timeoutContext)
	walletService := service.NewWalletService(walletRepository)

	handler.NewUserHandler(e, userService)
	handler.NewWalletHandler(e, walletService)

	e.GET("/api/check", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK!")
	})

	// Starting server
	go func() {
		if cfg.App.HTTPPort == "" {
			cfg.App.HTTPPort = os.Getenv("PORT")
		}

		err := e.Start(":" + cfg.App.HTTPPort)
		if err != nil {
			log.Fatal("error starting server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	// Block until a signal is received.
	<-quit

	log.Println("server shutdown of 5 second.")

	// gracefully shutdown the server, waiting max 5 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e.Shutdown(ctx)
}
