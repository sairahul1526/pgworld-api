package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// admin

// AdminGet .
func AdminGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	params := r.URL.Query()
	limitOffset := " "

	if _, ok := params["limit"]; ok {
		limitOffset += " limit " + params["limit"][0]
		delete(params, "limit")
	} else {
		limitOffset += " limit " + defaultLimit
	}

	offset := defaultOffset
	if _, ok := params["offset"]; ok {
		limitOffset += " offset " + params["offset"][0]
		offset = params["offset"][0]
		delete(params, "offset")
	} else {
		limitOffset += " offset " + defaultOffset
	}

	orderBy := " "

	if _, ok := params["orderby"]; ok {
		orderBy += " order by " + params["orderby"][0]
		delete(params, "orderby")
		if _, ok := params["sortby"]; ok {
			orderBy += " " + params["sortby"][0] + " "
			delete(params, "sortby")
		} else {
			orderBy += " asc "
		}
	} else {
		orderBy += " order by created_date_time desc "
	}

	resp := " * "
	if _, ok := params["resp"]; ok {
		resp = " " + params["resp"][0] + " "
		delete(params, "resp")
	}

	oneSignalID := ""
	if _, ok := params["oneSignalID"]; ok {
		oneSignalID = " " + params["oneSignalID"][0] + " "
		delete(params, "oneSignalID")
	}

	shouldMail := false
	shouldMailID := ""
	if _, ok := params["shouldMail"]; ok {
		shouldMailID = params["shouldMailID"][0]
		shouldMail = true
		delete(params, "shouldMail")
		delete(params, "shouldMailID")
	}

	where := ""
	init := false
	for key, val := range params {
		if init {
			where += " and "
		}
		if strings.EqualFold("name", key) {
			where += " `" + key + "` like '%" + val[0] + "%' "
		} else if strings.EqualFold("amount", key) {
			values := strings.Split(val[0], ",")
			where += " `" + key + "` between '" + values[0] + "' and '" + values[1] + "' "
		} else {
			where += " `" + key + "` = '" + val[0] + "' "
		}
		init = true
	}
	SQLQuery := " from `" + adminTable + "`"
	if strings.Compare(where, "") != 0 {
		SQLQuery += " where " + where
	}
	SQLQuery += orderBy
	if shouldMail {
		mailResults("select "+resp+SQLQuery, shouldMailID)
	}
	SQLQuery += limitOffset

	data, status, ok := selectProcess("select " + resp + SQLQuery)
	w.Header().Set("Status", status)
	if ok {
		response["data"] = data
		if len(data) > 0 && len(oneSignalID) > 0 {
			fmt.Println("oneSignalID", oneSignalID)
			// db.Exec("")
		}

		pagination := map[string]string{}
		if len(where) > 0 {
			count, _, _ := selectProcess("select count(*) as ctn from `" + adminTable + "` where " + where)
			pagination["total_count"] = count[0]["ctn"]
		} else {
			count, _, _ := selectProcess("select count(*) as ctn from `" + adminTable + "`")
			pagination["total_count"] = count[0]["ctn"]
		}
		pagination["count"] = strconv.Itoa(len(data))
		pagination["offset"] = offset
		response["pagination"] = pagination

		response["meta"] = setMeta(status, "ok", "")
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

// AdminAdd .
func AdminAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	body := map[string]string{}

	r.ParseMultipartForm(32 << 20)

	for key, value := range r.Form {
		body[key] = value[0]
	}
	fieldCheck := requiredFiledsCheck(body, adminRequiredFields)
	if len(fieldCheck) > 0 {
		SetReponseStatus(w, r, statusCodeBadRequest, fieldCheck+" required", dialogType, response)
		return
	}

	// log
	// logAction(body["admin_name"], "added admin", "1", body["hostel_id"])
	delete(body, "admin_name")

	body["status"] = "1"
	body["created_date_time"] = time.Now().UTC().String()

	var (
		ok     bool
		status string
	)
	for true {
		body["id"] = RandStringBytes(adminDigits)
		status, ok = insertSQL(adminTable, body)
		if !strings.EqualFold(status, statusCodeDuplicateEntry) {
			break
		}
	}
	w.Header().Set("Status", status)
	if ok {
		response["meta"] = setMeta(status, "Admin added", dialogType)
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

// AdminUpdate .
func AdminUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	body := map[string]string{}

	r.ParseMultipartForm(32 << 20)

	for key, value := range r.Form {
		body[key] = value[0]
	}

	for _, field := range adminRequiredFields {
		if _, ok := body[field]; ok {
			if len(body[field]) == 0 {
				SetReponseStatus(w, r, statusCodeBadRequest, field+" required", dialogType, response)
				return
			}
		}
	}

	// log
	// logAction(body["admin_name"], "updated admin", "1", body["hostel_id"])
	delete(body, "admin_name")

	body["modified_date_time"] = time.Now().UTC().String()

	status, ok := updateSQL(adminTable, r.URL.Query(), body)
	w.Header().Set("Status", status)
	if ok {
		response["meta"] = setMeta(status, "Admin updated", dialogType)
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
