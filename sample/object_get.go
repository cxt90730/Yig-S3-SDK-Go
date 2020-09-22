package sample

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
)

func GetObjectSample() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()
	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	// Create a bucket
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	// Put a file
	f, err := os.Open(localFilePath)
	defer f.Close()
	if err != nil {
		HandleError(err)
	}
	err = sc.PutObject(bucketName, objectKey, f)
	if err != nil {
		HandleError(err)
	}

	// Get the reader
	out, err := sc.GetObject(bucketName, objectKey)
	if err != nil {
		HandleError(err)
	}

	// Download to a file
	f2, err := os.OpenFile("sample/Download.jpeg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	defer f2.Close()
	if err != nil {
		HandleError(err)
	}
	io.Copy(f2, out)
	out.Close()

	fmt.Printf("GetObjectSample Run Success !\n\n")
}

func GetObjectByRange() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()
	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	object := "abcdefghi#jklmnopqIsMyWantedPart!!!r$stuvwxyz1@234567890-"
	err = sc.PutObject(bucketName, objectKey, strings.NewReader(object))
	if err != nil {
		HandleError(err)
	}

	var getObject *s3.GetObjectOutput
	// Get object part
	rangeString := "bytes=18-34"
	getObject, err = sc.GetObjectWithRange(bucketName, objectKey, rangeString)
	if err != nil {
		HandleError(err)
	}

	b, _ := ioutil.ReadAll(getObject.Body)
	value := string(b)
	if value != "IsMyWantedPart!!!" {
		fmt.Println("GetObject err: value is:", getObject, ",\nbut should be:", object)
	}
	fmt.Printf("GetObjectByRange Run Success !\n\n")}

func GetObjectByRange2() {
	DeleteTestBucketAndObject()
	defer DeleteTestBucketAndObject()
	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	// Create a bucket
	err := sc.MakeBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	f, err := os.Open(localFilePath)
	if err != nil {
		HandleError(err)
	}
	defer f.Close()

	err = sc.PutObject(bucketName, objectKey, f)
	if err != nil {
		HandleError(err)
	}

	out, err := sc.HeadObject(bucketName, objectKey)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("Object info: ", out)

	var partRange int64
	var num int
	var failedDownload []string
	objects := make(map[int]*s3.GetObjectOutput)
	objectSize := aws.Int64Value(out.ContentLength)
	// Get object by parts
	for {
		num++
		if partRange+10<<10-1 <= objectSize {
			rangeString := "bytes=" + strconv.FormatInt(partRange, 10) + "-" + strconv.FormatInt(partRange+(10<<10-1), 10)
			object, err := sc.GetObjectWithRange(bucketName, objectKey, rangeString)
			if err != nil {
				fmt.Println("part ", num, " download err: ", err, "\n rangeString: ", rangeString)
				failedDownload = append(failedDownload, rangeString)
				continue
			}

			fmt.Println("part ", num, " object info: ", object)
			objects[num] = object
		} else if partRange < objectSize {
			rangeString := "bytes=" + strconv.FormatInt(partRange, 10) + "-" + strconv.FormatInt(objectSize-1, 10)
			fmt.Println("rangeString: ", rangeString)
			object, err := sc.GetObjectWithRange(bucketName, objectKey, rangeString)
			if err != nil {
				fmt.Println("part ", num, " download err: ", err, "\n rangeString: ", rangeString)
				failedDownload = append(failedDownload, rangeString)
				continue
			}
			fmt.Println("part ", num, " object info: ", object)
			objects[num] = object
			break
		} else {
			// when range > object size, s3 wil return object
			break
		}
		partRange += 10 << 10
	}
	// get parts which download failed before
	for i := 0; i < len(failedDownload); i++ {
		object, err := sc.GetObjectWithRange(bucketName, objectKey, failedDownload[i])
		if err != nil {
			fmt.Println("part ", num, " download err: ", err, "\n rangeString: ", failedDownload[i])
			failedDownload = append(failedDownload, failedDownload[i])
			continue
		}
		fmt.Println("part ", num, " object info: ", object)
		objects[num] = object
	}

	f2, err := os.OpenFile("sample/Download.jpg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		HandleError(err)
	}
	defer f2.Close()

	for i := 1; i <= len(objects); i++ {
		io.Copy(f2, objects[i].Body)
		objects[i].Body.Close()
	}

	fmt.Printf("GetObjectByRange2 Run Success !\n\n")
}

