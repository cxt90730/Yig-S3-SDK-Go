package s3lib

import (
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
)

func (s3client *S3Client) SetBucketCORS(bucketName string, rules []*s3.CORSRule) error {
	input := &s3.PutBucketCorsInput{
		Bucket: aws.String(bucketName),
		CORSConfiguration: &s3.CORSConfiguration{
			CORSRules: rules,
		},
	}
	_, err := s3client.Client.PutBucketCors(input)
	if err != nil {
		return err
	}
	return nil
}

func (s3client *S3Client) GetBucketCORS(bucketName string) (string, error) {
	input := &s3.GetBucketCorsInput{
		Bucket: aws.String(bucketName),
	}

	out, err := s3client.Client.GetBucketCors(input)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func (s3client *S3Client) DeleteBucketCORS(bucketName string) error {
	input := &s3.DeleteBucketCorsInput{
		Bucket: aws.String(bucketName),
	}
	_, err := s3client.Client.DeleteBucketCors(input)
	if err != nil {
		return err
	}
	return nil
}

