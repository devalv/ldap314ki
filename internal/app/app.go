package app

import (
	"context"
	"os"

	"github.com/devalv/ldap314ki/internal/config"
	"github.com/rs/zerolog/log"
)

type Application struct {
	cfg *config.Config
}

func NewApplication(cfg *config.Config) *Application {
	app := &Application{cfg: cfg}

	return app
}

func (app *Application) Start(ctx context.Context) {
	log.Debug().Msg("Starting mail application")
	app.Stop(ctx)
}

func (app *Application) Stop(ctx context.Context) {
	log.Debug().Msg("Application stopped")
	os.Exit(0)
}
