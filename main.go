package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

func connectDatabase() {
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
}

func main() {

	dbConfig = "root:root@tcp(localhost:8889)/testDB?charset=utf8mb4"
	connectionPool = 10
	test = true
	migrate = false
	awsAccessKey = "asd"
	awsSecretKey = "sda"
	s3Bucket = ""
	// dbConfig = os.Getenv("dbConfig")
	// connectionPool, _ = strconv.Atoi(os.Getenv("connectionPool"))
	// test, _ = strconv.ParseBool(os.Getenv("test"))
	// migrate, _ = strconv.ParseBool(os.Getenv("migrate"))
	// awsAccessKey = os.Getenv("awsAccessKey")
	// awsSecretKey = os.Getenv("awsSecretKey")
	// s3Bucket = os.Getenv("s3Bucket")

	inits()
	defer db.Close()
	router := mux.NewRouter()

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

	// room
	router.HandleFunc("/room", checkHeaders(RoomGet)).Methods("GET")
	router.HandleFunc("/room", checkHeaders(RoomAdd)).Methods("POST")
	router.HandleFunc("/room", checkHeaders(RoomUpdate)).Methods("PUT")

	// transactions
	router.HandleFunc("/rent", checkHeaders(Rent)).Methods("POST")
	router.HandleFunc("/salary", checkHeaders(Salary)).Methods("POST")

	// user
	router.HandleFunc("/user", checkHeaders(UserGet)).Methods("GET")
	router.HandleFunc("/user", checkHeaders(UserAdd)).Methods("POST")
	router.HandleFunc("/user", checkHeaders(UserUpdate)).Methods("PUT")
	router.HandleFunc("/user", checkHeaders(UserDelete)).Methods("DELETE")

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
