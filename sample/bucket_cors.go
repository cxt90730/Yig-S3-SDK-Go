package sample

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
)

func BucketCORSSample() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()
	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	rule1 := &s3.CORSRule{
		AllowedOrigins: []*string{aws.String("*")},
		AllowedMethods: []*string{aws.String("PUT"), aws.String("GET"), aws.String("POST")},
		AllowedHeaders: []*string{},
		ExposeHeaders:  []*string{},
		MaxAgeSeconds:  aws.Int64(100),
	}

	rule2 := &s3.CORSRule{
		AllowedOrigins: []*string{aws.String("http://www.a.com"), aws.String("http://www.b.com")},
		AllowedMethods: []*string{aws.String("GET")},
		AllowedHeaders: []*string{aws.String("Authorization")},
		ExposeHeaders:  []*string{aws.String("x-amz-test"), aws.String("x-amz-test1")},
		MaxAgeSeconds:  aws.Int64(100),
	}

	err = sc.SetBucketCORS(bucketName, []*s3.CORSRule{rule1})
	if err != nil {
		HandleError(err)
	}

	out, err := sc.GetBucketCORS(bucketName)
	if err != nil {
		HandleError(err)
	}
	fmt.Println(out)

	err = sc.SetBucketCORS(bucketName, []*s3.CORSRule{rule1, rule2})
	if err != nil {
		HandleError(err)
	}

	out, err = sc.GetBucketCORS(bucketName)
	if err != nil {
		HandleError(err)
	}
	fmt.Println(out)

	err = sc.DeleteBucketCORS(bucketName)
	if err != nil {
		HandleError(err)
	}

	out, err = sc.GetBucketCORS(bucketName)
	if err != nil {
		HandleError(err)
	}
	fmt.Println(out)

	fmt.Printf("BucketCORSSample Run Success !\n\n")
}
