package sample

import (
	"fmt"
	"github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
)


func PutObjectEncryptSample() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()
	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	// Create a bucket
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	// 1. Put Encrypt Object With SSEC
	err = sc.PutEncryptObjectWithSSEC(bucketName, objectKey, "NewBucketAndObjectSample")
	if err != nil {
		HandleError(err)
	}

	result, err := sc.GetEncryptObjectWithSSEC(bucketName, objectKey)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("Get Encrypt Object With SSEC: ", result)

	fmt.Printf("PutObjectEncryptSample Run Success !\n\n")
}
