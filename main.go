package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/akrylysov/algnhsa"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func connectDatabase() {
	if db == nil {
		var err error
		db, err = sql.Open("mysql", dbConfig)
		if err != nil {
			log.Fatal(err)
		}
		db.SetMaxOpenConns(connectionPool)
		db.SetMaxIdleConns(connectionPool)
		db.SetConnMaxLifetime(time.Hour)
	}
}

// HealthCheck .
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("ok")
}

func main() {

	dbConfig = os.Getenv("dbConfig")
	connectionPool, _ = strconv.Atoi(os.Getenv("connectionPool"))
	test, _ = strconv.ParseBool(os.Getenv("test"))
	migrate, _ = strconv.ParseBool(os.Getenv("migrate"))
	s3Bucket = os.Getenv("s3Bucket")
	baseURL = os.Getenv("baseURL")
	supportEmailID = os.Getenv("supportEmailID")
	supportEmailPassword = os.Getenv("supportEmailPassword")
	supportEmailHost = os.Getenv("supportEmailHost")
	supportEmailPort, _ = strconv.Atoi(os.Getenv("supportEmailPort"))

	rand.Seed(time.Now().UnixNano())
	if db == nil {
		connectDatabase()
	}
	router := mux.NewRouter()

	// dashboard
	router.HandleFunc("/dashboard", checkHeaders(Dashboard)).Queries(
		"hostel_id", "{hostel_id}",
	).Methods("GET")

	// otp
	router.HandleFunc("/sendotp", checkHeaders(SendOTP)).Methods("GET")
	router.HandleFunc("/verifyotp", checkHeaders(VerifyOTP)).Methods("GET")

	// report
	router.HandleFunc("/report", checkHeaders(Report)).Queries(
		"hostel_id", "{hostel_id}",
	).Methods("GET")

	// admin
	router.HandleFunc("/admin", checkHeaders(AdminGet)).Queries(
		"username", "{username}",
	).Methods("GET")
	router.HandleFunc("/admin", checkHeaders(AdminAdd)).Methods("POST")
	router.HandleFunc("/admin", checkHeaders(AdminUpdate)).Methods("PUT")

	// bill
	router.HandleFunc("/bill", checkHeaders(BillGet)).Queries(
		"hostel_id", "{hostel_id}",
	).Methods("GET")
	router.HandleFunc("/bill", checkHeaders(BillAdd)).Methods("POST")
	router.HandleFunc("/bill", checkHeaders(BillUpdate)).Queries(
		"id", "{id}",
		"hostel_id", "{hostel_id}",
	).Methods("PUT")

	// employee
	router.HandleFunc("/employee", checkHeaders(EmployeeGet)).Queries(
		"hostel_id", "{hostel_id}",
	).Methods("GET")
	router.HandleFunc("/employee", checkHeaders(EmployeeAdd)).Methods("POST")
	router.HandleFunc("/employee", checkHeaders(EmployeeUpdate)).Queries(
		"id", "{id}",
		"hostel_id", "{hostel_id}",
	).Methods("PUT")

	// invoice
	router.HandleFunc("/food", checkHeaders(FoodGet)).Queries(
		"hostel_id", "{hostel_id}",
	).Methods("GET")
	router.HandleFunc("/food", checkHeaders(FoodAdd)).Methods("POST")
	router.HandleFunc("/food", checkHeaders(FoodUpdate)).Queries(
		"date", "{date}",
		"hostel_id", "{hostel_id}",
	).Methods("PUT")

	// invoice
	router.HandleFunc("/invoice", checkHeaders(InvoiceGet)).Queries(
		"admin_id", "{admin_id}",
	).Methods("GET")

	// issue
	router.HandleFunc("/issue", checkHeaders(IssueGet)).Queries(
		"hostel_id", "{hostel_id}",
	).Methods("GET")
	router.HandleFunc("/issue", checkHeaders(IssueAdd)).Methods("POST")
	router.HandleFunc("/issue", checkHeaders(IssueUpdate)).Queries(
		"id", "{id}",
		"hostel_id", "{hostel_id}",
	).Methods("PUT")

	// hostel
	router.HandleFunc("/hostel", checkHeaders(HostelGet)).Queries(
		"id", "{id}",
	).Methods("GET")
	router.HandleFunc("/hostel", checkHeaders(HostelGet)).Queries(
		"name", "{name}",
	).Methods("GET")
	router.HandleFunc("/hostel", checkHeaders(HostelAdd)).Methods("POST")
	router.HandleFunc("/hostel", checkHeaders(HostelUpdate)).Queries(
		"id", "{id}",
	).Methods("PUT")

	// log
	router.HandleFunc("/log", checkHeaders(LogGet)).Queries(
		"hostel_id", "{hostel_id}",
	).Methods("GET")

	// note
	router.HandleFunc("/note", checkHeaders(NoteGet)).Queries(
		"hostel_id", "{hostel_id}",
	).Methods("GET")
	router.HandleFunc("/note", checkHeaders(NoteAdd)).Methods("POST")
	router.HandleFunc("/note", checkHeaders(NoteUpdate)).Methods("PUT")

	// notice
	router.HandleFunc("/notice", checkHeaders(NoticeGet)).Queries(
		"hostel_id", "{hostel_id}",
	).Methods("GET")
	router.HandleFunc("/notice", checkHeaders(NoticeAdd)).Methods("POST")
	router.HandleFunc("/notice", checkHeaders(NoticeUpdate)).Queries(
		"id", "{id}",
		"hostel_id", "{hostel_id}",
	).Methods("PUT")

	//  payment
	router.HandleFunc("/payment", checkHeaders(Payment)).Methods("POST")

	// room
	router.HandleFunc("/room", checkHeaders(RoomGet)).Queries(
		"hostel_id", "{hostel_id}",
	).Methods("GET")
	router.HandleFunc("/room", checkHeaders(RoomAdd)).Methods("POST")
	router.HandleFunc("/room", checkHeaders(RoomUpdate)).Queries(
		"id", "{id}",
		"hostel_id", "{hostel_id}",
	).Methods("PUT")

	// set hostel status
	router.HandleFunc("/status", StatusGet).Queries(
		"hostel_id", "{hostel_id}",
	).Methods("GET")
	router.HandleFunc("/status", StatusSet).Methods("POST")

	// transactions
	router.HandleFunc("/rent", checkHeaders(Rent)).Methods("POST")
	router.HandleFunc("/salary", checkHeaders(Salary)).Methods("POST")

	// signup
	router.HandleFunc("/signup", checkHeaders(SignupGet)).Methods("GET")
	router.HandleFunc("/signup", checkHeaders(SignupAdd)).Methods("POST")
	router.HandleFunc("/signup", checkHeaders(SignupUpdate)).Methods("PUT")

	// support
	router.HandleFunc("/support", checkHeaders(SupportGet)).Methods("GET")
	router.HandleFunc("/support", checkHeaders(SupportAdd)).Methods("POST")
	router.HandleFunc("/support", checkHeaders(SupportUpdate)).Methods("PUT")

	// user
	router.HandleFunc("/user", checkHeaders(UserGet)).Queries(
		"hostel_id", "{hostel_id}",
	).Methods("GET")
	router.HandleFunc("/user", checkHeaders(UserAdd)).Methods("POST")
	router.HandleFunc("/user", checkHeaders(UserUpdate)).Queries(
		"id", "{id}",
		"hostel_id", "{hostel_id}",
	).Methods("PUT")
	router.HandleFunc("/user", checkHeaders(UserDelete)).Queries(
		"id", "{id}",
		"hostel_id", "{hostel_id}",
	).Methods("DELETE")
	router.HandleFunc("/userbook", checkHeaders(UserJoin)).Queries(
		"hostel_id", "{hostel_id}",
	).Methods("PUT")
	router.HandleFunc("/userbooked", checkHeaders(UserJoined)).Queries(
		"hostel_id", "{hostel_id}",
	).Methods("PUT")
	router.HandleFunc("/uservacate", checkHeaders(UserVacate)).Queries(
		"hostel_id", "{hostel_id}",
	).Methods("PUT")

	router.HandleFunc("/upload", checkHeaders(Upload)).Methods("POST")

	router.Path("/").HandlerFunc(HealthCheck).Methods("GET")

	// fmt.Println(http.ListenAndServe(":5000", &WithCORS{router}))

	algnhsa.ListenAndServe(router, nil)
}

func (s *WithCORS) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS,POST,PUT,DELETE")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type,apikey,appversion")

	if req.Method == "OPTIONS" {
		return
	}

	s.r.ServeHTTP(res, req)
}
