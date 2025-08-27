package minio

import (
	"context"
	"fmt"
	"tet/internals/models"
	"time"

	"github.com/minio/minio-go/v7"
)

// http://localhost:9000/talklet-media/profile/user_id/profile.png

// http://localhost:9000/talklet-media/chats/private/sender_id/receiver_id/msg_id/file_name.extension

func UploadFile_private_chats(msg models.Message) {
	object_name := fmt.Sprintf(`/chats/private/%s/%s/%d/%s`, msg.SenderID, msg.ReceiverID, msg.ID, msg.MetaData.FileName)
	actual_file := msg.ActualFile //file internally implements io.Reader so no need for conversion
	file_size := msg.MetaData.FileSize
	content_type := msg.MetaData.MimeType // image/jpeg ,audio/mpeg,application/pdf it doesn't mandatory but for the browser understanding we give this
	upload_status, err := Minio_client.PutObject(context.Background(), Bucket_name, object_name, actual_file, file_size, minio.PutObjectOptions{ContentType: content_type})
	if err != nil {
		fmt.Println("error while uploading privatechat files to minio - ", err)
		return
	}
	fmt.Println(" minio key for the key - ", upload_status.Key)
}

// http://localhost:9000/talklet-media/chats/groups/department_name/sender_id/group_id/msg_id/file_name.extension
func UploadFile_group_chats(group_id string, department string, msg models.Message) {
	object_name := fmt.Sprintf(`/chats/groups/%s/%s/%s/%d/%s`, department, msg.SenderID, group_id, msg.ID, msg.MetaData.FileName)
	actual_file := msg.ActualFile //file internally implements io.Reader so no need for conversion
	file_size := msg.MetaData.FileSize
	content_type := msg.MetaData.MimeType // image/jpeg ,audio/mpeg,application/pdf it doesn't mandatory but for the browser understanding we give this
	upload_status, err := Minio_client.PutObject(context.Background(), Bucket_name, object_name, actual_file, file_size, minio.PutObjectOptions{ContentType: content_type})
	if err != nil {
		fmt.Println("error while uploading groupchat files to minio - ", err)
		return
	}
	fmt.Println(" minio key for the key - ", upload_status.Key)
}

//http://localhost:9000/talklet-media/chats/private/sender_id/receiver_id/msg_id/file_name.extension

func GetFile_private_chats(msg models.Message) {

	var err error
	object_name := fmt.Sprintf(`/chats/private/%s/%s/%d/%s`, msg.SenderID, msg.ReceiverID, msg.ID, msg.MetaData.FileName)

	//generate the secure link , only one can access it with specified time (15 minutes)
	msg.MetaData.FileURL, err = Minio_client.PresignedGetObject(context.Background(), Bucket_name, object_name, 15*time.Minute, nil)
	if err != nil {
		fmt.Println("error while generating the privatechat presigned link from the minio - ", err)
	}
}

// http://localhost:9000/talklet-media/chats/groups/department_name/sender_id/group_id/msg_id/file_name.extension
func GetFile_group_chats(group_id string, department string, msg models.Message) {

	var err error
	object_name := fmt.Sprintf(`/chats/groups/%s/%s/%s/%d/%s`, department, msg.SenderID, group_id, msg.ID, msg.MetaData.FileName)

	msg.MetaData.FileURL, err = Minio_client.PresignedGetObject(context.Background(), Bucket_name, object_name, 15*time.Minute, nil)
	if err != nil {
		fmt.Println("error while generating the privatechat presigned link from the minio - ", err)
	}
}
