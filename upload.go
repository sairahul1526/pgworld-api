package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

// upload

// Upload .
func Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	body := map[string]string{}

	// file upload
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("photo")
	if err != nil {
		fmt.Println("Upload", err)
	}
	if file != nil {
		defer file.Close()

		fileName, uploaded := uploadToS3(docS3Path, file, filepath.Ext(handler.Filename))
		if !uploaded {
			fmt.Println("Upload", err)
			SetReponseStatus(w, r, statusCodeServerError, "", dialogType, response)
			return
		}
		body["image"] = fileName
	}
	//

	response["data"] = []map[string]string{body}
	response["meta"] = setMeta(statusCodeOk, "Image upload", "")

	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	meta, required := checkAppUpdate(r)
	if required {
		response["meta"] = meta
	}
	json.NewEncoder(w).Encode(response)
}
