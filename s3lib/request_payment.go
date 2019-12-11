package s3lib

import (
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
)

func (s3client *S3Client) PutBucketRequestPayment(bucketName string) error {
	params := &s3.PutBucketRequestPaymentInput{
		Bucket: aws.String(bucketName),
		RequestPaymentConfiguration: &s3.RequestPaymentConfiguration{
			Payer: aws.String("Requester"),
		},
	}
	_, err := s3client.Client.PutBucketRequestPayment(params)
	if err != nil {
		return err
	}
	return nil
}

func (s3client *S3Client) GetBucketRequestPayment(bucketName string) (result *s3.GetBucketRequestPaymentOutput, err error) {
	params := &s3.GetBucketRequestPaymentInput{
		Bucket: aws.String(bucketName),
	}
	result, err = s3client.Client.GetBucketRequestPayment(params)
	if err != nil {
		return nil, err
	}
	return result, nil
}
