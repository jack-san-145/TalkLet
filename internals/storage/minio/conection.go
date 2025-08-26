package minio

import (
	"context"
	"fmt"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	Minio_client *minio.Client
	err          error
)

func CreateMinioClient() {
	endpoint := "localhost:9000"
	username := os.Getenv("MINIO_CONTAINER_ROOT_USER")
	password := os.Getenv("MINIO_CONTAINER_ROOT_PASSWORD")
	isSLL := false
	Minio_client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(username, password, ""),
		Secure: isSLL,
	})
	if err != nil {
		fmt.Println("error while creating the minio client - ", err)
		return
	}

	_, err = Minio_client.ListBuckets(context.Background())
	if err != nil {
		fmt.Println("minio server is unreachable - ", err)
		return
	}
	fmt.Println("Minio conencted Successfully")
}
