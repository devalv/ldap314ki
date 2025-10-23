package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/youmark/pkcs8"
)

// UserCertInfo - атрибуты пользовательского сертификата.
type UserCertInfo struct {
	CommonName   string
	Emails       []string
	ValidityDays int
}

// loadPrivateKey загружает приватный ключ.
func loadPrivateKey(keyFile, password string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(filepath.Clean(keyFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	log.Debug().Msgf("key type: %v", block.Type)
	if block.Type == "ENCRYPTED PRIVATE KEY" {
		// Используем pkcs8 для декодирования PKCS#8 encrypted key
		privateKey, err := pkcs8.ParsePKCS8PrivateKey(block.Bytes, []byte(password))
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt PKCS#8 key: %w", err)
		}

		// Приводим к RSA ключу
		rsaKey, ok := privateKey.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("key is not RSA private key")
		}

		return rsaKey, nil
	}

	if block.Type == "RSA PRIVATE KEY" {
		return nil, errors.New("PKCS#1 format is not supported. Use PKCS#8 instead")
	}

	return nil, fmt.Errorf("unsupported key format: %s", block.Type)
}

// loadCACertificate загружает CA сертификат и приватный ключ.
func loadCACertificate(certPath, keyPath, keyPassword string) (*x509.Certificate, *rsa.PrivateKey, error) {
	// Загрузка CA сертификата
	certData, err := os.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read cert file: %w", err)
	}

	block, _ := pem.Decode(certData)
	if block == nil {
		return nil, nil, errors.New("failed to parse CA certificate PEM")
	}

	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	// Загрузка приватного ключа
	caKey, err := loadPrivateKey(keyPath, keyPassword)
	if err != nil {
		return nil, nil, err
	}

	return caCert, caKey, nil
}

// generateUserCertificateWithExtensions создает пользовательский сертификат с дополнительными расширениями.
func generateUserCertificateWithExtensions(
	caCert *x509.Certificate,
	caKey *rsa.PrivateKey,
	keySize int,
	commonName string,
	emailAddresses []string,
	validityDays int,
) ([]byte, *rsa.PrivateKey, error) {
	// Генерация ключа пользовательского сертификата
	userKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate user cert key: %w", err)
	}

	// Создание шаблона пользовательского сертификата
	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			Country:            caCert.Subject.Country,
			OrganizationalUnit: caCert.Subject.OrganizationalUnit,
			Organization:       caCert.Subject.Organization,
			Locality:           caCert.Subject.Locality,
			Province:           caCert.Subject.Province,
			StreetAddress:      caCert.Subject.StreetAddress,
			PostalCode:         caCert.Subject.PostalCode,
			CommonName:         commonName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 0, validityDays),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
		IsCA:                  false,
		DNSNames:              caCert.DNSNames,
		EmailAddresses:        emailAddresses,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, caCert, &userKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create user cert: %w", err)
	}

	return certDER, userKey, nil
}

// saveCertificateAndKey сохраняет сертификат и ключ в файлы.
func saveCertificateAndKey(certDER []byte, privateKey *rsa.PrivateKey, certPath, keyPath string) error {
	// Права доступа для файлов
	const (
		certPerm = 0o644 // -rw-r--r--
		keyPerm  = 0o600 // -rw-------
	)
	// Сохранение сертификата
	certOut, err := os.OpenFile(filepath.Clean(certPath), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, certPerm)
	if err != nil {
		return fmt.Errorf("failed to save user cert: %w", err)
	}
	defer func() {
		if closeErr := certOut.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Error closing cert file")
		}
	}()

	if err := pem.Encode(certOut, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	}); err != nil {
		return fmt.Errorf("failed to encode user cert: %w", err)
	}

	// Сохранение приватного ключа
	keyOut, err := os.OpenFile(filepath.Clean(keyPath), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, keyPerm)
	if err != nil {
		return fmt.Errorf("failed to save user cert key: %w", err)
	}

	defer func() {
		if closeErr := keyOut.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Error closing cert key file")
		}
	}()

	if err := pem.Encode(keyOut, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		return fmt.Errorf("failed to encode user key: %w", err)
	}

	return nil
}

// GenerateUserCertificate генерирует пользовательский сертификат.
func GenerateUserCertificate(
	caCertPath, caKeyPath, caPass string, keySize int, certInfo UserCertInfo,
) (err error) {
	// Загрузка промежуточного CA
	caCert, caKey, err := loadCACertificate(caCertPath, caKeyPath, caPass)
	if err != nil {
		return fmt.Errorf("ошибка загрузки CA: %w", err)
	}

	// Генерация пользовательского сертификата
	certDER, userKey, err := generateUserCertificateWithExtensions(
		caCert,
		caKey,
		keySize,
		certInfo.CommonName,
		certInfo.Emails,
		certInfo.ValidityDays,
	)
	if err != nil {
		return fmt.Errorf("ошибка генерации пользовательского сертификата: %w", err)
	}

	// Сохранение сертификата и ключа
	// TODO: путь до ключей читать из конфига
	err = saveCertificateAndKey(certDER, userKey, "user_certificate.crt", "user_private.key")
	if err != nil {
		return fmt.Errorf("ошибка сохранения пользовательского сертификата: %w", err)
	}

	log.Debug().Msgf("Сертификат для пользователя %s создан. Выполняем верификацию...", certInfo.CommonName)

	// Верификация сертификата
	userCert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return fmt.Errorf("ошибка верификации пользовательского сертификата: %w", err)
	}

	// Проверка подписи
	err = userCert.CheckSignatureFrom(caCert)
	if err != nil {
		return fmt.Errorf("ошибка проверки  пользовательского сертификата: %w", err)
	}

	log.Debug().Msgf("Сертификат для пользователя %s верифицирован.", certInfo.CommonName)

	return nil
}
