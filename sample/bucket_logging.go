package sample

import (
	"fmt"
	"github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
)

// WARN: Bucket Logging is not implemented now.
func BucketLoggingSample() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()

	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	// SetBucketLogging(bucketName, logBucketName, "prefix")
	rules := &s3.LoggingEnabled{
		TargetBucket: aws.String("targetbucket"),
		TargetGrants: []*s3.TargetGrant{
			{
				Grantee: &s3.Grantee{
					Type: aws.String("Group"),
					URI:  aws.String("http://acs.amazonaws.com/groups/global/AllUsers"),
				},
				Permission: aws.String("READ"),
			},
		},
		TargetPrefix: aws.String("MyBucketLogs/"),
	}
	err = sc.PutBucketLogging(bucketName, rules)
	if err != nil {
		HandleError(err)
	}
	// GetBucketLogging(bucketName)
	a, err := sc.GetBucketLogging(bucketName)
	if err != nil {
		HandleError(err)
	}
	fmt.Println(a)
	// DeleteBucketLogging(bucketName)
	err = sc.PutBucketLogging(bucketName, nil)
	if err != nil {
		HandleError(err)
	}

	fmt.Printf("BucketLoggingSample Run Success !\n\n")
}
