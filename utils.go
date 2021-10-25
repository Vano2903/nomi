package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Resp struct {
	Error   bool        `json:"error, omitempty"`
	Content interface{} `json:"content, omitempty"`
	Msg     string      `json:"msg, omitempty"`
	Code    int         `json:"code, omitempty"`
}

var fileNamePath string = "nomi.txt"

//will respond to the client with a json that explains the error
func sendError(w http.ResponseWriter, code int, msg string) {
	var resp Resp
	//set the status code in the header
	w.WriteHeader(code)
	resp.Error = true
	resp.Code = code
	resp.Msg = msg
	//convert the resp variable in json
	toSend, _ := json.Marshal(resp)
	w.Write(toSend)
}

//filter will filter a slice of string and return a slice with only the strings
//that starts with the given substring
func filter(array []string, sub string) []string {
	var subArray []string
	for _, s := range array {
		if strings.HasPrefix(s, sub) {
			subArray = append(subArray, s)
		}
	}
	return subArray
}

//will check if the string is inside the slice given as param
//true if found, false if not
func inSlice(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
