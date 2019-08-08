package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// user

// UserGet .
func UserGet(w http.ResponseWriter, r *http.Request) {
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
		orderBy += " order by " + userTable + "." + params["orderby"][0]
		delete(params, "orderby")
		if _, ok := params["sortby"]; ok {
			orderBy += " " + params["sortby"][0] + " "
			delete(params, "sortby")
		} else {
			orderBy += " asc "
		}
	} else {
		orderBy += " order by " + userTable + ".created_date_time desc "
	}

	resp := " users.*, rooms.roomno "
	if _, ok := params["resp"]; ok {
		resp = " " + params["resp"][0] + " "
		delete(params, "resp")
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
		if strings.EqualFold("name", key) ||
			strings.EqualFold("phone", key) ||
			strings.EqualFold("email", key) {
			where += " " + userTable + ".`" + key + "` like '%" + val[0] + "%' "
		} else if strings.EqualFold("rent", key) {
			values := strings.Split(val[0], ",")
			where += " " + userTable + ".`" + key + "` between '" + values[0] + "' and '" + values[1] + "' "
		} else {
			where += " " + userTable + ".`" + key + "` = '" + val[0] + "' "
		}
		init = true
	}
	SQLQuery := " from `" + userTable + "` left join `" + roomTable + "` on `" + userTable + "`.`room_id` = `" + roomTable + "`.`id` "
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
		newData := []map[string]string{}
		for _, val := range data {
			val["payment_status"] = "1"
			if len(val["expiry_date_time"]) > 0 && !strings.Contains(val["expiry_date_time"], "0000") {
				t, _ := time.Parse("2006-01-02", val["expiry_date_time"])
				diff := time.Since(t)
				if diff.Hours() > 0 {
					val["payment_status"] = "0"
				}
			}
			newData = append(newData, val)
		}
		response["data"] = newData

		pagination := map[string]string{}
		if len(where) > 0 {
			count, _, _ := selectProcess("select count(*) as ctn from `" + userTable + "` left join `" + roomTable + "` on `" + userTable + "`.`room_id` = `" + roomTable + "`.`id` where " + where)
			pagination["total_count"] = count[0]["ctn"]
		} else {
			count, _, _ := selectProcess("select count(*) as ctn from `" + userTable + "` left join `" + roomTable + "` on `" + userTable + "`.`room_id` = `" + roomTable + "`.`id` ")
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

// UserAdd .
func UserAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	body := map[string]string{}

	r.ParseMultipartForm(32 << 20)

	for key, value := range r.Form {
		body[key] = value[0]
	}
	fieldCheck := requiredFiledsCheck(body, userRequiredFields)
	if len(fieldCheck) > 0 {
		SetReponseStatus(w, r, statusCodeBadRequest, fieldCheck+" required", dialogType, response)
		return
	}

	// log
	logAction(body["admin_name"], "added user "+body["name"], "9", body["hostel_id"])
	delete(body, "admin_name")

	body["status"] = "1"
	body["created_date_time"] = time.Now().UTC().String()

	status, ok := insertSQL(userTable, body)
	w.Header().Set("Status", status)
	if ok {
		db.Exec("update " + roomTable + " set filled = filled + 1 where id = '" + body["room_id"] + "'")
		response["meta"] = setMeta(status, "User added", dialogType)
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

// UserUpdate .
func UserUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	body := map[string]string{}

	r.ParseMultipartForm(32 << 20)

	for key, value := range r.Form {
		body[key] = value[0]
	}

	for _, field := range userRequiredFields {
		if _, ok := body[field]; ok {
			if len(body[field]) == 0 {
				SetReponseStatus(w, r, statusCodeBadRequest, field+" required", dialogType, response)
				return
			}
		}
	}

	// log
	logAction(body["admin_name"], "updated user "+body["name"], "9", body["hostel_id"])
	delete(body, "admin_name")
	prevRoomID := body["prev_room_id"]
	delete(body, "prev_room_id")

	body["modified_date_time"] = time.Now().UTC().String()

	status, ok := updateSQL(userTable, r.URL.Query(), body)
	w.Header().Set("Status", status)
	if ok {
		db.Exec("update " + roomTable + " set filled = filled + 1 where id = '" + body["room_id"] + "'")
		db.Exec("update " + roomTable + " set filled = filled - 1 where id = '" + prevRoomID + "'")
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

// UserDelete .
func UserDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	body := map[string]string{}

	r.ParseMultipartForm(32 << 20)

	for key, value := range r.Form {
		body[key] = value[0]
	}

	// log
	logAction(body["admin_name"], "deleted user", "9", body["hostel_id"])
	delete(body, "admin_name")

	body["modified_date_time"] = time.Now().UTC().String()

	status, ok := updateSQL(userTable, r.URL.Query(), map[string]string{"status": "0"})
	w.Header().Set("Status", status)
	if ok {
		db.Exec("update " + roomTable + " set filled = filled - 1 where id = '" + body["room_id"] + "'")
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
