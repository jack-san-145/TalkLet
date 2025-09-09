package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"os"
)

var (
	Minio_client *minio.Client
	err          error
)

var Bucket_name = "talklet-media"

func CreateMinioClient() {
	endpoint := "localhost:9000"

	// endpoint := "10.10.5.153:9000"
	// endpoint := "192.168.170.106:9000"
	// endpoint := "127.0.0.0:9000"   //localhost:9000
	username := os.Getenv("MINIO_CONTAINER_ROOT_USER")      //root username
	password := os.Getenv("MINIO_CONTAINER_ROOT_PASSWORD")  //root password
	isSLL := false                                          //http or https
	Minio_client, err = minio.New(endpoint, &minio.Options{ //creating the new minio client
		Creds:  credentials.NewStaticV4(username, password, ""), // aws style authentication
		Secure: isSLL,                                           //decides the api either http or https
	})
	if err != nil {
		fmt.Println("error while creating the minio client - ", err)
		return
	}

	_, err = Minio_client.ListBuckets(context.Background()) //by make the single minio cmd we can see whether is pinging or not
	if err != nil {
		fmt.Println("minio server is unreachable - ", err)
		return
	}
	fmt.Println("Minio conencted Successfully")
	createMinioBucket()
}

func createMinioBucket() {
	exists, _ := Minio_client.BucketExists(context.Background(), Bucket_name) //checks if the bucket exist or not
	if exists {
		return
	}
	//if not exist then create the minio's bucket "talklet-media"
	err := Minio_client.MakeBucket(context.Background(), Bucket_name, minio.MakeBucketOptions{Region: "us-east-1"})
	// For AWS S3, regions matter because buckets can be created in different geographic regions
	//by default location region = "us-east-1", this location is for multiple servers
	if err != nil {
		fmt.Println("error while creating the minio bucket - ", err)
	}

}
