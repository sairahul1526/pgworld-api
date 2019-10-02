package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	// issue
	result, _, _ = selectProcess("select count(*) as ctn from " + issueTable + " where resolve = '0' and status = '1' and hostel_id = '" + r.FormValue("hostel_id") + "'")
	counts["issue"] = result[0]["ctn"]

	// pies
	pies := []map[string]interface{}{}

	// room filled and capacity
	result, _, _ = selectProcess("SELECT sum(capacity) as tot_cap, sum(filled) as tot_fill FROM " + roomTable + " where hostel_id = '" + r.FormValue("hostel_id") + "' and status = 1;")
	cap, _ := strconv.Atoi(result[0]["tot_cap"])
	fil, _ := strconv.Atoi(result[0]["tot_fill"])
	not := strconv.Itoa(cap - fil)
	pies = append(pies, map[string]interface{}{
		"title": "Beds",
		"type":  "2",
		"data": []map[string]interface{}{
			map[string]interface{}{
				"title": "Beds",
				"type":  "1",
				"data": []map[string]string{
					map[string]string{
						"title": "Filled",
						"shown": result[0]["tot_fill"],
						"value": result[0]["tot_fill"],
						"color": "#9AD7CB",
					},
					map[string]string{
						"title": "Vacant",
						"shown": not,
						"value": not,
						"color": "#E1A1AD",
					},
				},
			},
		},
	})

	// user active and expired
	result, _, _ = selectProcess("select count(*) as total_users, count(case when date(expiry_date_time) >= '" + strings.Split(time.Now().String(), " ")[0] + "' then 'active' end) as active_users, count(case when date(expiry_date_time) < '" + strings.Split(time.Now().String(), " ")[0] + "' then 'expired' end) as expired_users from " + userTable + " where hostel_id = '" + r.FormValue("hostel_id") + "' and status = 1")
	pies = append(pies, map[string]interface{}{
		"title": "Users",
		"type":  "1",
		"data": []map[string]interface{}{
			map[string]interface{}{
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
						"color": "#E1A1AD",
					},
				},
			},
		},
	})

	fromDate := strings.Split(time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local).String(), " ")[0]
	toDate := strings.Split(time.Now().AddDate(0, 0, 1).String(), " ")[0]

	result, _, _ = selectProcess("select paid, sum(amount) as `amount`  from " + billTable + " where hostel_id = '" + r.FormValue("hostel_id") + "' and status = 1 and date(paid_date_time) >= '" + fromDate + "' and date(paid_date_time) <= '" + toDate + "' group by MONTH(paid)")

	data := []map[string]string{}
	var income, expense int
	for _, res := range result {
		if strings.EqualFold(res["paid"], "1") {
			data = append(data, map[string]string{
				"title": "Expense",
				"value": res["amount"],
				"color": "#E1A1AD",
			})
			expense, _ = strconv.Atoi(res["amount"])
		} else {
			data = append(data, map[string]string{
				"title": "Income",
				"value": res["amount"],
				"color": "#9AD7CB",
			})
			income, _ = strconv.Atoi(res["amount"])
		}
	}
	if income-expense > 0 {
		data = append(data, map[string]string{
			"title": "Total",
			"value": "+" + strconv.Itoa(income-expense),
			"color": "#9AD7CB",
		})
	} else {
		data = append(data, map[string]string{
			"title": "Total",
			"value": strconv.Itoa(income - expense),
			"color": "#E1A1AD",
		})
	}

	pies = append(pies, map[string]interface{}{
		"title": "Income & Expenses",
		"type":  "3",
		"data": []map[string]interface{}{
			map[string]interface{}{
				"title": "Income & Expenses",
				"type":  "1",
				"data":  data,
			},
		},
	})

	response["graphs"] = pies
	response["data"] = []map[string]string{counts}
	response["meta"] = setMeta(statusCodeOk, "ok", "")

	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	meta, required := checkAppUpdate(r)
	if required {
		response["meta"] = meta
	}
	json.NewEncoder(w).Encode(response)
}
