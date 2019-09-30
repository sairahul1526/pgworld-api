package main

var dbConfig string
var connectionPool int

var awsAccessKey string
var awsSecretKey string

var baseURL string

var supportEmailID string
var supportEmailPassword string
var supportEmailHost string
var supportEmailPort int

var adminTable = "admins"
var billTable = "bills"
var employeeTable = "employees"
var foodTable = "foods"
var invoiceTable = "invoices"
var issueTable = "issues"
var hostelTable = "hostels"
var logTable = "logs"
var noteTable = "notes"
var noticeTable = "notices"
var roomTable = "rooms"
var signupTable = "signups"
var supportTable = "supports"
var userTable = "users"

var dialogType = "1"
var toastType = "2"
var appUpdateAvailable = "3"
var appUpdateRequired = "4"

var androidLive = "T9h9P6j2N6y9M3Q8"
var androidTest = "K7b3V4h3C7t6g6M7"
var iOSLive = "b4E6U9K8j6b5E9W3"
var iOSTest = "R4n7N8G4m9B4S5n2"

// for checking unauth request
var apikeys = map[string]string{
	androidLive: "1", // android live
	androidTest: "1", // android test
	iOSLive:     "1", // iOS live
	iOSTest:     "1", // iOS test
}

// required fields
var adminRequiredFields = []string{}
var billRequiredFields = []string{}
var employeeRequiredFields = []string{}
var foodRequiredFields = []string{}
var invoiceRequiredFields = []string{}
var issueRequiredFields = []string{}
var hostelRequiredFields = []string{}
var logRequiredFields = []string{}
var noteRequiredFields = []string{}
var noticeRequiredFields = []string{}
var roomRequiredFields = []string{}
var signupRequiredFields = []string{}
var supportRequiredFields = []string{}
var userRequiredFields = []string{}

// server codes
var statusCodeOk = "200"
var statusCodeCreated = "201"
var statusCodeBadRequest = "400"
var statusCodeForbidden = "403"
var statusCodeServerError = "500"
var statusCodeDuplicateEntry = "1000"

var defaultLimit = "25"
var defaultOffset = "0"

var test bool
var migrate bool

// versions
var iOSVersionCode = 1.0
var iOSForceVersionCode = 1.0

var androidVersionCode = 1.6
var androidForceVersionCode = 1.6

// s3
var s3Bucket string
var docS3Path = "document"
