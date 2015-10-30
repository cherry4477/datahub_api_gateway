package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"
)

func init() {
	_ = os.Getenv("Test")
}

//======================================================
//
//======================================================

const (
	ErrorCode_Unkown = 1
	ErrorCode_JsonCreateError = 100
	
	MaxErrorCodePlusOne = 256
)

var (
	ErrorMessages = make ([]string, MaxErrorCodePlusOne)
)

func init () {
	buildErrorMessage (ErrorCode_Unkown, "Unknown error")
	buildErrorMessage (ErrorCode_JsonCreateError, "Error on creating json")
}

func buildErrorMessage (code int, message string) {
	ErrorMessages [code] = fmt.Sprintf ("%d: %s", code, message)
}

//======================================================
//
//======================================================

type Result struct {
	Ok    bool        `json:"ok"`
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

var (
	Json_JsonError string
)

func init () {
	Json_JsonError = fmt.Sprintf (`{"ok":false,"error":"%s"}`, ErrorMessages [ErrorCode_JsonCreateError])
}

//======================================================
//
//======================================================

func jsonResult(errorMessage string, data interface{}) []byte {
	result := &Result{}
	
	if errorMessage == "" && data == nil {
		errorMessage = ErrorMessages [ErrorCode_Unkown]
	}
	
	if errorMessage != "" {
		result.Ok = false
		result.Error = errorMessage
	} else {
		result.Ok = true
		result.Data = data
	}

	jsondata, err := json.Marshal(&result)
	if err != nil {
		return []byte(Json_JsonError)
	} else {
		return jsondata
	}
}

//======================================================
//
//======================================================

func doLogging () {
   
}

//======================================================
//
//======================================================

func onSinaAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsonResult("Unsupported URL", nil))
}

func onSohuAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsonResult("Unsupported URL", nil))
}

func onServiceError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsonResult("Unsupported URL", nil))
}

//======================================================
//
//======================================================

var port = flag.Int("port", 7777, "server port")

func main() {
   go doLogging ()
   
	http.HandleFunc("/sina", onSinaAPI)
	http.HandleFunc("/sohu", onSohuAPI)
	http.HandleFunc("/", onServiceError)

	address := fmt.Sprintf(":%d", *port)
	log.Printf("Listening at %s\n", address)
	log.Fatal(http.ListenAndServe(address, nil)) // will block here
}
