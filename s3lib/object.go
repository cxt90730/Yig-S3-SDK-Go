package s3lib

import (
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
	"io"
	"time"
	"fmt"
)

func (s3client *S3Client) PutObject(bucketName, key string, value io.Reader) (err error) {
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   aws.ReadSeekCloser(value),
	}
	if _, err = s3client.Client.PutObject(params); err != nil {
		return err
	}
	return
}

func (s3client *S3Client) PutObjectMeta(bucketName, key string, value io.Reader, meta map[string]string) (err error) {
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   aws.ReadSeekCloser(value),
		Metadata:  aws.StringMap(meta),
	}
	if _, err = s3client.Client.PutObject(params); err != nil {
		return err
	}
	return
}

func (s3client *S3Client) PutObjectPreSignedWithSpecifiedBody(bucketName, key string, value io.Reader, expire time.Duration) (url string, err error) {
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   aws.ReadSeekCloser(value),
	}
	req, _ := s3client.Client.PutObjectRequest(params)
	return req.Presign(expire)
}

func (s3client *S3Client) PutObjectPreSignedWithoutSpecifiedBody(bucketName, key string, expire time.Duration) (url string, err error) {
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
	req, _ := s3client.Client.PutObjectRequest(params)
	return req.Presign(expire)
}

func (s3client *S3Client) HeadObject(bucketName, key string) (err error) {
	params := &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
	_, err = s3client.Client.HeadObject(params)
	if err != nil {
		return err
	}
	return
}

func (s3client *S3Client) GetObject(bucketName, key string) (value io.ReadCloser, err error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
	out, err := s3client.Client.GetObject(params)
	if err != nil {
		return nil, err
	}
	return out.Body, err
}

func (s3client *S3Client) GetObjectPreSigned(bucketName, key string, expire time.Duration) (url string, err error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
	req, _ := s3client.Client.GetObjectRequest(params)
	return req.Presign(expire)
}

func (s3client *S3Client) DeleteObject(bucketName, key string) (err error) {
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
	_, err = s3client.Client.DeleteObject(params)
	if err != nil {
		return err
	}
	return
}

func (s3client *S3Client) DeleteObjects(bucketName string, key ...string) (deletedKeys []string, err error) {
	var objects []*s3.ObjectIdentifier
	for _, k := range key {
		oi := &s3.ObjectIdentifier{
			Key: aws.String(k),
		}
		objects = append(objects, oi)
	}

	params := &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &s3.Delete{
			Objects: objects,
		},
	}
	out, err := s3client.Client.DeleteObjects(params)
	if err != nil {
		return nil, err
	}
	for _, dk := range out.Deleted {
		deletedKeys = append(deletedKeys, *dk.Key)
	}
	return
}

func (s3client *S3Client)CopyObject(destinationBucket,copySource, key string) (err error) {
	params := &s3.CopyObjectInput{
		Bucket: aws.String(destinationBucket),
		CopySource: aws.String(copySource),
		Key:    aws.String(key),
	}
	_, err = s3client.Client.CopyObject(params)
	if err != nil {
		return err
	}
	return
}

func (s3client *S3Client) AppendObject(bucketName, key string, value io.ReadSeeker, position int64) (nextPos int64, err error) {
	var out *s3.AppendObjectOutput
	params := &s3.AppendObjectInput{
		Bucket:   aws.String(bucketName),
		Key:      aws.String(key),
		Body:     value,
		Position: aws.Int64(position),
	}
	if out, err = s3client.Client.AppendObject(params); err != nil {
		return 0, err
	}

	return *out.NextPosition, nil
}

func (s3client *S3Client) AppendObjectWithAclAndMeta(bucketName, key string, value io.ReadSeeker, position int64, acl string, meta map[string]string) (nextPos int64, err error) {
	var out *s3.AppendObjectOutput
	params := &s3.AppendObjectInput{
		Bucket:   aws.String(bucketName),
		Key:      aws.String(key),
		Body:     value,
		Position: aws.Int64(position),
		ACL:      aws.String(acl),
		Metadata:  aws.StringMap(meta),
	}
	if out, err = s3client.Client.AppendObject(params); err != nil {
		return 0, err
	}

	return *out.NextPosition, nil
}

func (s3client *S3Client) GetObjectNextAppendPosition(bucketName, key string, position int64) (nextPos int64, err error) {
	var out *s3.AppendObjectOutput
	params := &s3.AppendObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Position: aws.Int64(position),
	}
	if out, err = s3client.Client.AppendObject(params); err != nil {
		return 0, err
	}

	return *out.NextPosition, nil
}

func (s3client *S3Client) ListObjects(bucketName string, prefix string, delimiter string, maxKey int64)(
	keys []string, isTruncated bool, nextMarker string, err error) {
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
		MaxKeys: aws.Int64(maxKey),
		Delimiter: aws.String(delimiter),
		Prefix: aws.String(prefix),
	}
	out, err := s3client.Client.ListObjects(params)
	if err != nil {
		return
	}
	isTruncated = *out.IsTruncated
	if out.NextMarker != nil {
		nextMarker = *out.NextMarker
	}
	for _, v := range out.CommonPrefixes {
		keys = append(keys, *v.Prefix)
		fmt.Println("Prefix:", *v.Prefix)
	}
	for _, v := range out.Contents {
		keys = append(keys, *v.Key)
		fmt.Println("Key:", *v.Key)
	}

	return
}
