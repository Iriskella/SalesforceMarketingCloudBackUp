package main

import (
    "context"
    "golang.org/x/oauth2/clientcredentials"
    "net/http"
    "os"
    "fmt"
)

func getMarketingCloudClient() (*http.Client, error) {

    // Load API credentials securely from environment variables
    salesforceClientID := os.Getenv("SALESFORCE_CLIENT_ID")
    salesforceClientSecret := os.Getenv("SALESFORCE_CLIENT_SECRET")
    salesforceTokenURL := os.Getenv("SALESFORCE_TOKEN_URL")

    if salesforceClientID == "" || salesforceClientSecret == "" || salesforceTokenURL == "" {
        fmt.Println("Missing Salesforce API credentials. Please set environment variables.")
        return nil, fmt.Errorf("Missing Salesforce API credentials. Please set environment variables.")
    }

    config := &clientcredentials.Config{
        ClientID:    salesforceClientID,
        ClientSecret: salesforceClientSecret,
        TokenURL:    salesforceTokenURL,
    }

    client := config.Client(context.Background())

    return client, nil
}
