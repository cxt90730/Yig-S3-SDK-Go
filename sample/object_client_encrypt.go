package sample

import (
	"bytes"
	"fmt"
	"github.com/journeymidnight/Yig-S3-SDK-Go/s3lib"
	"io/ioutil"
	"strings"
)

func ClientEncrypt(){
	MasterKey := []byte{1,2,3,4}
	// please generate your own encrypted 32bytes keys
	// just mock implement
	mockEncryptedKey :=  bytes.Repeat(MasterKey, 8) // 4bytes * 8 = 32 bytes

	ec := s3lib.NewEncryptS3Client(endpoint, accessKey, secretKey, mockEncryptedKey)
	err := ec.PutObjectWithSelfEncrypt(bucketName, objectKey, strings.NewReader("test_value"))
	if err != nil {
		panic(err)
	}

	sc := s3lib.NewS3(endpoint, accessKey, secretKey)
	out1, err := sc.GetObject(bucketName, objectKey)
	if err != nil {
		panic(err)
	}
	defer out1.Close()
	data1, err := ioutil.ReadAll(out1)
	if err != nil {
		panic(err)
	}
	fmt.Println("data1:", string(data1))

	dc := s3lib.NewDecryptS3Client(endpoint, accessKey, secretKey)
	out, err := dc.GetObjectWithSelfDecrypt(bucketName, objectKey)
	if err != nil {
		panic(err)
	}
	defer out.Body.Close()
	data, err := ioutil.ReadAll(out.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("data:", string(data))
}
