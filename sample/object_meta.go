package sample

import (
    "fmt"
    "github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
    "strings"
)

func ObjectMetaSample() {
	DeleteTestBucketAndObject()

	defer DeleteTestBucketAndObject()

	// Set Custom Meta
        sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	// Create a bucket
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}
	c := make(map[string]string)
	c["a"]="b"
	err = sc.PutObjectMeta(bucketName, objectKey, strings.NewReader("NewBucketAndObjectSample"),c)
	if err != nil {
		HandleError(err)
	}

	fmt.Printf("ObjectMetaSample Run Success !\n\n")
}
