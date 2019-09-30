package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func connectDatabase() {
	var err error
	db, err = sql.Open("mysql", dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(connectionPool)
	db.SetMaxIdleConns(connectionPool)
	db.SetConnMaxLifetime(time.Hour)
}

// HealthCheck .
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("ok")
}

func inits() {

	os.Setenv("AWS_ACCESS_KEY_ID", awsAccessKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", awsSecretKey)
	os.Setenv("AWS_REGION", "ap-south-1")

	rand.Seed(time.Now().UnixNano())

	connectDatabase()
	initcache()
}

func main() {

	// root:root@tcp(localhost:8889)/testDB?charset=utf8mb4
	// pgcruxapp:pg4c!123@tcp(pgcruxapp.clsfriejtsvw.ap-south-1.rds.amazonaws.com:3306)/pgcruxapp?charset=utf8mb4
	// dbConfig = "pgcruxapp:pg4c!123@tcp(pgcruxapp.clsfriejtsvw.ap-south-1.rds.amazonaws.com:3306)/pgcruxapp?charset=utf8mb4"
	// connectionPool = 10
	// test = true
	// migrate = false
	// awsAccessKey = "AKIAUMOCANH676PVBR7X"
	// awsSecretKey = "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"
	// s3Bucket = "test-pgworld"
	// baseURL = "https://test-pgworld.s3.ap-south-1.amazonaws.com/"
	// supportEmailID = "support@cloudpg.in"
	// supportEmailPassword = "sp90rt!pSs#ds"
	// supportEmailHost = "smtp.zoho.com"
	// supportEmailPort = 587
	dbConfig = os.Getenv("dbConfig")
	connectionPool, _ = strconv.Atoi(os.Getenv("connectionPool"))
	test, _ = strconv.ParseBool(os.Getenv("test"))
	migrate, _ = strconv.ParseBool(os.Getenv("migrate"))
	awsAccessKey = os.Getenv("awsAccessKey")
	awsSecretKey = os.Getenv("awsSecretKey")
	s3Bucket = os.Getenv("s3Bucket")
	baseURL = os.Getenv("baseURL")
	supportEmailID = os.Getenv("supportEmailID")
	supportEmailPassword = os.Getenv("supportEmailPassword")
	supportEmailHost = os.Getenv("supportEmailHost")
	supportEmailPort, _ = strconv.Atoi(os.Getenv("supportEmailPort"))

	inits()
	defer db.Close()
	router := mux.NewRouter()

	// dashboard
	router.HandleFunc("/dashboard", checkHeaders(Dashboard)).Methods("GET")

	// otp
	router.HandleFunc("/sendotp", checkHeaders(SendOTP)).Methods("GET")
	router.HandleFunc("/verifyotp", checkHeaders(VerifyOTP)).Methods("GET")

	// report
	router.HandleFunc("/report", checkHeaders(Report)).Methods("GET")

	// admin
	router.HandleFunc("/admin", checkHeaders(AdminGet)).Methods("GET")
	router.HandleFunc("/admin", checkHeaders(AdminAdd)).Methods("POST")
	router.HandleFunc("/admin", checkHeaders(AdminUpdate)).Methods("PUT")

	// bill
	router.HandleFunc("/bill", checkHeaders(BillGet)).Methods("GET")
	router.HandleFunc("/bill", checkHeaders(BillAdd)).Methods("POST")
	router.HandleFunc("/bill", checkHeaders(BillUpdate)).Methods("PUT")

	// employee
	router.HandleFunc("/employee", checkHeaders(EmployeeGet)).Methods("GET")
	router.HandleFunc("/employee", checkHeaders(EmployeeAdd)).Methods("POST")
	router.HandleFunc("/employee", checkHeaders(EmployeeUpdate)).Methods("PUT")

	// invoice
	router.HandleFunc("/food", checkHeaders(FoodGet)).Methods("GET")
	router.HandleFunc("/food", checkHeaders(FoodAdd)).Methods("POST")
	router.HandleFunc("/food", checkHeaders(FoodUpdate)).Methods("PUT")

	// invoice
	router.HandleFunc("/invoice", checkHeaders(InvoiceGet)).Methods("GET")
	router.HandleFunc("/invoice", checkHeaders(InvoiceAdd)).Methods("POST")
	router.HandleFunc("/invoice", checkHeaders(InvoiceUpdate)).Methods("PUT")

	// issue
	router.HandleFunc("/issue", checkHeaders(IssueGet)).Methods("GET")
	router.HandleFunc("/issue", checkHeaders(IssueAdd)).Methods("POST")
	router.HandleFunc("/issue", checkHeaders(IssueUpdate)).Methods("PUT")

	// hostel
	router.HandleFunc("/hostel", checkHeaders(HostelGet)).Methods("GET")
	router.HandleFunc("/hostel", checkHeaders(HostelAdd)).Methods("POST")
	router.HandleFunc("/hostel", checkHeaders(HostelUpdate)).Methods("PUT")

	// log
	router.HandleFunc("/log", checkHeaders(LogGet)).Methods("GET")

	// note
	router.HandleFunc("/note", checkHeaders(NoteGet)).Methods("GET")
	router.HandleFunc("/note", checkHeaders(NoteAdd)).Methods("POST")
	router.HandleFunc("/note", checkHeaders(NoteUpdate)).Methods("PUT")

	// notice
	router.HandleFunc("/notice", checkHeaders(NoticeGet)).Methods("GET")
	router.HandleFunc("/notice", checkHeaders(NoticeAdd)).Methods("POST")
	router.HandleFunc("/notice", checkHeaders(NoticeUpdate)).Methods("PUT")

	// room
	router.HandleFunc("/room", checkHeaders(RoomGet)).Methods("GET")
	router.HandleFunc("/room", checkHeaders(RoomAdd)).Methods("POST")
	router.HandleFunc("/room", checkHeaders(RoomUpdate)).Methods("PUT")

	// set hostel status
	router.HandleFunc("/hostelstatus", HostelStatus).Methods("GET")

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
	router.HandleFunc("/user", checkHeaders(UserGet)).Methods("GET")
	router.HandleFunc("/user", checkHeaders(UserAdd)).Methods("POST")
	router.HandleFunc("/user", checkHeaders(UserUpdate)).Methods("PUT")
	router.HandleFunc("/user", checkHeaders(UserDelete)).Methods("DELETE")
	router.HandleFunc("/userbook", checkHeaders(UserJoin)).Methods("PUT")
	router.HandleFunc("/userbooked", checkHeaders(UserJoined)).Methods("PUT")
	router.HandleFunc("/uservacate", checkHeaders(UserVacate)).Methods("PUT")

	router.HandleFunc("/upload", checkHeaders(Upload)).Methods("POST")

	router.Path("/").HandlerFunc(HealthCheck).Methods("GET")

	fmt.Println(http.ListenAndServe(":5000", &WithCORS{router}))
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
