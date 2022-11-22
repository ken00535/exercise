package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	dhttp "assignment/internal/app/delivery/http"
	"assignment/internal/app/repository/db"
	"assignment/internal/app/usecase"
	"assignment/internal/infra/config"
	"assignment/internal/infra/fxlogger"
	wsgin "assignment/internal/infra/gin"
	wszerolog "assignment/internal/infra/zerolog"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	zerolog "github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

var (
	srv *http.Server
)

func main() {
	config.Init("./config", "app")
	cfg := config.Get()
	wszerolog.Init(cfg.Log)

	zerolog.Info().Msg("starting shorten url")
	defer func() {
		zerolog.Info().Msg("exiting shorten url")
	}()
	defer recoverFn()
	var e *gin.Engine
	app := fx.New(
		fx.Supply(cfg),
		fx.Provide(
			wsgin.New,
			db.New,
			usecase.New,
			dhttp.New,
		),
		fx.WithLogger(fxlogger.NewLogger),
		fx.Populate(&e),
		fx.Invoke(
			startServe,
		),
	)
	if err := app.Start(context.Background()); err != nil {
		zerolog.Err(err).Msg(err.Error())
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan

	stopServe(app)

}

func startServe(e *gin.Engine, hd *dhttp.Delivery) {
	go func() {
		cfg := config.Get()
		addr := strings.Split(cfg.Http.Address, ":")
		srv = &http.Server{
			Addr:     addr[0] + ":80",
			Handler:  e,
			ErrorLog: log.New(io.Discard, "", 0),
		}
		zerolog.Info().Msg("starting http server")
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			zerolog.Error().Err(err).Msgf("%v", err)
		}
	}()
}

func stopServe(app *fx.App) {
	zerolog.Info().Msg("stopping http server")
	err := srv.Shutdown(context.Background())
	if err != nil {
		zerolog.Err(err).Msgf("%s", err.Error())
	}
	err = app.Stop(context.Background())
	if err != nil {
		zerolog.Err(err).Msgf("%s", err.Error())
	}
}

func recoverFn() {
	if r := recover(); r != nil {
		var msg string
		for i := 2; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			msg += fmt.Sprintf("[%s:%d]\n", file, line)
		}
		err, ok := r.(error)
		if !ok {
			zerolog.Error().Str("stack", msg).Msg("panic")
			return
		}
		zerolog.Err(err).Str("stack", msg).Msg(err.Error())
	}
}
