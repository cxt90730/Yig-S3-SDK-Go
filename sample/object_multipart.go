package sample

import (
	"fmt"
	. "github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
)

func MultiPartUploadSample() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()
	sc := NewS3(endpoint, accessKey, secretKey)
	// Create a bucket
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	// 1.Create Multipart Upload
	uploadId, err := sc.CreateMultiPartUpload(bucketName, objectKey, s3.ObjectStorageClassStandard)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("UploadId is: ", aws.StringValue(aws.String(uploadId)))
	partCount := 3
	completedUpload := &s3.CompletedMultipartUpload{
		Parts: make([]*s3.CompletedPart, partCount),
	}
	for i := 0; i < partCount; i++ {
		partNumber := int64(i + 1)
		etag, err := sc.UploadPart(bucketName, objectKey, GenMinimalPart(), uploadId, partNumber)
		if err != nil {
			HandleError(err)
		}
		completedUpload.Parts[i] = &s3.CompletedPart{
			ETag:       aws.String(etag),
			PartNumber: aws.Int64(partNumber),
		}
	}
	// 2.Upload Part
	err = sc.CompleteMultiPartUpload(bucketName, objectKey, uploadId, completedUpload)
	if err != nil {
		HandleError(err)
		err = sc.AbortMultiPartUpload(bucketName, objectKey, uploadId)
		if err != nil {
			HandleError(err)
		}
	}
	fmt.Printf("MultipartUploadSample Run Success !\n\n")
}

func MultiPartUploadSampleWithForbidOverwrite() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()
	sc := NewS3(endpoint, accessKey, secretKey)
	// Create a bucket
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	uploadId, err := sc.CreateMultiPartUpload(bucketName, objectKey, s3.ObjectStorageClassStandard)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("UploadId is: ", aws.StringValue(aws.String(uploadId)))
	partCount := 3
	completedUpload := &s3.CompletedMultipartUpload{
		Parts: make([]*s3.CompletedPart, partCount),
	}
	for i := 0; i < partCount; i++ {
		partNumber := int64(i + 1)
		etag, err := sc.UploadPart(bucketName, objectKey, GenMinimalPart(), uploadId, partNumber)
		if err != nil {
			HandleError(err)
		}
		completedUpload.Parts[i] = &s3.CompletedPart{
			ETag:       aws.String(etag),
			PartNumber: aws.Int64(partNumber),
		}
	}
	err = sc.CompleteMultiPartUpload(bucketName, objectKey, uploadId, completedUpload)
	if err != nil {
		HandleError(err)
		err = sc.AbortMultiPartUpload(bucketName, objectKey, uploadId)
		if err != nil {
			HandleError(err)
		}
	}
	fmt.Printf("MultipartUploadSample Run Success !\n\n")

	// forbid overwrite
	uploadId, err = sc.CreateMultiPartUploadWithForbidOverwrite(bucketName, objectKey, s3.ObjectStorageClassStandard, true)
	if err == nil {
		fmt.Println("forbid overwrite success:", err)
	}

	// overwrite
	uploadId, err = sc.CreateMultiPartUploadWithForbidOverwrite(bucketName, objectKey, s3.ObjectStorageClassStandard,false)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("UploadId is: ", aws.StringValue(aws.String(uploadId)))
	partCount = 5
	completedUpload = &s3.CompletedMultipartUpload{
		Parts: make([]*s3.CompletedPart, partCount),
	}
	for i := 0; i < partCount; i++ {
		partNumber := int64(i + 1)
		etag, err := sc.UploadPart(bucketName, objectKey, GenMinimalPart(), uploadId, partNumber)
		if err != nil {
			HandleError(err)
		}
		completedUpload.Parts[i] = &s3.CompletedPart{
			ETag:       aws.String(etag),
			PartNumber: aws.Int64(partNumber),
		}
	}
	err = sc.CompleteMultiPartUpload(bucketName, objectKey, uploadId, completedUpload)
	if err != nil {
		HandleError(err)
		err = sc.AbortMultiPartUpload(bucketName, objectKey, uploadId)
		if err != nil {
			HandleError(err)
		}
	}

	fmt.Printf("overwrite Run Success !\n\n")
}
