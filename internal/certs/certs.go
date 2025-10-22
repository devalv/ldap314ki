package certs

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/youmark/pkcs8"
)

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

// SaveCertificateAndKey сохраняет сертификат и ключ в файлы
// func SaveCertificateAndKey(certDER []byte, privateKey *rsa.PrivateKey, certPath, keyPath string) error {
// 	// Сохранение сертификата
// 	certOut, err := os.Create(certPath)
// 	if err != nil {
// 		return err
// 	}
// 	defer certOut.Close()

// 	if err := pem.Encode(certOut, &pem.Block{
// 		Type:  "CERTIFICATE",
// 		Bytes: certDER,
// 	}); err != nil {
// 		return err
// 	}

// 	// Сохранение приватного ключа
// 	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
// 	if err != nil {
// 		return err
// 	}
// 	defer keyOut.Close()

// 	if err := pem.Encode(keyOut, &pem.Block{
// 		Type:  "RSA PRIVATE KEY",
// 		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
// 	}); err != nil {
// 		return err
// 	}

// 	return nil
// }

// GenerateUserCertificateWithExtensions создает сертификат с дополнительными расширениями
// func GenerateUserCertificateWithExtensions(
// 	caCert *x509.Certificate,
// 	caKey *rsa.PrivateKey,
// 	commonName string,
// 	dnsNames []string,
// 	emailAddresses []string,
// 	validityDays int,
// ) ([]byte, *rsa.PrivateKey, error) {

// 	userKey, err := rsa.GenerateKey(rand.Reader, 2048)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	template := &x509.Certificate{
// 		SerialNumber: big.NewInt(time.Now().UnixNano()),
// 		Subject: pkix.Name{
// 			CommonName:   commonName,
// 			Organization: []string{"My Organization"},
// 		},
// 		NotBefore: time.Now(),
// 		NotAfter:  time.Now().AddDate(0, 0, validityDays),
// 		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
// 		ExtKeyUsage: []x509.ExtKeyUsage{
// 			x509.ExtKeyUsageClientAuth,
// 			x509.ExtKeyUsageServerAuth,
// 		},
// 		BasicConstraintsValid: true,
// 		IsCA:                  false,
// 		DNSNames:              dnsNames,
// 		EmailAddresses:        emailAddresses,
// 	}

// 	certDER, err := x509.CreateCertificate(rand.Reader, template, caCert, &userKey.PublicKey, caKey)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return certDER, userKey, nil
// }

// func QuickGenerateUserCert(caCertPath, caKeyPath, userCommonName string) error {
// 	// Загрузка CA
// 	caCert, caKey, err := LoadCACertificate(caCertPath, caKeyPath)
// 	if err != nil {
// 		return err
// 	}

// 	// Генерация пользовательского сертификата
// 	certDER, userKey, err := GenerateUserCertificate(caCert, caKey, userCommonName, 365)
// 	if err != nil {
// 		return err
// 	}

// 	// Сохранение
// 	certFile := userCommonName + ".crt"
// 	keyFile := userCommonName + ".key"

// 	return SaveCertificateAndKey(certDER, userKey, certFile, keyFile)
// }

func GenerateUserCertificate(caCertPath, caKeyPath, caPassowrd, userCommonName string) {
	// Загрузка промежуточного CA
	caCert, caKey, err := loadCACertificate(caCertPath, caKeyPath, caPassowrd)
	if err != nil {
		fmt.Printf("\n\nError loading CA: %v\n\n", err)
		return
	}

	// Генерация пользовательского сертификата
	// certDER, userKey, err := GenerateUserCertificateWithExtensions(
	// 	caCert,
	// 	caKey,
	// 	userCommonName,
	// 	[]string{"example.com", "www.example.com"},
	// 	[]string{"user@example.com"},
	// 	365, // 1 год
	// )
	// if err != nil {
	// 	fmt.Printf("Error generating certificate: %v\n", err)
	// 	return
	// }

	// Сохранение сертификата и ключа
	// err = SaveCertificateAndKey(certDER, userKey, "user_certificate.crt", "user_private.key")
	// if err != nil {
	// 	fmt.Printf("Error saving certificate and key: %v\n", err)
	// 	return
	// }

	// fmt.Println("User certificate generated successfully!")

	// Верификация сертификата
	// userCert, err := x509.ParseCertificate(certDER)
	// if err != nil {
	// 	fmt.Printf("Error parsing generated certificate: %v\n", err)
	// 	return
	// }

	// Проверка подписи
	// err = userCert.CheckSignatureFrom(caCert)
	// if err != nil {
	// 	fmt.Printf("Certificate signature verification failed: %v\n", err)
	// } else {
	// 	fmt.Println("Certificate signature verified successfully!")
	// }
}
