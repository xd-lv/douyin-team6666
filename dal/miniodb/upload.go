package miniodb

import (
	"github.com/minio/minio-go"
	"io"
	"net/url"
	"time"
)

func Upload(file io.Reader, bukectName string, objectName string, fileSize int64) (string, error) {
	var URL *url.URL
	var res string
	_, err := MinioDB.PutObject(bukectName, objectName, file, fileSize, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return res, err
	}

	URL, err = MinioDB.PresignedGetObject(bukectName, objectName, time.Second*24*60*60, make(url.Values))
	if err != nil {
		return res, err
	}
	res = URL.String()
	return res, nil
}

func Create(bucketName string) error {
	err := MinioDB.MakeBucket(bucketName, "us-east-1")
	//MinioDB.SetBucketPolicy()

	if err != nil {
		return err
	}
	return nil
}

func Test(bucketName string) error {
	policy := "{\n  \"Statement\": [\n    {\n      \"Action\": [\n        \"s3:GetBucketLocation\",\n        \"s3:ListBucket\"\n      ],\n      \"Effect\": \"Allow\",\n      \"Principal\": \"*\",\n      \"Resource\": \"arn:aws:s3:::my-bucketname\"\n    },\n    {\n      \"Action\": \"s3:GetObject\",\n      \"Effect\": \"Allow\",\n      \"Principal\": \"*\",\n      \"Resource\": \"arn:aws:s3:::my-bucketname/myobject*\"\n    }\n  ],\n  \"Version\": \"2012-10-17\"\n}\n\n"
	err := MinioDB.SetBucketPolicy(bucketName, policy)
	return err
}

func IsExists(bucketName string) bool {
	exists, _ := MinioDB.BucketExists(bucketName)
	return exists
}
