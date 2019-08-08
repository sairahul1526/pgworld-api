package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// SendOTP .
func SendOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	response["status"] = "200"
	response["message"] = "OTP sent"
	response["message_type"] = "1"

	json.NewEncoder(w).Encode(response)
}

// VerifyOTP .
func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	if strings.EqualFold(r.FormValue("otp"), "4242") {
		response["status"] = "200"
		response["message"] = ""
		response["message_type"] = ""
	} else {
		response["status"] = "400"
		response["message"] = "Incorrect OTP"
		response["message_type"] = "1"
	}

	json.NewEncoder(w).Encode(response)
}
