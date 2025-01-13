package s3buckets

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// afficher les buckets
// afficher le contenu des buckets
// uploader un fichier
// downloader un fichier
// supprimer un fichier
// uploader un dossier ?
// downloader un dossier

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
