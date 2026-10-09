package handler

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"x/onte-server/internal"
)

func TestUploadHandlerSize(t *testing.T) {
	buf := &bytes.Buffer{}
	identifier := internal.GenerateIdentifier()

	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile("file", "am.txt")
	if err != nil {
		t.Fatalf("Error %v", err)
	}

	dummyData := make([]byte, 1200<<20)
	part.Write(dummyData)

	writer.Close()

	req := httptest.NewRequest("POST", "/upload", buf)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("identifier", identifier)

	rr := httptest.NewRecorder()

	UploadHandler(rr, req, identifier)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler passed file that is more than 1 GB but shouldn't.")
	}
}

func TestUploadHandler(t *testing.T) {
	buf := &bytes.Buffer{}
	identifier := internal.GenerateIdentifier()

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
	req.Header.Set("identifier", identifier)

	rr := httptest.NewRecorder()

	UploadHandler(rr, req, identifier)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler couldn't uploaded file.")
	}
}

func TestUploadHandlerIdentifierVericifation(t *testing.T) {
	buf := &bytes.Buffer{}
	identifier := internal.GenerateIdentifier()

	req := httptest.NewRequest("POST", "/upload", buf)

	req.Header.Set("Identifier", "1")

	rr := httptest.NewRecorder()

	UploadHandler(rr, req, identifier)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("Handler passed Header but shouldn't.")
	}

}

func TestDownloadHandler(t *testing.T) {
	defer os.RemoveAll("uploads")
	identifier := internal.GenerateIdentifier()
	filename := "am.txt"
	req := httptest.NewRequest("GET", "/download?file=am.txt", nil)
	rr := httptest.NewRecorder()

	headerExpected := fmt.Sprintf("attachment; filename=\"%s\"", filename)

	req.Header.Set("identifier", identifier)

	DownloadHandler(rr, req, identifier)

	if status := rr.Code; status != http.StatusOK {
		t.Error("Handler couldn't donwload file")
	}

	headerResult := rr.Header().Get("Content-Disposition")

	if headerExpected != headerResult {
		t.Error("Headers are different")
	}
}

func TestDownloadHandlerIdentifierVericifation(t *testing.T) {
	buf := &bytes.Buffer{}
	identifier := internal.GenerateIdentifier()

	req := httptest.NewRequest("GET", "/download", buf)

	req.Header.Set("Identifier", "1")

	rr := httptest.NewRecorder()

	DownloadHandler(rr, req, identifier)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("Handler passed Header but shouldn't.")
	}

}
