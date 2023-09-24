package main

import (
    "fmt"
    "time"
    "net/http"
    "io"
    "encoding/json"
)

type ContentBlock struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    CreatedAt string `json:"created_at"`
	Content   string `json:"content"`
}

func fetchUpdatedContent(lastRun time.Time, client *http.Client) ([]ContentBlock, error) {
    apiEndpoint := "https://marketing-cloud-api-endpoint.com/content-blocks"

    req, err := http.NewRequest("GET", apiEndpoint, nil)
    if err != nil {
        return nil, err
    }

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

    // Filter content blocks based on the last run date
    updatedContentBlocks := []ContentBlock{}
    for _, block := range contentBlocks {
        createdAt, err := time.Parse(time.RFC3339, block.CreatedAt)
        if err != nil {
            continue
        }

        if createdAt.After(lastRun) {
            updatedContentBlocks = append(updatedContentBlocks, block)
        }
    }

    return updatedContentBlocks, nil
}
