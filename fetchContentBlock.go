package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ContentBlock struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	Content   string `json:"content"`
}

func fetchUpdatedContent(lastRun time.Time, client *http.Client) ([]ContentBlock, error) {
	// Define the Content Builder API endpoint
	apiEndpoint := "https://your-marketing-cloud-instance.rest.marketingcloudapis.com/asset/v1/content/assets/query"

	// Create a JSON payload for the Content Builder API request
	requestPayload := map[string]interface{}{
		"query": map[string]interface{}{
			"property":       "createdAt",
			"simpleOperator": "greaterThanOrEqual",
			"value":          lastRun.Format(time.RFC3339),
		},
	}

	requestBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, err
	}

	// Create an HTTP POST request to the Content Builder API
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}

    salesforceTokenURL := os.Getenv("SALESFORCE_TOKEN_URL")


	req.Header.Set("Authorization", salesforceTokenURL)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	// Read the API response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response
	var contentBlocks []ContentBlock
	if err := json.Unmarshal(body, &contentBlocks); err != nil {
		return nil, err
	}

	return contentBlocks, nil
}
