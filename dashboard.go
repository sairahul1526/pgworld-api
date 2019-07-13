package main

import (
	"encoding/json"
	"net/http"
)

// Dashboard .
func Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	counts := map[string]string{}

	// user
	result, _, _ := selectProcess("select count(*) as ctn from " + userTable + " where status = '1' and hostel_id = '" + r.FormValue("hostel_id") + "'")
	counts["user"] = result[0]["ctn"]
	// room
	result, _, _ = selectProcess("select count(*) as ctn from " + roomTable + " where status = '1' and hostel_id = '" + r.FormValue("hostel_id") + "'")
	counts["room"] = result[0]["ctn"]
	// bill
	result, _, _ = selectProcess("select count(*) as ctn from " + billTable + " where status = '1' and hostel_id = '" + r.FormValue("hostel_id") + "'")
	counts["bill"] = result[0]["ctn"]
	// note
	result, _, _ = selectProcess("select count(*) as ctn from " + noteTable + " where status = '1' and hostel_id = '" + r.FormValue("hostel_id") + "'")
	counts["note"] = result[0]["ctn"]
	// employee
	result, _, _ = selectProcess("select count(*) as ctn from " + employeeTable + " where status = '1' and hostel_id = '" + r.FormValue("hostel_id") + "'")
	counts["employee"] = result[0]["ctn"]

	response["data"] = []map[string]string{counts}
	response["meta"] = setMeta(statusCodeOk, "ok", "")

	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	meta, required := checkAppUpdate(r)
	if required {
		response["meta"] = meta
	}
	json.NewEncoder(w).Encode(response)
}
