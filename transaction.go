package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// transaction

// Rent .
func Rent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	body := map[string]string{}

	r.ParseMultipartForm(32 << 20)

	for key, value := range r.Form {
		body[key] = value[0]
	}

	status, ok := updateSQL(userTable, url.Values{"hostel_id": {body["hostel_id"]}, "id": {body["user_id"]}}, map[string]string{"last_paid_date_time": body["last_paid_date_time"]})
	w.Header().Set("Status", status)
	if ok {
		if len(body["bill_id"]) == 0 {
			// log
			logAction(body["admin_name"], "accepted rent")
			insertSQL(billTable, map[string]string{"hostel_id": body["hostel_id"], "user_id": body["user_id"], "title": body["title"], "description": body["description"], "amount": body["amount"], "bill_date_time": body["bill_date_time"], "status": "1", "paid": "0"})
		} else {
			// log
			logAction(body["admin_name"], "updated rent")
			updateSQL(billTable, url.Values{"id": {body["bill_id"]}, "hostel_id": {body["hostel_id"]}, "user_id": {body["user_id"]}}, map[string]string{"amount": body["amount"], "bill_date_time": body["bill_date_time"]})
		}
		response["meta"] = setMeta(status, "User updated", dialogType)
	} else {
		response["meta"] = setMeta(status, "", dialogType)
	}

	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	meta, required := checkAppUpdate(r)
	if required {
		response["meta"] = meta
	}
	json.NewEncoder(w).Encode(response)
}

// Salary .
func Salary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	body := map[string]string{}

	r.ParseMultipartForm(32 << 20)

	for key, value := range r.Form {
		body[key] = value[0]
	}

	status, ok := updateSQL(employeeTable, url.Values{"hostel_id": {body["hostel_id"]}, "id": {body["employee_id"]}}, map[string]string{"last_paid_date_time": body["last_paid_date_time"]})
	w.Header().Set("Status", status)
	if ok {
		if len(body["bill_id"]) == 0 {
			// log
			logAction(body["admin_name"], "accepted rent")
			insertSQL(billTable, map[string]string{"hostel_id": body["hostel_id"], "employee_id": body["employee_id"], "title": body["title"], "description": body["description"], "amount": body["amount"], "bill_date_time": body["bill_date_time"], "status": "1", "paid": "0"})
		} else {
			// log
			logAction(body["admin_name"], "updated rent")
			updateSQL(billTable, url.Values{"id": {body["bill_id"]}, "hostel_id": {body["hostel_id"]}, "employee_id": {body["employee_id"]}}, map[string]string{"amount": body["amount"], "bill_date_time": body["bill_date_time"]})
		}
		response["meta"] = setMeta(status, "User updated", dialogType)
	} else {
		response["meta"] = setMeta(status, "", dialogType)
	}

	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	meta, required := checkAppUpdate(r)
	if required {
		response["meta"] = meta
	}
	json.NewEncoder(w).Encode(response)
}
