package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Report .
func Report(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// pies
	pies := []map[string]interface{}{}

	// room filled and capacity
	result, _, _ := selectProcess("SELECT sum(capacity) as tot_cap, sum(filled) as tot_fill FROM " + roomTable + " where hostel_id = '" + r.FormValue("hostel_id") + "' and status = 1;")
	cap, _ := strconv.Atoi(result[0]["tot_cap"])
	fil, _ := strconv.Atoi(result[0]["tot_fill"])
	not := strconv.Itoa(cap - fil)
	pies = append(pies, map[string]interface{}{
		"title": "Rooms",
		"type":  "1",
		"data": []map[string]string{
			// map[string]string{
			// 	"title": "Capacity",
			// 	"shown": result[0]["tot_cap"],
			// 	"value": result[0]["tot_cap"],
			// 	"color": "#AED6F1",
			// },
			map[string]string{
				"title": "Filed",
				"shown": result[0]["tot_fill"],
				"value": result[0]["tot_fill"],
				"color": "#A2D9CE",
			},
			map[string]string{
				"title": "Vacant",
				"shown": not,
				"value": not,
				"color": "#F5B7B1",
			},
		},
	})

	// user active and expired
	result, _, _ = selectProcess("select count(*) as total_users, count(case when date(expiry_date_time) >= '" + strings.Split(time.Now().String(), " ")[0] + "' then 'active' end) as active_users, count(case when date(expiry_date_time) < '" + strings.Split(time.Now().String(), " ")[0] + "' then 'expired' end) as expired_users from " + userTable + " where hostel_id = '" + r.FormValue("hostel_id") + "' and status = 1")
	pies = append(pies, map[string]interface{}{
		"title": "Users",
		"type":  "1",
		"data": []map[string]string{
			map[string]string{
				"title": "Total",
				"shown": result[0]["total_users"],
				"value": result[0]["total_users"],
				"color": "#AED6F1",
			},
			map[string]string{
				"title": "Active",
				"shown": result[0]["active_users"],
				"value": result[0]["active_users"],
				"color": "#A2D9CE",
			},
			map[string]string{
				"title": "Due",
				"shown": result[0]["expired_users"],
				"value": result[0]["expired_users"],
				"color": "#F5B7B1",
			},
		},
	})

	result, _, _ = selectProcess("select sum(amount) as `amount`, MONTH(paid_date_time) as `month`  from " + billTable + " where hostel_id = '" + r.FormValue("hostel_id") + "' and status = 1 and paid = 0 and user_id > 0 and hostel_id = '1' and date(paid_date_time) >= '" + r.FormValue("from") + "' and date(paid_date_time) <= '" + r.FormValue("to") + "' group by MONTH(paid_date_time)")

	data := []map[string]string{}
	max := 0
	for _, res := range result {
		amount, _ := strconv.Atoi(res["amount"])
		if max < amount {
			max = amount
		}
		var m, _ = strconv.Atoi(res["month"])
		data = append(data, map[string]string{
			"title": strings.Split(time.Date(2019, time.Month(m), 1, 0, 0, 0, 0, time.Local).String(), " ")[0],
			"value": res["amount"],
		})
	}

	pies = append(pies, map[string]interface{}{
		"title":      "Rents",
		"color":      "#F5B7B1",
		"data_title": "/-",
		"type":       "2",
		"steps":      strconv.Itoa(max / 10),
		"data":       data,
	})

	response["graphs"] = pies
	response["meta"] = setMeta(statusCodeOk, "ok", "")

	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	meta, required := checkAppUpdate(r)
	if required {
		response["meta"] = meta
	}
	json.NewEncoder(w).Encode(response)
}
