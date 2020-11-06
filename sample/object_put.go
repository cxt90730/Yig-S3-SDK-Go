package sample

import (
	"fmt"
	"github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
	"os"
	"strings"
)

func PutObjectSample() {
	//DeleteTestBucketAndObject()
	//defer DeleteTestBucketAndObject()
	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	// Create a bucket
	//err := sc.MakeBucket(bucketName)
	//if err != nil {
	//	HandleError(err)
	//}

	// 1. Put a string object
	err := sc.PutObject(bucketName, objectKey, strings.NewReader("NewBucketAndObjectSample"))
	if err != nil {
		HandleError(err)
	}

	err = sc.DeleteObject(bucketName, objectKey)
	if err != nil {
		HandleError(err)
	}

	// 2. Put a file
	f, err := os.Open(localFilePath)
	defer f.Close()
	if err != nil {
		HandleError(err)
	}
	err = sc.PutObject(bucketName, objectKey, f)
	if err != nil {
		HandleError(err)
	}
	out, err := sc.GetObject(bucketName, objectKey)
	if err != nil {
		HandleError(err)
	}

	out.Close()

	fmt.Printf("NewObjectSample Run Success !\n\n")
}
