package helper

import (
	"io"
	"mime/multipart"
)

type FileData struct {
	Content io.Reader
	FileName string
	Size int64
	Header *multipart.FileHeader
}