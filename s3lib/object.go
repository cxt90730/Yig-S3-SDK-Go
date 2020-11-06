package s3lib

import (
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
	"io"
	"time"
	"fmt"
	"github.com/journeymidnight/yig/api/datatype"
	"encoding/base64"
	"mime/multipart"
	"bytes"
	"math/rand"
	"encoding/json"
	"net/http"
	"errors"
	"encoding/hex"
	"crypto/hmac"
	"crypto/sha256"
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

func (s3client *S3Client) PutObjectPresigned(bucketName, key string) (url string, err error) {
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
	req, _ := s3client.Client.PutObjectRequest(params);
	url, err = req.Presign(3600*time.Second)
	if err != nil {
		return "", err
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

func (s3client *S3Client) ListObjects(bucketName string, prefix string, delimiter string, maxKey int64) (
	keys []string, isTruncated bool, nextMarker string, err error) {
	params := &s3.ListObjectsInput{
		Bucket:    aws.String(bucketName),
		MaxKeys:   aws.Int64(maxKey),
		Delimiter: aws.String(delimiter),
		Prefix:    aws.String(prefix),
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

type postPolicyElem struct {
	Expiration string        `json:"expiration"`
	Conditions []interface{} `json:"conditions"`
}

func (s3Client *S3Client) newPostFormPolicy(expiration time.Time, conditions map[string]string, matches [][]string) (string, error) {
	var cons []interface{}
	for k, v := range conditions {
		m := make(map[string]string)
		m[k] = v
		cons = append(cons, m)
	}
	for _, v := range matches {
		cons = append(cons, v)
	}
	ppe := &postPolicyElem{
		Expiration: expiration.Format(time.RFC3339Nano),
		Conditions: cons,
	}

	body, err := json.Marshal(ppe)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(body)

	return encoded, nil
}

const (
	chars    = "0123456789_abcdefghijkl-mnopqrstuvwxyz" //ABCDEFGHIJKLMNOPQRSTUVWXYZ
	charsLen = len(chars)
	mask     = 1<<6 - 1
)

var rng = rand.NewSource(time.Now().UnixNano())

// RandBytes return the random byte sequence.
func RandBytes(ln int) []byte {
	/* chars 38 characters.
	 * rng.Int64() we can use 10 time since it produces 64-bit random digits and we use 6bit(2^6=64) each time.
	 */
	buf := make([]byte, ln)
	for idx, cache, remain := ln-1, rng.Int63(), 10; idx >= 0; {
		if remain == 0 {
			cache, remain = rng.Int63(), 10
		}
		buf[idx] = chars[int(cache&mask)%charsLen]
		cache >>= 6
		remain--
		idx--
	}
	return buf
}

func (s3client *S3Client) newPostFormBody(params map[string]string, fieldName, fileName string, fileSize int) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// set the params before file segment.
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, "", err
	}
	contents := RandBytes(fileSize)
	file := bytes.NewReader(contents)
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	return body, writer.FormDataContentType(), nil
}

// AWS Signature Version '4' constants.
const (
	signV4Algorithm = "AWS4-HMAC-SHA256"
)

// getSigningKey hmac seed to calculate final signature.
func getSigningKey(secretKey string, t time.Time, region string) []byte {
	date := sumHMAC([]byte("AWS4"+secretKey), []byte(t.Format(datatype.YYYYMMDD)))
	regionBytes := sumHMAC(date, []byte(region))
	service := sumHMAC(regionBytes, []byte("s3"))
	signingKey := sumHMAC(service, []byte("aws4_request"))
	return signingKey
}

// getSignature final signature in hexadecimal form.
func getSignature(signingKey []byte, stringToSign string) string {
	return hex.EncodeToString(sumHMAC(signingKey, []byte(stringToSign)))
}

// sumHMAC calculate hmac between two input byte array.
func sumHMAC(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

type PostObjectInput struct {
	Endpoint   string
	Bucket     string
	ObjName    string
	Expiration time.Time
	Date       time.Time
	Region     string
	AK         string
	SK         string
	FileSize   int
	Body       io.Reader
}

func (s3Client *S3Client) PostObject(pbi *PostObjectInput) error {
	commons := make(map[string]string)
	conditions := make(map[string]string)
	formParams := make(map[string]string)
	var matches [][]string
	// credential, pls refer to https://docs.amazonaws.cn/en_us/general/latest/gr/sigv4-create-string-to-sign.html
	cred := fmt.Sprintf("%s/%s/%s/s3/aws4_request", pbi.AK, pbi.Date.Format(datatype.YYYYMMDD), pbi.Region)
	encMethod := "AES256"

	commons["acl"] = "public-read"
	commons["x-amz-meta-uuid"] = "14365123651274"
	commons["x-amz-server-side-encryption"] = encMethod
	commons["x-amz-credential"] = cred
	commons["x-amz-algorithm"] = "AWS4-HMAC-SHA256"
	commons["x-amz-date"] = pbi.Date.Format(datatype.Iso8601Format)
	// Support custom fields
	commons["x:hehe"] = "hehe"

	conditions["bucket"] = pbi.Bucket
	for k, v := range commons {
		conditions[k] = v
		formParams[k] = v
	}
	// set key
	matches = append(matches, []string{"starts-with", "$key", pbi.ObjName})

	policyStr, err := s3Client.newPostFormPolicy(pbi.Expiration, conditions, matches)
	if err != nil {
		return err
	}
	t, _ := time.Parse(datatype.Iso8601Format, commons["x-amz-date"])
	signKey := getSigningKey(pbi.SK, t, pbi.Region)
	sig := getSignature(signKey, policyStr)

	formParams["Policy"] = policyStr
	formParams["x-amz-signature"] = sig
	formParams["key"] = pbi.ObjName

	body, contentType, err := s3Client.newPostFormBody(formParams, "file", pbi.ObjName, pbi.FileSize)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", pbi.Endpoint, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("post object for url %s returns %d", pbi.Endpoint, resp.StatusCode))
	}

	return nil
}
