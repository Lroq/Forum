package handlers

import ( 
    "net/http"
    "errors"
    "mime/multipart"
)

func checkFileTypeAndSize(r *http.Request) (multipart.File, error) {
    r.ParseMultipartForm(20 << 20) 

    file, header, err := r.FormFile("image")
    if err != nil {
        return nil, err
    }

    buffer := make([]byte, 512)
    file.Read(buffer)
    fileType := http.DetectContentType(buffer)

    if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/gif" && fileType != "image/svg+xml" {
        return nil, errors.New("invalid file type")
    }

    if header.Size > (20 << 20) {
        return nil, errors.New("file size is too big")
    }

    return file, nil
}