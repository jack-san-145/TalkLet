package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"tet/internals/models"
	"tet/internals/storage/minio"
	"tet/internals/storage/postgres"
	"time"
)

func ChatFileUploads(w http.ResponseWriter, r *http.Request) {

	isFound, senderId := FindCookie(r)
	if !isFound {
		return
	}

	var (
		media_msg   models.Message
		file_header *multipart.FileHeader
		// group_id    string
		msg_time time.Time
	)
	msg_time = time.Now()
	media_msg.CreatedAt = msg_time.Format("2006-01-02 15:04:05")
	media_msg.SenderID = senderId

	err := r.ParseMultipartForm(200 << 20) //users can send file upto the size of 200mb
	if err != nil {
		fmt.Println("error while parsing the media-msg request - ", err)
		return
	}
	media_msg.ActualFile, file_header, err = r.FormFile("actual_file")
	if err != nil {
		fmt.Println("error occured in the media-file - ", err)
		return
	}

	media_msg.Type = r.FormValue("type")
	media_msg.MetaData.MimeType = r.FormValue("mime_type")
	media_msg.Content = r.FormValue("content")
	media_msg.IsAck = "ack"
	media_msg.MetaData.FileName = file_header.Filename
	media_msg.MetaData.FileSize = file_header.Size

	is_group_str := r.FormValue("is_group")
	is_group_bool, _ := strconv.ParseBool(is_group_str)
	if is_group_bool {
		// group_id = r.FormValue("group_id")
		// code for store msg in db
	} else {
		media_msg.ReceiverID = r.FormValue("receiver_id")
		media_msg.ID = postgres.Store_Privatechat_MessagesPostDB(media_msg)
		if media_msg.ID != 0 {
			status := minio.UploadFile_private_chats(&media_msg)
			WriteJSON(w, r, status)
		}
		minio.GetFile_private_chats(&media_msg)
		// fmt.Println("presigned url - ", media_msg.MetaData.FileURL)

	}

}
