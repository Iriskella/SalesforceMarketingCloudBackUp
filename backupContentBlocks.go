package main 
import (
	"fmt"
	"time"
	"os"
    "strings"

	
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/credentials"	
)


func backupContentBlocks(updatedContentBlocks []ContentBlock){
	
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
    awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
    awsBucketName := os.Getenv("AWS_BUCKET_NAME")
	

    // Create a subfolder with the current date (e.g., "backup_010623")
    currentTime := time.Now()
    backupFolder := fmt.Sprintf("backup_%s", currentTime.Format("010206_1504"))

    // Ensure the backup folder exists or create it
    if err := os.MkdirAll(backupFolder, os.ModePerm); err != nil {
        fmt.Printf("Error creating backup folder: %v\n", err)
        return
    }

	// Initialize an AWS session
	awsSession, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-1"), 
        Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
    })
    if err != nil {
        fmt.Printf("Error creating AWS session: %v\n", err)
        return
    }

    // Create an S3 client
    s3Client := s3.New(awsSession)
    
    // Store fetched content in the local storage
    for _, block := range updatedContentBlocks {
		  // Generate a unique key for each content block within the subfolder
		  key := fmt.Sprintf("%s/%s.txt", backupFolder, block.ID)

		  // Create an S3 object with the content block data
		  content := block.Content// Replace with actual content
		  contentReader := strings.NewReader(content) // For text content
		  // If content is binary, use bytes.NewReader(contentBytes) instead
  
		  _, err := s3Client.PutObject(&s3.PutObjectInput{
			  Bucket: aws.String(awsBucketName),
			  Key:    aws.String(key),
			  Body:   contentReader,
		  })
		  if err != nil {
			  fmt.Printf("Error backing up content block %s to S3: %v\n", block.ID, err)
			  continue // Move on to the next content block on error
		  }
  
		  fmt.Printf("Content block %s backed up to S3 with key: %s\n", block.ID, key)
	  }
  
	  fmt.Println("Backup completed successfully.")
     
}