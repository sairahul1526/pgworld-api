package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Payment .
func Payment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	body := map[string]string{}

	r.ParseMultipartForm(32 << 20)

	for key, value := range r.Form {
		body[key] = value[0]
	}

	for _, field := range paymentRequiredFields {
		if _, ok := body[field]; ok {
			if len(body[field]) == 0 {
				SetReponseStatus(w, r, statusCodeBadRequest, field+" required", dialogType, response)
				return
			}
		}
	}

	// log
	// logAction(body["admin_name"], "updated invoice", "4", body["invoice_id"])
	delete(body, "admin_name")

	body["status"] = "1"
	body["created_date_time"] = time.Now().UTC().String()

	// capture payment
	req, _ := http.NewRequest("POST", "https://api.razorpay.com/v1/payments/"+body["payment_id"]+"/capture", nil)
	req.Header.Add("Authorization", razorpayAuth)
	q := req.URL.Query()
	q.Add("amount", body["amount"])
	q.Add("currency", "INR")
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		SetReponseStatus(w, r, statusCodeBadRequest, "", dialogType, response)
		return
	}

	defer res.Body.Close()
	resp, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(resp))
	if strings.Contains(string(resp), "BAD_REQUEST_ERROR") {
		SetReponseStatus(w, r, statusCodeBadRequest, "", dialogType, response)
		return
	}

	var (
		status string
	)
	for true {
		body["id"] = RandStringBytes(invoiceDigits)
		status, _ = insertSQL(invoiceTable, body)
		if !strings.EqualFold(status, statusCodeDuplicateEntry) {
			break
		}
	}

	// update hostels
	var hostels string
	db.QueryRow("select hostels from " + adminTable + " where id = '" + body["admin_id"] + "'").Scan(&hostels)

	var expiry string
	db.QueryRow("select expiry_date_time from " + hostelTable + " where id = '" + strings.Split(hostels, ",")[0] + "'").Scan(&expiry)
	t, _ := time.Parse("2006-01-02", expiry)
	diff := time.Since(t)
	if diff.Hours() > 0 {
		db.Exec("update " + hostelTable + " set expiry_date_time = '" + t.AddDate(0, pro[body["amount"]], 0).Format("2006-01-02") + "' where id in ('" + strings.Join(strings.Split(hostels, ","), "','") + "')")
	} else {
		db.Exec("update " + hostelTable + " set expiry_date_time = '" + time.Now().AddDate(0, pro[body["amount"]], 0).Format("2006-01-02") + "' where id in ('" + strings.Join(strings.Split(hostels, ","), "','") + "')")
	}

	for _, hostelID := range strings.Split(hostels, ",") {
		setHostelExpiry(hostelID)
	}

	w.Header().Set("Status", status)
	response["meta"] = setMeta(statusCodeOk, "", "")

	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	meta, required := checkAppUpdate(r)
	if required {
		response["meta"] = meta
	}
	json.NewEncoder(w).Encode(response)
}
