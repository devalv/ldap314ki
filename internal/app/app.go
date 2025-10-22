package app

import (
	"context"
	"os"

	"github.com/devalv/ldap314ki/internal/certs"
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
	log.Debug().Msg("Starting the application")

	// lu, err := transport.GetLDAPUsers(ctx, app.cfg)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("ошибка получения пользователей из ldap")
	// }

	// Создание сертификатов для каждого полученного пользователя
	// for _, u := range lu {
	// 	log.Debug().Msgf("Found ldap User: %v", u)
	// }

	certs.GenerateUserCertificate(app.cfg.CACertPath, app.cfg.CAKeyPath, app.cfg.CAPassword, "test")
	app.Stop(ctx)
}

func (app *Application) Stop(ctx context.Context) {
	log.Debug().Msg("Application stopped")
	os.Exit(0)
}
