package sample

import (
	"github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
	"time"
	"fmt"
)

func PostObjectSample() {
	//DeleteTestBucketAndObject()
	//defer DeleteTestBucketAndObject()
	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	// Create a bucket
	//err := sc.MakeBucket(bucketName)
	//if err != nil {
	//	HandleError(err)
	//}

	pbi := &s3lib.PostObjectInput{
		Endpoint:   "http://" + bucketName + "." + endpoint,
		Bucket:     bucketName,
		ObjName:    objectKey,
		Expiration: time.Now().UTC().Add(time.Duration(1 * time.Hour)),
		Date:       time.Now().UTC(),
		Region:     "r",
		AK:         accessKey,
		SK:         secretKey,
		FileSize:   1024,
	}

	// 1. Put a string object
	err := sc.PostObject(pbi)
	if err != nil {
		HandleError(err)
	}

	err = sc.DeleteObject(bucketName, objectKey)
	if err != nil {
		HandleError(err)
	}

	fmt.Printf("PostObjectSample Run Success !\n\n")
}
