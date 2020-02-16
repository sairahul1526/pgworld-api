package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
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

	update := map[string]string{"last_paid_date_time": body["paid_date_time"], "expiry_date_time": body["expiry_date_time"]}
	if strings.EqualFold(body["joining"], "1") {
		update["joining"] = "0"
	}
	status, ok := updateSQL(userTable, url.Values{"hostel_id": {body["hostel_id"]}, "id": {body["user_id"]}}, update)
	w.Header().Set("Status", status)
	if ok {
		if len(body["bill_id"]) == 0 {
			// log
			logAction(body["admin_name"], "accepted rent ₹"+body["amount"]+" from "+body["name"], "7", body["hostel_id"])
			for true {
				status, _ := insertSQL(billTable, map[string]string{"id": RandStringBytes(hostelDigits), "hostel_id": body["hostel_id"], "user_id": body["user_id"], "title": body["title"], "description": body["description"], "amount": body["amount"], "paid_date_time": body["paid_date_time"], "status": "1", "paid": "0", "document": body["document"], "type": body["type"], "payment": body["payment"], "transaction_id": body["transaction_id"], "billid": body["billid"]})
				if !strings.EqualFold(status, statusCodeDuplicateEntry) {
					break
				}
			}
		} else {
			// log
			logAction(body["admin_name"], "updated rent ₹"+body["amount"]+" from "+body["name"], "7", body["hostel_id"])
			updateSQL(billTable, url.Values{"id": {body["bill_id"]}, "hostel_id": {body["hostel_id"]}, "user_id": {body["user_id"]}}, map[string]string{"amount": body["amount"], "paid_date_time": body["paid_date_time"], "document": body["document"], "type": body["type"], "payment": body["payment"], "transaction_id": body["transaction_id"], "billid": body["billid"]})
		}
		if strings.EqualFold(body["joining"], "1") {
			db.Exec("update " + roomTable + " set room_joining = room_joining - 1 where id = '" + body["room_id"] + "' and room_joining > 0")
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

	status, ok := updateSQL(employeeTable, url.Values{"hostel_id": {body["hostel_id"]}, "id": {body["employee_id"]}}, map[string]string{"last_paid_date_time": body["paid_date_time"], "expiry_date_time": body["expiry_date_time"]})
	w.Header().Set("Status", status)
	if ok {
		if len(body["bill_id"]) == 0 {
			// log
			logAction(body["admin_name"], "gave salary ₹"+body["amount"]+" to "+body["name"], "8", body["hostel_id"])
			for true {
				status, _ := insertSQL(billTable, map[string]string{"id": RandStringBytes(hostelDigits), "hostel_id": body["hostel_id"], "employee_id": body["employee_id"], "title": body["title"], "description": body["description"], "amount": body["amount"], "paid_date_time": body["paid_date_time"], "status": "1", "paid": "1", "document": body["document"], "type": body["type"], "payment": body["payment"], "transaction_id": body["transaction_id"], "billid": body["billid"]})
				if !strings.EqualFold(status, statusCodeDuplicateEntry) {
					break
				}
			}
		} else {
			// log
			logAction(body["admin_name"], "updated rent ₹"+body["amount"]+" to "+body["name"], "8", body["hostel_id"])
			updateSQL(billTable, url.Values{"id": {body["bill_id"]}, "hostel_id": {body["hostel_id"]}, "employee_id": {body["employee_id"]}}, map[string]string{"amount": body["amount"], "paid_date_time": body["paid_date_time"], "document": body["document"], "type": body["type"], "payment": body["payment"], "transaction_id": body["transaction_id"], "billid": body["billid"]})
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
