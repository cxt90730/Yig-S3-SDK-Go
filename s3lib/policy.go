package s3lib

import (
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
	"strings"
	"github.com/journeymidnight/yig/api/datatype/policy"
	"github.com/journeymidnight/yig/api/datatype/policy/condition"
	"fmt"
	"bytes"
	"encoding/json"
)

func newPolicyWithStatement(s ...policy.Statement) (p *policy.Policy) {
	return &policy.Policy{
		Version:    "2012-10-17",
		Statements: s,
	}
}

func makeRefererPolicy(currentPolicy string, bucketName string, referers []string) ([]byte, error) {
	if currentPolicy == "" {
		s, err := newRefererStatement(bucketName, referers)
		if err != nil {
			return nil, err
		}
		return newPolicyWithStatement(*s).MarshalJSON()
	}

	p, err := policy.ParseConfig(strings.NewReader(currentPolicy), bucketName)
	if err != nil {
		return nil, err
	}
	if p.Statements == nil {
		s, err := newRefererStatement(bucketName, referers)
		if err != nil {
			return nil, err
		}
		return newPolicyWithStatement(*s).MarshalJSON()
	}

	for i, s := range p.Statements {
		if s.SID == "Referer" {
			if len(referers) == 0 {
				p.Statements = append(p.Statements[:i], p.Statements[i+1:]...)
				return p.MarshalJSON()

			} else {
				cond, err := condition.NewStringLikeFunc(condition.AWSReferer, referers...)
				if err != nil {
					return nil, err
				}
				s.Conditions = condition.Functions{cond}
				p.Statements[i] = s
				return p.MarshalJSON()
			}
		}
	}

	s, err := newRefererStatement(bucketName, referers)
	if err != nil {
		return nil, err
	}
	p.Statements = append(p.Statements, *s)
	return p.MarshalJSON()
}

func newRefererStatement(bucketName string, referers []string) (*policy.Statement, error) {
	var conds condition.Functions
	if len(referers) == 0 {
		conds = nil
	} else {
		cond, err := condition.NewStringLikeFunc(condition.AWSReferer, referers...)
		if err != nil {
			return nil, err
		}
		conds = condition.Functions{cond}
	}

	s := policy.NewStatement(
		policy.Allow,
		policy.NewPrincipal("*"),
		policy.NewActionSet(policy.GetObjectAction),
		policy.NewResourceSet(
			policy.NewResource(bucketName, "*"),
		),
		conds,
	)
	s.SID = "Referer"
	return &s, nil
}

func (s3client *S3Client) PutReferer(bucketName string, referers []string) (err error) {
	currentPolicy, err := s3client.GetBucketPolicy(bucketName)
	if err != nil {
		return
	}
	refJson, err := makeRefererPolicy(currentPolicy, bucketName, referers)
	if err != nil {
		return
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, refJson, "", "    ")
	if err != nil {
		return
	}
	refJson2 := buf.Bytes()
	fmt.Println(len(refJson))
	fmt.Println(string(refJson2))
	return s3client.PutBucketPolicy(bucketName, string(refJson))
}

func (s3client *S3Client) GetReferer(bucketName string) (referers []string, err error) {
	currentPolicy, err := s3client.GetBucketPolicy(bucketName)
	if err != nil {
		return
	}
	if currentPolicy == "" {
		return
	}
	p, err := policy.ParseConfig(strings.NewReader(currentPolicy), bucketName)
	if err != nil {
		return
	}
	if p.Statements == nil {
		return
	}
	for _, s := range p.Statements {
		if s.SID == "Referer" {
			var referers []string
			for _,r := range s.Conditions {
				referers = append(referers, strings.Replace(r.String(), "StringLike:aws:Referer:","", -1)[1:len(r.String())-24])
			}
			return referers, nil
		}
	}
	return
}

func (s3client *S3Client) PutBucketPolicy(bucketName, policy string) (err error) {
	params := &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(policy),
	}
	_, err = s3client.Client.PutBucketPolicy(params)
	if err != nil {
		return
	}
	return
}

func (s3client *S3Client) GetBucketPolicy(bucketName string) (policy string, err error) {
	params := &s3.GetBucketPolicyInput{
		Bucket: aws.String(bucketName),
	}
	out, err := s3client.Client.GetBucketPolicy(params)
	if err != nil {
		return "", err
	}
	return *out.Policy, err
}

func (s3client *S3Client) DeleteBucketPolicy(bucketName string) (err error) {
	params := &s3.DeleteBucketPolicyInput{
		Bucket: aws.String(bucketName),
	}
	_, err = s3client.Client.DeleteBucketPolicy(params)
	if err != nil {
		return
	}
	return
}
