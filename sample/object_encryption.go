package sample

import (
	"fmt"
	"github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
)

func PutEncryptObjectWithSSECSample() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()
	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	err = sc.PutEncryptObjectWithSSEC(bucketName, objectKey, "PutEncryptObjectWithSSEC")
	if err != nil {
		HandleError(err)
	}

	v, err := sc.GetEncryptObjectWithSSEC(bucketName, objectKey)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("GetEncryptObjectWithSSEC:",v)

	fmt.Printf("PutEncryptObjectWithSSECSample Run Success !\n\n")
}

func PutEncryptObjectWithSSES3Sample() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()
	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	// Create a bucket
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	err = sc.PutEncryptObjectWithSSES3(bucketName, objectKey, "PutEncryptObjectWithSSES3")
	if err != nil {
		HandleError(err)
	}

	v, err := sc.GetEncryptObjectWithSSES3(bucketName, objectKey)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("GetEncryptObjectWithSSES3:",v)

	fmt.Printf("PutEncryptObjectWithSSES3ample Run Success !\n\n")
}
