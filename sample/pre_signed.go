package sample

import (
	"fmt"
	"github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
	"strings"
	"time"
)

func PreSignedSample() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()

	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	// Create a bucket
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	// 1. Put object pre signed url with specified body
	url, err := sc.PutObjectPreSignedWithSpecifiedBody(bucketName, objectKey, strings.NewReader("NewBucketAndObjectSample"), (5 * time.Minute))
	if err != nil {
		HandleError(err)
	}
	fmt.Println("put object pre signed with specified body url is :", url)
	// 2. Get object pre signed url
	url, err = sc.GetObjectPreSigned(bucketName, objectKey, (5 * time.Minute))
	if err != nil {
		HandleError(err)
	}
	fmt.Println("get object pre signed url is :", url)

	err = sc.DeleteObject(bucketName, objectKey)
	if err != nil {
		HandleError(err)
	}

	//3. Put object pre signed url without specified body
	url, err = sc.PutObjectPreSignedWithoutSpecifiedBody(bucketName, objectKey, (5 * time.Minute))
	if err != nil {
		HandleError(err)
	}
	fmt.Println("put object pre signed without specified body url is :", url)

	fmt.Printf("PreSignedSample Run Success!\n\n")
}
