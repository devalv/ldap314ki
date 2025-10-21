package app

import (
	"context"
	"fmt"
	"os"

	"github.com/devalv/ldap314ki/internal/config"
	transport "github.com/devalv/ldap314ki/internal/transport/ldap"
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
	log.Debug().Msg("Starting the application")

	lu, err := app.getLdapUsers(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get ldap users")
	}

	for _, u := range lu {
		log.Debug().Msgf("Found ldap User: %v", u)
	}
	app.Stop(ctx)
}

func (app *Application) Stop(ctx context.Context) {
	log.Debug().Msg("Application stopped")
	os.Exit(0)
}

func (app *Application) getLdapUsers(ctx context.Context) (lu []transport.LDAPUser, err error) {
	f, err := transport.GetLDAPUsers(ctx, app.cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователей из ldap: %w", err)
	}

	return f, nil
}
