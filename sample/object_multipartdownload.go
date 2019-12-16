package sample

import (
	"bytes"
	"fmt"
	"github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
)

func MultiPartDownloadSample() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()
	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	// Create a bucket
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	// Put a 5kib file
	RANGE := make([]byte, 5<<10)
	err = sc.PutObject(bucketName, objectKey, bytes.NewReader(RANGE))
	if err != nil {
		HandleError(err)
	}
	//Slice download by range
	ranges := map[string]string{"0": "1000", "1001": "2000", "2001": "3000", "3001": "4000", "4001": "5119"}
	for range1, range2 := range ranges {
		out, err := sc.GetObjectWithRange(bucketName, objectKey, "bytes="+range1+"-"+range2)
		if err != nil {
			HandleError(err)
		}
		fmt.Println("Download range is :", out)

	}
	fmt.Printf("MultiPartDownloadSample Run Success !\n\n")
}
