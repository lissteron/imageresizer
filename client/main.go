package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const (
	fileName    = "1.jpg"
	newFileName = "new.jpg"

	serverURL = "http://localhost:8080"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate)

	var (
		b bytes.Buffer
		w = multipart.NewWriter(&b)
	)

	part, err := w.CreateFormFile("image", fileName)
	if err != nil {
		log.Println("[error]", err)

		return
	}

	fh, err := os.Open(fileName)
	if err != nil {
		log.Println("[error]", err)

		return
	}

	defer fh.Close()

	if _, err := io.Copy(part, fh); err != nil {
		log.Println("[error]", err)

		return
	}

	if err := w.Close(); err != nil {
		log.Println("[error]", err)

		return
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", serverURL, &b)
	if err != nil {
		log.Println("[error]", err)

		return
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		log.Println("[error]", err)

		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[error]", err)

		return
	}

	if err := ioutil.WriteFile(newFileName, body, 0600); err != nil {
		log.Println("[error]", err)

		return
	}
}
