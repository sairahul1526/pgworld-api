package main

import (
	"encoding/json"
	"net/http"
)

// StatusGet .
func StatusGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	response["meta"] = setMeta(statusCodeOk, "", "")
	// if getHostelStatus(r.FormValue("hostel_id")) {
	// 	response["meta"] = setMeta(statusCodeOk, "", "")
	// } else {
	// 	response["meta"] = setMeta(statusCodeForbidden, "", "")
	// }

	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	json.NewEncoder(w).Encode(response)
}

// StatusSet .
func StatusSet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	setHostelExpiry(r.FormValue("hostel_id"))

	response["meta"] = setMeta(statusCodeOk, "Updated", dialogType)
	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	json.NewEncoder(w).Encode(response)
}
