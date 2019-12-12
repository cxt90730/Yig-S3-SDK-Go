package sample

import (
	"fmt"
	"github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
	"github.com/journeymidnight/aws-sdk-go/aws"
	"strings"
)

func MultipartUploadSample() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()
	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	// Create a bucket
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	// 1.Create Multipart Upload
	uploadId, err := sc.CreateMultipartUpload(bucketName, objectKey)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("UploadId is: ", aws.StringValue(uploadId))

	// 2.Upload Part
	ETag, err := sc.UploadPart(bucketName, objectKey, uploadId, strings.NewReader("NewBucketAndObjectSample"))
	if err != nil {
		HandleError(err)
	}
	fmt.Println("ETag is: ", aws.StringValue(ETag))

	// 3.CompleteMultipartUpload
	err = sc.CompleteMultipartUpload(bucketName, objectKey, ETag, uploadId)
	if err != nil {
		HandleError(err)
	}

	fmt.Printf("MultipartUploadSample Run Success !\n\n")
}
