package miniodb

import (
	"github.com/minio/minio-go"
	"main/constants"
)

var MinioDB *minio.Client

func Init() {
	var err error

	MinioDB, err = minio.New(constants.MinioEndpoint, constants.MinioAccessKeyID, constants.MinioSecretAccessKey, constants.MinioUseSSL)

	if err != nil {
		panic(err)
	}

}
