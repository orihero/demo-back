package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"../models"
	"../utils"
	"github.com/gorilla/mux"
)

func MultipleFileUpload(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile
	fhs := req.MultipartForm.File["files"]
	var urls []string
	for _, fh := range fhs {
		f, err := fh.Open()
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			var err models.Error
			err = utils.SetError(err, "File damaged or corrupted")
			return
		}
		// f is one of the files

		data, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
		}
		err = ioutil.WriteFile(fmt.Sprintf("./public/uploads/%s", fh.Filename), data, 0644)
		if err != nil {
			fmt.Println(err)
		}
		urls = append(urls, fh.Filename)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(urls)
}

func GetUploadedFiles(w http.ResponseWriter, r *http.Request) {
	fileName := mux.Vars(r)["name"]
	img, err := os.Open(fmt.Sprintf("./public/uploads/%s", fileName))
	if err != nil {
		log.Fatal(err)
	}
	var bytes []byte
	_, err = img.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	//data, err := ioutil.ReadFile(fmt.Sprintf("public/uploads/%s", "2.jpg"))
	//log.Println(data)
	//defer img.Close()
	w.Header().Set("Content-Type", "text")
	io.Copy(w, img)
}
