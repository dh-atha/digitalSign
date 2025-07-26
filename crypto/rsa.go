package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// RSAKeyPair represents an RSA key pair
type RSAKeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// GenerateRSAKeyPair generates a new RSA key pair
func GenerateRSAKeyPair(keySize int) (*RSAKeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %v", err)
	}

	return &RSAKeyPair{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

// PrivateKeyToPEM converts private key to PEM format
func (kp *RSAKeyPair) PrivateKeyToPEM() string {
	privateKeyBytes, _ := x509.MarshalPKCS8PrivateKey(kp.PrivateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return string(privateKeyPEM)
}

// PublicKeyToPEM converts public key to PEM format
func (kp *RSAKeyPair) PublicKeyToPEM() string {
	publicKeyBytes, _ := x509.MarshalPKIXPublicKey(kp.PublicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return string(publicKeyPEM)
}

// ParsePrivateKeyFromPEM parses private key from PEM string
func ParsePrivateKeyFromPEM(pemString string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}

	return rsaPrivateKey, nil
}

// ParsePublicKeyFromPEM parses public key from PEM string
func ParsePublicKeyFromPEM(pemString string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return rsaPublicKey, nil
}

// SignDocument signs a document using RSA private key
func SignDocument(document []byte, privateKey *rsa.PrivateKey) (string, error) {
	hash := sha256.Sum256(document)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign document: %v", err)
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifySignature verifies a document signature using RSA public key
func VerifySignature(document []byte, signature string, publicKey *rsa.PublicKey) error {
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %v", err)
	}

	hash := sha256.Sum256(document)
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signatureBytes)
	if err != nil {
		return fmt.Errorf("signature verification failed: %v", err)
	}

	return nil
}
