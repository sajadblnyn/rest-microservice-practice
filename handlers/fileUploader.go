package handlers

import (
	"mime"
	"net/http"

	"github.com/google/uuid"
	"github.com/sajadblnyn/rest-microservice-practice/services/storage"
)

func Upload(rs http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(256 * 1024)

	if err != nil {
		http.Error(rs, "invalid multipart files", http.StatusBadRequest)
		return
	}

	fileName := uuid.NewString()

	f, fh, err := r.FormFile("file")

	if err != nil {
		http.Error(rs, "invalid file upload", http.StatusBadRequest)
		return
	}

	ct := fh.Header.Get("Content-Type")
	mimeType, _, err := mime.ParseMediaType(ct)
	if err != nil {
		http.Error(rs, "invalid file content type", http.StatusBadRequest)
		return
	}

	extension, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		http.Error(rs, "error in getting file extension", http.StatusInternalServerError)
		return
	}

	fileName = fileName + extension[0]

	st := storage.NewStorage()

	err = st.UploadFile(fileName, f)
	if err != nil {
		http.Error(rs, "error in uploading file", http.StatusInternalServerError)
		return
	}
}
