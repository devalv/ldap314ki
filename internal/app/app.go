package app

import (
	"context"
	"os"

	"github.com/devalv/ldap314ki/internal/certs"
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

	ldapUsers, err := transport.GetLDAPUsers(ctx, app.cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("ошибка получения пользователей из ldap")
	}

	// Создание сертификатов для каждого полученного пользователя
	for _, user := range ldapUsers {
		log.Debug().Msgf("Found ldap User: %v", user.CN)
		err := certs.GenerateUserCertificate(
			app.cfg.CACertPath, app.cfg.CAKeyPath, app.cfg.CAPassword, app.cfg.CertKeySize, certs.UserCertInfo{
				CommonName:         user.CN,
				Emails:             []string{user.Mail},
				ValidityDays:       app.cfg.CertValidityDays,
				SAMAccountName:     user.SAMAccountName,
				UserCertSaveToPath: app.cfg.UserCertSaveToPath,
			})
		if err != nil {
			log.Fatal().Err(err).Msg("ошибка создания сертификата")
		}
		log.Info().Msgf("Сертификат для пользователя %s создан", user.CN)
	}

	app.Stop(ctx)
}

func (app *Application) Stop(ctx context.Context) {
	log.Debug().Msg("Application stopped")
	os.Exit(0)
}
