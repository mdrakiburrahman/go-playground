package utils

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/crypto/pkcs12"
)

func decodePkcs12Chain(pfxData []byte, password string) (certChain []*x509.Certificate, err error) {
	pemBlocks, err := pkcs12.ToPEM(pfxData, password)
	if err != nil {
		return nil, err
	}

	for _, block := range pemBlocks {
		if block.Type != "CERTIFICATE" {
			continue
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("error parsing certificate(s): %v", err)
		}
		certChain = append(certChain, cert)
	}

	return certChain, nil
}

func decodePkcs12(pkcs []byte, password string) (*x509.Certificate, *rsa.PrivateKey, error) {
	privateKey, certificate, err := pkcs12.Decode(pkcs, password)
	if err != nil {
		return nil, nil, err
	}

	rsaPrivateKey, isRsaKey := privateKey.(*rsa.PrivateKey)
	if !isRsaKey || rsaPrivateKey == nil {
		return nil, nil, fmt.Errorf("PKCS#12 certificate must contain an RSA private key")
	}

	return certificate, rsaPrivateKey, nil
}

// GetCertByFilePath guarantees to return valid (non-nil) certificate and private key on nil error
func GetCertByFilePath(certificatePath string) (certificate *x509.Certificate, rsaPrivateKey crypto.PrivateKey, err error) {
	certbytes, err := os.ReadFile(certificatePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read the certificate file (%s): %v", certificatePath, err)
	}
	certBase64Data, err := base64.StdEncoding.DecodeString(string(certbytes))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode certificate (1): %v", err)
	}
	certificate, rsaPrivateKey, err = decodePkcs12(certBase64Data, "")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode certificate (2): %v", err)
	}

	return certificate, rsaPrivateKey, nil
}

// GetCertChainByFilePath reads a certificate file, decodes it from base64, and returns the certificate chain.
func GetCertChainByFilePath(certificatePath string) (certChain []*x509.Certificate, err error) {

	certbytes, err := os.ReadFile(certificatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read the certificate file (%s): %v", certificatePath, err)
	}
	certBase64Data, err := base64.StdEncoding.DecodeString(string(certbytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode certificate (1): %v", err)
	}
	certChain, err = decodePkcs12Chain(certBase64Data, "")
	if err != nil {
		return nil, fmt.Errorf("failed to decode certificate chain: %v", err)
	}
	if len(certChain) == 0 {
		return nil, fmt.Errorf("no certificate found in the chain")
	}

	return certChain, nil
}
