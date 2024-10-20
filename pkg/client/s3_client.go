package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"nftvc-profile/pkg/config"
	"nftvc-profile/pkg/logger"
	"path/filepath"
)

type S3Client struct {
	cfg *config.Config
	log logger.Logger
}

func NewS3Client(cfg *config.Config, log logger.Logger) *S3Client {
	return &S3Client{cfg: cfg, log: log}
}

func (s *S3Client) UploadFile(file []byte, filename string) (string, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, bytes.NewReader(file))
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/upload", s.cfg.S3Url)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
