package main

import (
	"encoding/json"
	"net/http"
)

// HostelStatus .
func HostelStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	setHostelExpiry(r.FormValue("hostel_id"))

	response["meta"] = setMeta(statusCodeOk, "Updated", dialogType)
	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	json.NewEncoder(w).Encode(response)
}
