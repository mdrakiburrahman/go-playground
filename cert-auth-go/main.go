package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"cert-auth-go/utils"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/spf13/cobra"
)

func main() {
	var certAbsPath, tenantID, clientID, scope string

	var rootCmd = &cobra.Command{
		Use:   "cert-auth-go",
		Short: "A CLI tool for certificate authentication",
		Run: func(cmd *cobra.Command, args []string) {
			if certAbsPath == "" {
				fmt.Println("Error: --cert-abs-path is required")
				os.Exit(1)
			}
			if tenantID == "" {
				fmt.Println("Error: --tenant-id is required")
				os.Exit(1)
			}
			if clientID == "" {
				fmt.Println("Error: --client-id is required")
				os.Exit(1)
			}
			if scope == "" {
				fmt.Println("Error: --scope is required")
				os.Exit(1)
			}

			certChain, err := utils.GetCertChainByFilePath(certAbsPath)
			if err != nil {
				fmt.Printf("Failed to get certificate chain: %v\n", err)
				os.Exit(1)
			}

			cert, rsaPrivateKey, err := utils.GetCertByFilePath(certAbsPath)
			if err != nil {
				fmt.Printf("Failed to get certificate chain: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Certificate Chain:")
			for i, cert := range certChain {
				fmt.Printf("Certificate %d:\n", i+1)
				fmt.Printf("  Subject: %s\n", cert.Subject)
				fmt.Printf("  Issuer: %s\n", cert.Issuer)
			}

			fmt.Println("Certificate:")
			fmt.Printf("  Subject: %s\n", cert.Subject)
			fmt.Printf("  Issuer: %s\n", cert.Issuer)
			rsaKey, ok := rsaPrivateKey.(*rsa.PrivateKey)
			if !ok {
				fmt.Println("Error: Private key is not of type RSA")
				os.Exit(1)
			}
			fmt.Printf("  RSA Private Key Length: %d bits\n", rsaKey.Size()*8)
			fmt.Println("Certificate and RSA Private Key retrieved successfully.")

			cred, err := azidentity.NewClientCertificateCredential(
				tenantID,
				clientID,
				[]*x509.Certificate{cert},
				rsaPrivateKey,
				&azidentity.ClientCertificateCredentialOptions{ClientOptions: policy.ClientOptions{}, DisableInstanceDiscovery: true, SendCertificateChain: true},
			)
			if err != nil {
				fmt.Printf("Failed to create ClientCertificateCredential: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("ClientCertificateCredential created successfully.")

			token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
				Scopes: []string{scope},
			})
			if err != nil {
				fmt.Printf("Failed to get token: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Access Token: %s\n", token.Token)
			fmt.Printf("Expires On: %s\n", token.ExpiresOn.Format(time.RFC3339))
			fmt.Println("Token retrieved successfully.")
		},
	}

	rootCmd.Flags().StringVar(&certAbsPath, "cert-abs-path", "", "Absolute path to the certificate file")
	rootCmd.Flags().StringVar(&tenantID, "tenant-id", "", "Azure Tenant ID")
	rootCmd.Flags().StringVar(&clientID, "client-id", "", "Azure Client ID")
	rootCmd.Flags().StringVar(&scope, "scope", "https://management.core.windows.net/.default", "Azure scope for token")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
