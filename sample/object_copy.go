package sample

import (
	"fmt"
	"github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
	"io/ioutil"
	"strings"
)

func CopyObjectSample() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()

        sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	// Create a sourceBucket and descBucket
	var descBucketName ="descbucketname"
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	err = sc.MakeBucket(descBucketName)
	if err != nil {
		HandleError(err)
	}

	// 1. Put a string object
	err = sc.PutObject(bucketName, objectKey, strings.NewReader("NewBucketAndObjectSample"))
	if err != nil {
                fmt.Println("111")
		HandleError(err)
	}

	// 2: Copy an existing object
	var descObjectKey = "descobject"
	//var copySource="/"+bucketName+"/"+objectKey
        err = sc.CopyObject(descBucketName,"/go-sdk-test/go-sdk-key",descObjectKey)
	if err != nil {
                fmt.Println("222")
		HandleError(err)
	}

    // 3. Get copy bucket object
	out,err := sc.GetObject(descBucketName,descObjectKey)
	if err != nil {
		HandleError(err)
	}
	b, _ := ioutil.ReadAll(out)
	fmt.Println("Get appended string:", string(b))
	out.Close()
     
        err = sc.DeleteObject(descBucketName, descObjectKey)
	if err != nil {
		HandleError(err)
	}

	err = sc.DeleteBucket(descBucketName)
	if err != nil {
		HandleError(err)
	}
	fmt.Printf("CopyObjectSample Run Success !\n\n")
}
