package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type BlockWithError struct {
	BlockID string
	Err   error
}
type BackupResult struct {
	LastCompletedBlockID string
	Errors               []BlockWithError
}

func backupContentBlocks(updatedContentBlocks []ContentBlock)(*BackupResult,error){
	
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
    awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
    awsBucketName := os.Getenv("AWS_BUCKET_NAME")
	

    // Create a subfolder with the current date (e.g., "backup_010623")
    currentTime := time.Now()
    backupFolder := fmt.Sprintf("backup_%s", currentTime.Format("010206_1504"))

    // Ensure the backup folder exists or create it
    if err := os.MkdirAll(backupFolder, os.ModePerm); err != nil {
        fmt.Printf("Error creating backup folder: %v\n", err)
        return nil, errors.New(fmt.Sprintf("Error creating backup folder: %v\n", err))
    }

	// Initialize an AWS session
	awsSession, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-1"), 
        Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
    })
    if err != nil {
        fmt.Printf("Error creating AWS session: %v\n", err)
		return nil, errors.New(fmt.Sprintf("Error creating AWS session: %v\n", err))

    }

    // Create an S3 client
    s3Client := s3.New(awsSession)
	blockErrSlice := []BlockWithError{}
	lastCompletedBlockID := ""


    // Store fetched content in the local storage
    for _, block := range updatedContentBlocks {
		  // Generate a unique key for each content block within the subfolder
		  key := fmt.Sprintf("%s/%s.txt", backupFolder, block.ID)

		  // Create an S3 object with the content block data
		  content := block.Content
		  contentReader := strings.NewReader(content) // For text content
		  
  
		  _, err := s3Client.PutObject(&s3.PutObjectInput{
			  Bucket: aws.String(awsBucketName),
			  Key:    aws.String(key),
			  Body:   contentReader,
		  })
		  if err != nil {
			  fmt.Printf("Error backing up content block %s to S3: %v\n", block.ID, err)
			  blockErrSlice = append(blockErrSlice, BlockWithError{BlockID: block.ID, Err: err})
			  continue 
		  }
		  lastCompletedBlockID = block.ID
		  fmt.Printf("Content block %s backed up to S3 with key: %s\n", block.ID, key)
	  }
	  return &BackupResult{
		LastCompletedBlockID: lastCompletedBlockID,
		Errors:               blockErrSlice,
	}, nil     
}