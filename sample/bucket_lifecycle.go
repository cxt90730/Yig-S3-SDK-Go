package sample

import (
	"fmt"
	"github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
)

func BucketLifecycleSample() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()

	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	//put bucket lifeCycle
	rules := &s3.LifecycleConfiguration{
		Rules: []*s3.Rule{
			{
				Expiration: &s3.LifecycleExpiration{
					Days: aws.Int64(3650),
				},
				Prefix: aws.String("documents/"),
				ID:     aws.String("TestOnly"),
				Status: aws.String("Enabled"),
				Transition: &s3.Transition{
					Days:         aws.Int64(365),
					StorageClass: aws.String("GLACIER"),
				},
			},
		},
	}

	err = sc.PutBucketLifeCycle(bucketName, rules)
	if err != nil {
		HandleError(err)
	}

	//get bucket lifeCycle
	out, err := sc.GetBucketLifeCycle(bucketName)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("Get Bucket Life Cycle", out)

	//Delete bucket lifeCycle
	err = sc.DeleteBucketLifeCycle(bucketName)
	if err != nil {
		HandleError(err)
	}

	fmt.Printf("BucketLifecycleSample Run Success !\n\n")
}
