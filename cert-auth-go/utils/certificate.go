package utils

import (
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
