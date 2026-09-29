package handler

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadHandlerSize(t *testing.T) {
	buf := &bytes.Buffer{}

	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile("file", "am.txt")
	if err != nil {
		t.Fatalf("Error %v", err)
	}

	dummyData := make([]byte, 20<<20)
	part.Write(dummyData)

	writer.Close()

	req := httptest.NewRequest("POST", "/upload", buf)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	UploadHandler(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler passed file that is more than 15 MB but shouldn't")
	}
}

func TestUploadHandler(t *testing.T) {
	buf := &bytes.Buffer{}

	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile("file", "am.txt")
	if err != nil {
		t.Fatalf("Error %v", err)
	}

	data := []byte("Hello World")
	part.Write(data)

	writer.Close()

	req := httptest.NewRequest("POST", "/upload", buf)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	UploadHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler couldn't uploaded file")
	}
}

func TestDownloadHandler(t *testing.T) {
	filename := "am.txt"
	req := httptest.NewRequest("GET", "/download?file=am.txt", nil)
	rr := httptest.NewRecorder()

	headerExpected := fmt.Sprintf("attachment; filename=\"%s\"", filename)

	DownloadHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Error("Handler couldn't donwload file")
	}

	headerResult := rr.Header().Get("Content-Disposition")

	if headerExpected != headerResult {
		t.Error("Headers are different")
	}

}
