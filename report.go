package main

import (
	"encoding/json"
	"net/http"
)

// Report .
func Report(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// pies
	pies := [][]map[string]string{}
	result, _, _ := selectProcess("select paid as title, count(*) as value from " + billTable + " where status = '1' and hostel_id = '" + r.FormValue("hostel_id") + "' group by paid")
	pies = append(pies, result)
	response["pies"] = pies

	// bars
	bars := []map[string]interface{}{}
	// bar := map[string]interface{}{}
	// bar["barData"], _, _ = selectProcess("select paid as title, count(*) as value from " + billTable + " where status = '1' and hostel_id = '" + r.FormValue("hostel_id") + "' group by paid")
	// bars = append(bars, bar)
	response["bars"] = bars

	response["meta"] = setMeta(statusCodeOk, "ok", "")

	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	meta, required := checkAppUpdate(r)
	if required {
		response["meta"] = meta
	}
	json.NewEncoder(w).Encode(response)
}
