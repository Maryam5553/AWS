package main

import (
	"aws/pkg/s3buckets"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Ask user for bucket name and print the content of this bucket.
func listContentBucket(s3client *s3.Client) error {
	// get bucket name from user
	fmt.Print("bucket name: ")
	var bucketName string
	_, err := fmt.Scan(&bucketName)
	if err != nil {
		return fmt.Errorf("error reading user input: %w", err)
	}

	// print objects in that bucket
	return s3buckets.ListObjectsInBucket(s3client, bucketName)
}

func main() {
	// loads AWS user configuration from the files ~/.aws/config (to retrieve the AWS region)
	// and ~/.aws/credentials (to retrieve the user AWS access key)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// creates a service client to perform actions on S3
	s3client := s3.NewFromConfig(cfg)

	// menu
	fmt.Println("This program can perfom S3 operations on your account. Possible actions:")
	fmt.Println("1- List all buckets")
	fmt.Println("2- List the content of a bucket")
	fmt.Println("3- List all objects in S3")
	fmt.Println("4- Exit program")

	// ask user which operation they want to perform
	for {
		fmt.Print("\n> Enter a number: ")

		var answer string
		_, err = fmt.Scan(&answer)
		if err != nil {
			log.Printf("error reading user input: %v", err)
			continue
		}

		switch answer {
		case "1":
			err = s3buckets.ListBuckets(s3client)
			if err != nil {
				log.Println(err)
			}
		case "2":
			err = listContentBucket(s3client)
			if err != nil {
				log.Println(err)
			}
		case "3":
			err = s3buckets.ListAllObjects(s3client)
			if err != nil {
				log.Println(err)
			}
		case "4":
			return
		default:
			fmt.Print("Invalid output, please answer a number between 1 and 6.")
		}
	}
}
