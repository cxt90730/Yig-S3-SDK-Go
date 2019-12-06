package s3lib

import (
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
)

func (s3client *S3Client) PutBucketLifeCycle(bucketName string, rules *s3.LifecycleConfiguration) (err error) {
	params := &s3.PutBucketLifecycleInput{
		Bucket: aws.String(bucketName),
		LifecycleConfiguration: rules,

	}
	if _, err = s3client.Client.PutBucketLifecycle(params); err != nil {
		return err
	}
	return
}

func (s3client *S3Client) GetBucketLifeCycle(bucketName string) (rules *s3.GetBucketLifecycleOutput,err error) {
	params := &s3.GetBucketLifecycleInput{
		Bucket: aws.String(bucketName),
	}
	rules, err = s3client.Client.GetBucketLifecycle(params)
	if  err != nil {
		return nil,err
	}
	return rules,nil
}
func (s3client *S3Client) DeleteBucketLifeCycle(bucketName string) (err error) {
	params := &s3.DeleteBucketLifecycleInput{
		Bucket: aws.String(bucketName),
	}
	_, err = s3client.Client.DeleteBucketLifecycle(params)
	if  err != nil {
		return err
	}
	return nil
}

