package sample

import (
	"fmt"
	"github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
)

//Not Implemented
func BucketRequestPaymentSample() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()

	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	//put bucket website
	err = sc.PutBucketRequestPayment(bucketName)
	if err != nil {
		HandleError(err)
	}

	//Get Bucket Website
	result, err := sc.GetBucketRequestPayment(bucketName)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("Request Payment is: ", result)

	fmt.Printf("BucketRequestPaymentSample Run Success !\n\n")
}
