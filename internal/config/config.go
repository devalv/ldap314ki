package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Debug              bool   `yaml:"debug"`
	URL                string `yaml:"url"`
	BaseDN             string `yaml:"baseDn"`
	LDAPUsername       string `yaml:"ldapUsername"`
	LDAPassword        string `yaml:"ldapPassword"`
	LDAPFilter         string `yaml:"ldapFilter"`
	CACertPath         string `yaml:"caCertPath"`
	CAKeyPath          string `yaml:"caKeyPath"`
	CAPassword         string `yaml:"caPassword"`
	CertValidityDays   int    `yaml:"certValidityDays"`
	CertKeySize        int    `yaml:"certKeySize"`
	UserCertSaveToPath string `yaml:"userCertSaveToPath"`

	ConfigPath string
}

func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to validate config path: %w", err)
	}

	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file", path)
	}

	return nil
}

// Проверяем наличие каталога указанного для сохранения пользовательских сертификатов.
func validateCertsPath(path string) error {
	// Права доступа на каталог
	const (
		cerDirPerm = 0o700 // -rwx--
	)

	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Info().Msgf("Certs path '%s' does not exist, will be created", path)
		err := os.MkdirAll(path, cerDirPerm)
		if err != nil {
			return fmt.Errorf("failed to create certs path: %w", err)
		}

		return nil
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("'%s' is not a directory", path)
	}

	return nil
}

func parseFlags() (path string, err error) {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "./config.yml", "path to config file")
	flag.Parse()

	if err := validateConfigPath(cfgPath); err != nil {
		return "", fmt.Errorf("failed to parse application flags: %w", err)
	}

	return cfgPath, nil
}

func NewConfig() (*Config, error) {
	cfgPath, err := parseFlags()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse flags")
	}

	var cfg Config
	err = cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	cfg.ConfigPath = cfgPath
	cfg.ConfigureLogger()

	if err := validateCertsPath(cfg.UserCertSaveToPath); err != nil {
		log.Fatal().Err(err).Msg("failed to validate certs path")
	}

	return &cfg, nil
}

func (cfg *Config) ConfigureLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug mode enabled")
	}
}
