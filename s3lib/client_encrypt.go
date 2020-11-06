package s3lib

import (
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/aws/credentials"
	"github.com/journeymidnight/aws-sdk-go/aws/session"
	"github.com/journeymidnight/aws-sdk-go/service/s3/s3crypto"
)

type EncryptionClient struct {
	Client *s3crypto.EncryptionClient
}

func NewEncryptS3Client(endpoint, accessKey, secretKey string, encryptKey []byte) *EncryptionClient {
	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	sess := session.Must(session.NewSession(
		&aws.Config{
			Credentials: creds,
			DisableSSL:  aws.Bool(true),
			Endpoint:    aws.String(endpoint),
			Region:      aws.String("ap-northeast-1"),
		},
	),
	)

	generator := &CustomKeyHandler{EncryptedKey: encryptKey}
	ec := s3crypto.NewEncryptionClient(sess, s3crypto.AESGCMContentCipherBuilder(generator))
	return &EncryptionClient{ec}
}

type DecryptionClient struct {
	Client *s3crypto.DecryptionClient
}

func NewDecryptS3Client(endpoint, accessKey, secretKey string) *DecryptionClient {
	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	sess := session.Must(session.NewSession(
		&aws.Config{
			Credentials: creds,
			DisableSSL:  aws.Bool(true),
			Endpoint:    aws.String(endpoint),
			Region:      aws.String("ap-northeast-1"),
		},
	),
	)
	dc := s3crypto.NewDecryptionClient(sess,func(svc *s3crypto.DecryptionClient){
		svc.LoadStrategy = s3crypto.HeaderV2LoadStrategy{}
		svc.WrapRegistry[UnisWrap] = CustomKeyHandler{}.decryptHandler
	})
	return &DecryptionClient{dc}
}