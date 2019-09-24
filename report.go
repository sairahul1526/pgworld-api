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

	// rents and salary
	result, _, _ := selectProcess("select sum(amount) as `amount`, MONTH(paid_date_time) as `month`  from " + billTable + " where hostel_id = '" + r.FormValue("hostel_id") + "' and status = 1 and paid = 0 and user_id > 0 and date(paid_date_time) >= '" + r.FormValue("from") + "' and date(paid_date_time) <= '" + r.FormValue("to") + "' group by MONTH(paid_date_time)")

	rents := []map[string]string{}

	for _, data := range result {
		rents = append(rents, map[string]string{
			"title": getMonthName(data["month"]),
			"value": data["amount"],
		})
	}

	result, _, _ = selectProcess("select sum(amount) as `amount`, MONTH(paid_date_time) as `month`  from " + billTable + " where hostel_id = '" + r.FormValue("hostel_id") + "' and status = 1 and paid = 1 and employee_id > 0 and date(paid_date_time) >= '" + r.FormValue("from") + "' and date(paid_date_time) <= '" + r.FormValue("to") + "' group by MONTH(paid_date_time)")

	salaries := []map[string]string{}

	for _, data := range result {
		salaries = append(salaries, map[string]string{
			"title": getMonthName(data["month"]),
			"value": data["amount"],
		})
	}

	pies = append(pies, map[string]interface{}{
		"title": "Rents & Salary",
		"color": "#F5B7B1",
		"type":  "2",
		"data": []map[string]interface{}{
			map[string]interface{}{
				"title": "Rents",
				"color": "#F5B7B1",
				"data":  rents,
			},
			map[string]interface{}{
				"title": "Salary",
				"color": "#F5B7B1",
				"data":  salaries,
			},
		},
	})

	// income and expense
	result, _, _ = selectProcess("select sum(amount) as `amount`, MONTH(paid_date_time) as `month`  from " + billTable + " where hostel_id = '" + r.FormValue("hostel_id") + "' and status = 1 and paid = 0 and date(paid_date_time) >= '" + r.FormValue("from") + "' and date(paid_date_time) <= '" + r.FormValue("to") + "' group by MONTH(paid_date_time)")

	incomes := []map[string]string{}

	for _, data := range result {
		incomes = append(incomes, map[string]string{
			"title": getMonthName(data["month"]),
			"value": data["amount"],
		})
	}

	result, _, _ = selectProcess("select sum(amount) as `amount`, MONTH(paid_date_time) as `month`  from " + billTable + " where hostel_id = '" + r.FormValue("hostel_id") + "' and status = 1 and paid = 1 and date(paid_date_time) >= '" + r.FormValue("from") + "' and date(paid_date_time) <= '" + r.FormValue("to") + "' group by MONTH(paid_date_time)")

	expenses := []map[string]string{}

	for _, data := range result {
		expenses = append(expenses, map[string]string{
			"title": getMonthName(data["month"]),
			"value": data["amount"],
		})
	}

	pies = append(pies, map[string]interface{}{
		"title": "Income & Expense",
		"color": "#F5B7B1",
		"type":  "2",
		"data": []map[string]interface{}{
			map[string]interface{}{
				"title": "Income",
				"color": "#F5B7B1",
				"data":  incomes,
			},
			map[string]interface{}{
				"title": "Expense",
				"color": "#F5B7B1",
				"data":  expenses,
			},
		},
	})

	// payment types
	result, _, _ = selectProcess("select sum(amount) as `amount`, payment from " + billTable + " where hostel_id = '" + r.FormValue("hostel_id") + "' and status = 1 and date(paid_date_time) >= '" + r.FormValue("from") + "' and date(paid_date_time) <= '" + r.FormValue("to") + "' group by payment;")
	if len(result) > 0 {
		payments := []map[string]string{}
		for _, val := range result {
			switch val["payment"] {
			case "1":
				payments = append(payments, map[string]string{
					"title": "Credit Card",
					"shown": val["amount"],
					"value": val["amount"],
					"color": "#A2D9CE",
				})
			case "2":
				payments = append(payments, map[string]string{
					"title": "Debit Card",
					"shown": val["amount"],
					"value": val["amount"],
					"color": "#A2D9CE",
				})
			case "3":
				payments = append(payments, map[string]string{
					"title": "Net Banking",
					"shown": val["amount"],
					"value": val["amount"],
					"color": "#A2D9CE",
				})
			case "4":
				payments = append(payments, map[string]string{
					"title": "Google Pay",
					"shown": val["amount"],
					"value": val["amount"],
					"color": "#A2D9CE",
				})
			case "5":
				payments = append(payments, map[string]string{
					"title": "PhonePe",
					"shown": val["amount"],
					"value": val["amount"],
					"color": "#A2D9CE",
				})
			case "6":
				payments = append(payments, map[string]string{
					"title": "PayTM",
					"shown": val["amount"],
					"value": val["amount"],
					"color": "#A2D9CE",
				})
			case "7":
				payments = append(payments, map[string]string{
					"title": "Cash",
					"shown": result[0]["amount"],
					"value": result[0]["amount"],
					"color": "#A2D9CE",
				})
			case "8":
				payments = append(payments, map[string]string{
					"title": "Others",
					"shown": result[0]["amount"],
					"value": result[0]["amount"],
					"color": "#A2D9CE",
				})
			}
		}
		pies = append(pies, map[string]interface{}{
			"title":      "Payment Modes",
			"color":      "#F5B7B1",
			"type":       "2",
			"horizontal": "1",
			"data": []map[string]interface{}{
				map[string]interface{}{
					"title": "Amount",
					"color": "#F5B7B1",
					"data":  payments,
				},
			},
		})
	}

	// room filled and capacity
	result, _, _ = selectProcess("SELECT sum(capacity) as tot_cap, sum(filled) as tot_fill FROM " + roomTable + " where hostel_id = '" + r.FormValue("hostel_id") + "' and status = 1;")
	if len(result) > 0 {
		cap, _ := strconv.Atoi(result[0]["tot_cap"])
		fil, _ := strconv.Atoi(result[0]["tot_fill"])
		not := strconv.Itoa(cap - fil)
		pies = append(pies, map[string]interface{}{
			"title": "Rooms",
			"type":  "1",
			"data": []map[string]interface{}{
				map[string]interface{}{
					"title": "Rooms",
					"color": "#F5B7B1",
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
				},
			},
		})
	}

	// user active and expired
	result, _, _ = selectProcess("select count(*) as total_users, count(case when date(expiry_date_time) >= '" + strings.Split(time.Now().String(), " ")[0] + "' then 'active' end) as active_users, count(case when date(expiry_date_time) < '" + strings.Split(time.Now().String(), " ")[0] + "' then 'expired' end) as expired_users from " + userTable + " where hostel_id = '" + r.FormValue("hostel_id") + "' and status = 1")
	if len(result) > 0 {
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
							"color": "#F5B7B1",
						},
					},
				},
			},
		})
	}

	response["graphs"] = pies
	response["meta"] = setMeta(statusCodeOk, "ok", "")

	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	meta, required := checkAppUpdate(r)
	if required {
		response["meta"] = meta
	}
	json.NewEncoder(w).Encode(response)
}
