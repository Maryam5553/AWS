package s3buckets

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// afficher le contenu des buckets
// uploader un fichier
// downloader un fichier
// supprimer un fichier
// uploader un dossier ?
// downloader un dossier
// proposer de créer un bucket mais pour ça il faut gérer la nomenclature compliqué d'amazon S3

// Returns true if the bucket exists, false otherwise.
func BucketExists(s3client *s3.Client, bucketName string) bool {
	// use the HeadBucket method to see if bucket exists and is accessible
	headBucketInput := s3.HeadBucketInput{Bucket: &bucketName}
	_, err := s3client.HeadBucket(context.TODO(), &headBucketInput)
	return err == nil
}

// Print all generate purpose buckets owned on the account.
func ListBuckets(s3client *s3.Client) error {
	// use paginated requests
	maxBuckets := 1000 // maximum number of buckets returned by the request.

	// continue fetching the next buckets
	// until the request indicate that there is no more
	moreBuckets := true
	nbRequest, nbBuckets := 0, 0
	lastContinuationToken := ""
	for moreBuckets {
		// make list buckets request
		listBucketsInput := s3.ListBucketsInput{
			MaxBuckets:        aws.Int32(int32(maxBuckets)),
			ContinuationToken: &lastContinuationToken,
		}
		resp, err := s3client.ListBuckets(context.TODO(), &listBucketsInput)
		if err != nil {
			return fmt.Errorf("couldn't list buckets: %w", err)
		}
		nbRequest++
		nbBuckets += len(resp.Buckets)

		// print them
		for _, bucket := range resp.Buckets {
			fmt.Printf("-%s (%s)\n", *bucket.Name, *bucket.BucketRegion)
		}

		// update continuation token value
		if resp.ContinuationToken == nil {
			moreBuckets = false
			break
		}
		lastContinuationToken = *resp.ContinuationToken
	}

	if nbRequest == 1 && nbBuckets == 0 {
		fmt.Println("No bucket owned.")
	} else {
		fmt.Println("Done listing buckets.")
	}
	return nil
}

// Print all objects contained in a bucket.
func ListObjectsInBucket(s3client *s3.Client, bucketName string) error {
	// check first if bucket exists
	bucketExists := BucketExists(s3client, bucketName)
	if !bucketExists {
		return fmt.Errorf("bucket \"%s\" doesn't exist", bucketName)
	}
	fmt.Printf("%s:\n", bucketName)

	// ListObjectsV2 can return at maximum 1000 results.
	// We use the paginator to make all the necessary list object request
	// to retrieve all the objects from the bucket
	params := &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	}

	listObjectsPaginator := s3.NewListObjectsV2Paginator(s3client, params)
	nbRequest := 0
	nbObjects := 0

	for listObjectsPaginator.HasMorePages() {
		// Makes listObjectsV2 request with pagination
		page, err := listObjectsPaginator.NextPage(context.TODO())
		if err != nil {
			return fmt.Errorf("failed to list all objects from bucket %s: %w", bucketName, err)
		}
		nbRequest++
		nbObjects += len(page.Contents)

		// Print objects
		for _, object := range page.Contents {
			fmt.Printf(" -%s\n", *object.Key)
		}
	}

	if nbRequest == 1 && nbObjects == 0 {
		fmt.Println(" (empty)")
	}
	return nil
}
