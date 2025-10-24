package transport

import (
	"context"
	"fmt"

	"github.com/devalv/ldap314ki/internal/config"
	ldap "github.com/go-ldap/ldap/v3"
	"github.com/rs/zerolog/log"
)

// LDAPUser - структура пользователя хранящаяся в домене.
type LDAPUser struct {
	DN             string
	CN             string
	UID            string
	Mail           string
	GivenName      string
	Surname        string
	DisplayName    string
	SAMAccountName string
}

// GetLDAPUsers выполняет подключение к серверу LDAP и поиск пользователей.
func GetLDAPUsers(ctx context.Context, cfg *config.Config) ([]LDAPUser, error) {
	conn, err := ldap.DialURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к ldap: %w", err)
	}

	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Error closing connection with ldap")
		}
	}()

	// Аутентификация
	err = conn.Bind(cfg.LDAPUsername, cfg.LDAPassword)
	if err != nil {
		return nil, fmt.Errorf("ошибка аутентификации в ldap: %w", err)
	}

	// Поиск пользователей
	searchRequest := ldap.NewSearchRequest(
		cfg.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		cfg.LDAPFilter,
		[]string{"dn", "cn", "uid", "mail", "givenName", "sn", "displayName", "sAMAccountName"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска: %w", err)
	}

	// Обработка результатов
	users := make([]LDAPUser, 0)
	for _, entry := range result.Entries {
		log.Debug().Msgf("Found ldap User: %v", entry)
		user := LDAPUser{
			DN:             entry.DN,
			CN:             entry.GetAttributeValue("cn"),
			UID:            entry.GetAttributeValue("uid"),
			Mail:           entry.GetAttributeValue("mail"),
			GivenName:      entry.GetAttributeValue("givenName"),
			Surname:        entry.GetAttributeValue("sn"),
			DisplayName:    entry.GetAttributeValue("displayName"),
			SAMAccountName: entry.GetAttributeValue("sAMAccountName"),
		}
		users = append(users, user)
	}

	return users, nil
}
