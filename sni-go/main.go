package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

func main() {
	vaultUrl := getEnv("VAULT_URL")
	certName := getEnv("CERT_NAME")
	clientId := getEnv("CLIENT_ID")
	tenantId := getEnv("TENANT_ID")
	scope := getEnv("SCOPE")

	cliCred, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to obtain a credential: %v\n", err)
		os.Exit(1)
	}

	sniCred, err := getCertificateTokenCredential(tenantId, clientId, vaultUrl, certName, cliCred)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get certificate token credential: %v\n", err)
		os.Exit(1)
	}

	token, err := getToken(sniCred, scope)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get token: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(token.Token)
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		fmt.Fprintf(os.Stderr, "environment variable %s is required\n", key)
		os.Exit(1)
	}
	return value
}

func getToken(cred azcore.TokenCredential, scope string) (azcore.AccessToken, error) {
	ctx := context.Background()
	requestContext := policy.TokenRequestOptions{
		Scopes: []string{scope},
	}
	return cred.GetToken(ctx, requestContext)
}

func getCertificateFromKeyVault(vaultUrl, certName string, cred azcore.TokenCredential) ([]byte, error) {
	client, err := azsecrets.NewClient(vaultUrl, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret client: %v", err)
	}

	resp, err := client.GetSecret(context.Background(), certName, "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %v", err)
	}

	privateKeyBytes, err := base64.StdEncoding.DecodeString(*resp.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret value: %v", err)
	}

	return privateKeyBytes, nil
}

func getCertificateTokenCredential(tenantId, clientId, vaultUrl, certName string, cred azcore.TokenCredential) (azcore.TokenCredential, error) {
	privateKeyBytes, err := getCertificateFromKeyVault(vaultUrl, certName, cred)
	if err != nil {
		return nil, err
	}

	c, err := cache.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %v", err)
	}

	certs, key, err := azidentity.ParseCertificates(privateKeyBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificates: %v", err)
	}

	options := &azidentity.ClientCertificateCredentialOptions{
		SendCertificateChain: true,
		Cache:                c,
	}

	return azidentity.NewClientCertificateCredential(tenantId, clientId, certs, key, options)
}
