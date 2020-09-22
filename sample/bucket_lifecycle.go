package sample

import (
	"fmt"
	"github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
	"time"
)

func BucketLifecycleSample() {
	DeleteTestBucket()
	defer DeleteTestBucket()

	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	date, err := time.Parse(time.RFC3339, "2020-04-16T00:00:00+08:00")
	if err != nil {
		HandleError(err)
	}

	lifecycle := &s3.BucketLifecycleConfiguration{
		Rules: []*s3.LifecycleRule{
			// AbortIncompleteMultipartUpload,action on all bucket
			{
				AbortIncompleteMultipartUpload: &s3.AbortIncompleteMultipartUpload{
					DaysAfterInitiation: aws.Int64(7),
				},
				ID:     aws.String("id1"),
				Status: aws.String("Enabled"),
			},
			//Expiration,action on "log" dir
			{
				Expiration: &s3.LifecycleExpiration{
					Date: aws.Time(date),
				},
				Filter: &s3.LifecycleRuleFilter{
					Prefix: aws.String("log/"),
				},
				ID:     aws.String("id2"),
				Status: aws.String("Enabled"),
			},
			// NoncurrentVersionExpiration,action on "time" dir
			{
				NoncurrentVersionExpiration: &s3.NoncurrentVersionExpiration{
					NoncurrentDays: aws.Int64(30),
				},
				Filter: &s3.LifecycleRuleFilter{
					And: &s3.LifecycleRuleAndOperator{
						Prefix: aws.String("time/"),
						Tags: []*s3.Tag{
							{
								Key:   aws.String("Key1"),
								Value: aws.String("Value1"),
							},
							{
								Key:   aws.String("Key2"),
								Value: aws.String("Value2"),
							},
						},
					},
				},
				ID:     aws.String("id3"),
				Status: aws.String("Enabled"),
			},
			// Transitions, storageClass turn into GLACIER, action on "prefix" dir
			{
				Transitions: []*s3.Transition{
					{
						Date:         aws.Time(date),
						StorageClass: aws.String("GLACIER"),
					},
				},
				Filter: &s3.LifecycleRuleFilter{
					And: &s3.LifecycleRuleAndOperator{
						Tags: []*s3.Tag{
							{
								Key:   aws.String("Key1"),
								Value: aws.String("Value1"),
							},
							{
								Key:   aws.String("Key2"),
								Value: aws.String("Value2"),
							},
						},
					},
				},
				ID:     aws.String("id4"),
				Status: aws.String("Enabled"),
			},
			// NoncurrentVersionTransitions,storageClass turn into STANDARD_IA,but it's not work
			//{
			//	Filter: &s3.LifecycleRuleFilter{
			//		Prefix: aws.String("test/"),
			//	},
			//	ID:     aws.String("id5"),
			//	Status: aws.String("Disable"),
			//	NoncurrentVersionTransitions: []*s3.NoncurrentVersionTransition{
			//		{
			//			NoncurrentDays: aws.Int64(30),
			//			StorageClass:   aws.String("STANDARD_IA"),
			//		},
			//	},
			//},
		},
	}

	err = sc.PutBucketLifecycle(bucketName,lifecycle)
	if err != nil {
		HandleError(err)
	}

	out, err := sc.GetBucketLifecycle(bucketName)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("Get Bucket LifecycleConfiguration", out)

	out, err = sc.DeleteBucketLifecycle(bucketName)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("Delete Bucket LifecycleConfiguration", out)


	fmt.Printf("BucketLifecycleSample Run Success !\n\n")
}
