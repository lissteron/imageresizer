package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"net/http"

	"github.com/disintegration/imaging"
)

const (
	fullQuality = 100
	cropSize    = 1

	httpPath = "/"
	httpPort = 8080

	fileName = "image"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate)

	http.HandleFunc(httpPath, resizeHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil))
}

func resizeHandler(w http.ResponseWriter, r *http.Request) {
	fh, _, err := r.FormFile(fileName)
	if err != nil {
		log.Println("[error]", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	img, err := imaging.Decode(fh)
	if err != nil {
		log.Println("[error]", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	img = imaging.CropAnchor(
		img,
		img.Bounds().Size().X-cropSize,
		img.Bounds().Size().Y-cropSize,
		imaging.Center,
	)

	if err := jpeg.Encode(w, img, &jpeg.Options{Quality: fullQuality}); err != nil {
		log.Println("[error]", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
