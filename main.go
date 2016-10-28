package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Image defines the location and name of an image..
type Image struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var baseFolder = "images"
var baseURL = "/api/static"

func imageListHandler(w http.ResponseWriter, r *http.Request) {

	queryValues := r.URL.Query()

	startValue, startErr := strconv.Atoi(queryValues.Get("start"))
	sizeValue, sizeErr := strconv.Atoi(queryValues.Get("size"))
	if startErr != nil || sizeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	images := make([]Image, sizeValue)
	for i := 0; i < sizeValue; i++ {
		images[i] = Image{
			Name: fmt.Sprintf("image %d", startValue+i),
			URL:  fmt.Sprintf("%s/%s.%s", baseURL, "image", "jpg"),
		}
	}

	j, err := json.Marshal(images)
	if err != nil {
		log.Printf("Error in imageListHandler [%s]", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func main() {

	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/image", imageListHandler).Methods("GET")

	r.PathPrefix(baseURL).Handler(http.StripPrefix(baseURL, http.FileServer(http.Dir(baseFolder))))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Listening at 8080")
	server.ListenAndServe()
}
