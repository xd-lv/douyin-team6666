package miniodb

import (
	"github.com/minio/minio-go"
	"main/constants"

	"mime/multipart"
)

func Upload(file multipart.File, bukectName string, objectName string, fileSize int64) (string, error) {
	var url string
	_, err := MinioDB.PutObject(bukectName, objectName, file, fileSize, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return url, err
	}
	url = "http://" + constants.MinioEndpoint + "/" + bukectName + "/" + objectName
	return url, nil
}

func Create(bucketName string) error {
	err := MinioDB.MakeBucket(bucketName, "us-east-1")
	if err != nil {
		return err
	}
	return nil
}

func IsExists(bucketName string) bool {
	exists, _ := MinioDB.BucketExists(bucketName)
	return exists
}
