package main 

import (
	"fmt"
    "time"


)

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
	backupContentBlocks(updatedContentBlocks) //might be used for lastRun logic 
    
	if err != nil {
        fmt.Printf("Error fetching updated content blocks: %v\n", err)
        return
    }
    
    // Log script activity and errors
    fmt.Println("Backup completed successfully.")
}
