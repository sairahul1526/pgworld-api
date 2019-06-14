package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func logAction(by string, logData string) {
	// go func(by string, logData string) {
	// 	insertSQL(logTable, map[string]string{"log": logData, "by": by})
	// }(by, logData)
}

func logger(str interface{}) {
	if test {
		fmt.Println(str)
	}
}

func sqlErrorCheck(code uint16) string {
	if code == 1054 { // Error 1054: Unknown column
		return statusCodeBadRequest
	} else if code == 1062 { // Error 1062: Duplicate entry
		return statusCodeDuplicateEntry
	}
	return statusCodeServerError
}
func setMeta(status string, msg string, msgType string) map[string]string {
	if len(msg) == 0 {
		if status == statusCodeBadRequest {
			msg = "Bad Request Body"
		} else if status == statusCodeServerError {
			msg = "Internal Server Error"
		}
	}
	return map[string]string{
		"status":       status,
		"message":      msg,
		"message_type": msgType, // 1 : dialog or 2 : toast if msg
	}
}

func getHTTPStatusCode(code string) int {
	switch code {
	case statusCodeOk:
		return http.StatusOK
	case statusCodeCreated:
		return http.StatusCreated
	case statusCodeBadRequest:
		return http.StatusBadRequest
	case statusCodeServerError:
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

// SetReponseStatus .
func SetReponseStatus(w http.ResponseWriter, r *http.Request, status string, msg string, msgType string, response map[string]interface{}) {
	w.Header().Set("Status", status)
	response["meta"] = setMeta(status, msg, msgType)
	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	json.NewEncoder(w).Encode(response)
}

func checkAppUpdate(r *http.Request) (map[string]string, bool) {
	if strings.EqualFold(r.Header.Get("apikey"), androidLive) || strings.EqualFold(r.Header.Get("apikey"), androidTest) {
		appversion, _ := strconv.ParseFloat(r.Header.Get("appversion"), 64)
		if appversion < androidForceVersionCode {
			return setMeta(statusCodeOk, "App update required", appUpdateRequired), true
		} else if appversion < androidVersionCode {
			return setMeta(statusCodeOk, "App update available", appUpdateAvailable), true
		}
	} else if strings.EqualFold(r.Header.Get("apikey"), iOSLive) || strings.EqualFold(r.Header.Get("apikey"), iOSTest) {
		appversion, _ := strconv.ParseFloat(r.Header.Get("appversion"), 64)
		if appversion < iOSForceVersionCode {
			return setMeta(statusCodeOk, "App update required", appUpdateRequired), true
		} else if appversion < iOSVersionCode {
			return setMeta(statusCodeOk, "App update available", appUpdateAvailable), true
		}
	}
	return map[string]string{}, false
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

func checkHeaders(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var response = make(map[string]interface{})
		if len(r.Header.Get("apikey")) == 0 || len(r.Header.Get("appversion")) == 0 {
			SetReponseStatus(w, r, statusCodeBadRequest, "apikey, appversion required", "", response)
			return
		} else if len(apikeys[r.Header.Get("apikey")]) == 0 {
			SetReponseStatus(w, r, statusCodeBadRequest, "Unauthorized request. Not valid apikey", "", response)
			return
		}

		if migrate { // statusCodeBadRequest because app will hit again if 500
			SetReponseStatus(w, r, statusCodeBadRequest, "Server is busy. Please try after some time", dialogType, response)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func requiredFiledsCheck(body map[string]string, required []string) string {
	for _, field := range required {
		if len(body[field]) == 0 {
			return field
		}
	}
	return ""
}

// RandStringBytes .
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func getMD5Hash(savedFileName string) string {
	openedFile, _ := os.Open(savedFileName)
	defer openedFile.Close()

	h := md5.New()
	if _, err := io.Copy(h, openedFile); err != nil {
		fmt.Println("getMD5Hash", err)
	}

	return hex.EncodeToString(h.Sum(nil))
}

func uploadToS3(path string, file multipart.File, extension string, maxage ...string) (string, bool) {

	savedFileName, saved := saveToDisk(file, extension)
	if !saved {
		return "", false
	}

	// conf := aws.Config{Region: aws.String("ap-south-1")}
	// sess := session.New(&conf)

	// svc := s3manager.NewUploader(sess)

	logger("Uploading file to S3...")

	openedFile, err := os.Open(savedFileName)
	if err != nil {
		fmt.Println("uploadToS3", err)
		return "", false
	}
	defer openedFile.Close()

	fileName := path + "/" + getMD5Hash(savedFileName) + extension

	// _, err = svc.Upload(&s3manager.UploadInput{
	// 	Bucket:      aws.String(s3Bucket),
	// 	Key:         aws.String(fileName),
	// 	Body:        openedFile,
	// 	ContentType: aws.String(getFileMIMEType(strings.ToLower(extension))),
	// 	ACL:         aws.String("public-read"),
	// })

	// os.Remove(savedFileName)
	// if err != nil {
	// 	fmt.Println("uploadToS3 "+path, err)
	// 	return "", false
	// }
	return fileName, true
}

func saveToDisk(file multipart.File, extension string) (string, bool) {
	logger("Saving to disk...")

	fileName := "./" + RandStringBytes(10) + extension

	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println("saveToDisk", err)
		return "", false
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println("saveToDisk", err)
		return "", false
	}

	return fileName, true
}

func getFileMIMEType(extension string) string {
	switch extension {
	// video
	case ".mp2":
		return "video/mpeg"
	case ".mpa":
		return "video/mpeg"
	case ".mpe":
		return "video/mpeg"
	case ".mpeg":
		return "video/mpeg"
	case ".mpg":
		return "video/mpeg"
	case ".mpv2":
		return "video/mpeg"
	case ".mp4":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	case ".qt":
		return "video/quicktime"
	case ".lsf":
		return "video/x-la-asf"
	case ".lsx":
		return "video/x-la-asf"
	case ".asf":
		return "video/x-ms-asf"
	case ".asr":
		return "video/x-ms-asf"
	case ".asx":
		return "video/x-ms-asf"
	case ".avi":
		return "video/x-msvideo"
	case ".movie":
		return "video/x-sgi-movie"
	case ".3gp":
		return "video/3gpp"
	case ".3gpp":
		return "video/3gpp"
	case ".3gpp2":
		return "video/3gpp2"
	case ".3g2":
		return "video/3gpp2"
	// image
	case ".bmp":
		return "image/bmp"
	case ".cod":
		return "image/cis-cod"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".ief":
		return "image/ief"
	case ".jpe":
		return "image/jpeg"
	case ".jpeg":
		return "image/jpeg"
	case ".jpg":
		return "image/jpeg"
	case ".jfif":
		return "image/pipeg"
	case ".svg":
		return "image/svg+xml"
	case ".tif":
		return "image/tiff"
	case ".tiff":
		return "image/tiff"
	case ".ras":
		return "image/x-cmu-raster"
	case ".cmx":
		return "image/x-cmx"
	case ".ico":
		return "image/x-icon"
	case ".pnm":
		return "image/x-portable-anymap"
	case ".pbm":
		return "image/x-portable-bitmap"
	case ".pgm":
		return "image/x-portable-graymap"
	case ".ppm":
		return "image/x-portable-pixmap"
	case ".rgb":
		return "image/x-rgb"
	case ".xbm":
		return "image/x-xbitmap"
	case ".xpm":
		return "image/x-xpixmap"
	case ".xwd":
		return "image/x-xwindowdump"
	default:
		return ""
	}
}
