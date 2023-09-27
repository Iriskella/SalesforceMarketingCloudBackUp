package main 

import (
	"fmt"
    "time"


)
var lastUpdatedContentBlock string

func main() {
    //Authenticate with Salesforce Marketing Cloud API
    client, err := getMarketingCloudClient()
    if err != nil {
        fmt.Printf("Error authenticating with Salesforce Marketing Cloud: %v\n", err)
        return
    }

    // Fetch content blocks since the last run
    lastRun := time.Now().Add(-24 * time.Hour) // Assuming the schedule is once per day as required  
    updatedContentBlocks, err := fetchUpdatedContent(lastRun, client)
    
	if err != nil {
        fmt.Printf("Error fetching updated content blocks: %v\n", err)
        return
    }

    result, err := backupContentBlocks(updatedContentBlocks)
	if err != nil {
		fmt.Println("Backup failed:", err)
	} else {
        lastUpdatedContentBlock = result.LastCompletedBlockID
		fmt.Println("Backup completed successfully. Last completed block ID:", result.LastCompletedBlockID) //to be used as lastRun logic
		if len(result.Errors) > 0 {
			fmt.Println("Backup completed with errors. Error count:", len(result.Errors))
		}
	}

    fmt.Println("Backup completed successfully.")
}
