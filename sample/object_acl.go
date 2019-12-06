package sample

import (
	"fmt"
	"github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
	"strings"
)

func ObjectACLSample() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()
	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	// Create a bucket
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	// Test ObjectACL public-read
	err = sc.PutObject(bucketName, objectKey, strings.NewReader("NewBucketAndObjectSample"))
	if err != nil {
		HandleError(err)
	}
	err = sc.PutObjectAcl(bucketName, objectKey, "public-read")
	if err != nil {
		HandleError(err)
	}
	out, err := sc.GetObjectAcl(bucketName, objectKey)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("Get Bucket ACL:", out)
	err = sc.DeleteObject(bucketName, objectKey)
	if err != nil {
		HandleError(err)
	}

	// Test ObjectACL public-read-write
	err = sc.PutObject(bucketName, objectKey, strings.NewReader("NewBucketAndObjectSample"))
	if err != nil {
		HandleError(err)
	}
	err = sc.PutObjectAcl(bucketName, objectKey, "public-read-write")
	if err != nil {
		HandleError(err)
	}
	out, err = sc.GetObjectAcl(bucketName, objectKey)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("Get Bucket ACL:", out)
	err = sc.DeleteObject(bucketName, objectKey)
	if err != nil {
		HandleError(err)
	}

	fmt.Printf("ObjectACLSample Run Success!\n\n")
}
