package main

import (
	"fmt"
	"os"

	"cert-auth-go/utils" // Import the utils package

	"github.com/spf13/cobra"
)

func main() {
	var certAbsPath string

	var rootCmd = &cobra.Command{
		Use:   "cert-auth-go",
		Short: "A CLI tool for certificate authentication",
		Run: func(cmd *cobra.Command, args []string) {
			if certAbsPath == "" {
				fmt.Println("Error: --cert-abs-path is required")
				os.Exit(1)
			}

			certChain, err := utils.GetCertChainByFilePath(certAbsPath)
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
		},
	}

	rootCmd.Flags().StringVar(&certAbsPath, "cert-abs-path", "", "Absolute path to the certificate file")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
