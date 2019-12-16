package sample

import (
	"fmt"
	"github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
)

func MultiPartUploadSample() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()
	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
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
	fmt.Println("UploadId is: ", aws.StringValue(uploadId))
	partCount := 3
	completedUpload := &s3.CompletedMultipartUpload{
		Parts: make([]*s3.CompletedPart, partCount),
	}
	for i := 0; i < partCount; i++ {
		partNumber := int64(i + 1)
		etag, err := sc.UploadPart(bucketName, objectKey, partNumber, uploadId, GenMinimalPart())
		if err != nil {
			HandleError(err)
		}
		completedUpload.Parts[i] = &s3.CompletedPart{
			ETag:       aws.String(etag),
			PartNumber: aws.Int64(partNumber),
		}
	}
	// 2.Upload Part
	err = sc.CompleteMultiPartUpload(bucketName, objectKey, completedUpload, uploadId)
	if err != nil {
		HandleError(err)
		err = sc.AbortMultiPartUpload(bucketName, objectKey, uploadId)
		if err != nil {
			HandleError(err)
		}
	}
	fmt.Printf("MultipartUploadSample Run Success !\n\n")

}

// Generate 5M part data
func GenMinimalPart() []byte {
	return make([]byte, 5<<20)
}
