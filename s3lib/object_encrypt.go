package s3lib

import (
	"crypto/rand"
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
	"github.com/journeymidnight/aws-sdk-go/service/s3/s3crypto"
	"io"
)

const UnisWrap = "unis"

type CustomKeyHandler struct {
	s3crypto.CipherData
	EncryptedKey []byte
}

func (ckh *CustomKeyHandler) DecryptKey(key []byte) ([]byte, error) {
	return key, nil
}

func (ckh *CustomKeyHandler) GenerateCipherData(keySize, ivSize int) (data s3crypto.CipherData, err error) {
	iv := generateBytes(ivSize)
	data = s3crypto.CipherData{
		Key:           ckh.EncryptedKey,
		IV:            iv,
		WrapAlgorithm: UnisWrap,
		EncryptedKey:  ckh.EncryptedKey,
	}
	return
}

func generateBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}


func (ckh CustomKeyHandler) decryptHandler(env s3crypto.Envelope) (s3crypto.CipherDataDecrypter, error) {
	ckh.WrapAlgorithm = UnisWrap
	return &ckh, nil
}

func (ec *EncryptionClient) PutObjectWithSelfEncrypt(bucketName, key string, value io.Reader) error {
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   aws.ReadSeekCloser(value),
		ACL:    aws.String(ObjectCannedACLPublicRead),
	}
	_, err := ec.Client.PutObject(params)
	return err
}

func (dc *DecryptionClient) GetObjectWithSelfDecrypt(bucketName, key string) (out *s3.GetObjectOutput, err error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
	return dc.Client.GetObject(params)
}