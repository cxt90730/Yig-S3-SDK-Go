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
	var descBucketName = "descbucket"
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
		HandleError(err)
	}

	// 2: Copy an existing object
	var descObjectKey = "descObject"
	//var copySource="/"+bucketName+"/"+objectKey
	err = sc.CopyObject(descBucketName, bucketName+"/"+objectKey, descObjectKey)
	if err != nil {
		HandleError(err)
	}

	// 3. Get copy bucket object
	out, err := sc.GetObject(descBucketName, descObjectKey)
	if err != nil {
		HandleError(err)
	}
	b, _ := ioutil.ReadAll(out)
	fmt.Println("Get string:", string(b))
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

func CopyObjectWithForbidOverwriteSample() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()

	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	// Create a sourceBucket and descBucket
	var descBucketName = "descbucket"
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	err = sc.MakeBucket(descBucketName)
	if err != nil {
		HandleError(err)
	}

	err = sc.PutObject(bucketName, objectKey, strings.NewReader("NewBucketAndObjectSample"))
	if err != nil {
		HandleError(err)
	}

	var descObjectKey = "descObject"
	//var copySource="/"+bucketName+"/"+objectKey
	err = sc.CopyObject(descBucketName, bucketName+"/"+objectKey, descObjectKey)
	if err != nil {
		HandleError(err)
	}

	output, err := sc.CopyObjectWithForbidOverwrite(descBucketName, bucketName+"/"+objectKey, descObjectKey, true)
	if err == nil {
		fmt.Println("output:",output)
		HandleError(err)
	}
	fmt.Println("Forbid overwrite Success!", err)

	output, err = sc.CopyObjectWithForbidOverwrite(descBucketName, bucketName+"/"+objectKey, descObjectKey, false)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("PutObjectWithOverwrite Success!", output)

	out, err := sc.GetObject(descBucketName, descObjectKey)
	if err != nil {
		HandleError(err)
	}
	b, _ := ioutil.ReadAll(out)
	fmt.Println("Get string:", string(b))
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
