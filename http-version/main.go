package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Resp struct {
	Error   bool        `json:"error, omitempty"`
	Content interface{} `json:"content, omitempty"`
	Msg     string      `json:"msg, omitempty"`
	Code    int         `json:"code, omitempty"`
}

var fileNamePath string = "nomi.txt"

func sendError(w http.ResponseWriter, code int, msg string) {
	var resp Resp
	w.WriteHeader(code)
	resp.Error = true
	resp.Code = code
	resp.Msg = msg
	toSend, _ := json.Marshal(resp)
	w.Write(toSend)
	return
}

func filter(array []string, sub string) []string {
	var subArray []string
	for _, s := range array {
		if strings.HasPrefix(s, sub) {
			subArray = append(subArray, s)
		}
	}
	return subArray
}

func inSlice(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func namesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		rand.Seed(time.Now().UnixNano())
		w.Header().Set("Content-Type", "application/json")

		r.ParseForm()
		if len(r.Form["n"]) <= 0 {
			sendError(w, http.StatusLengthRequired, "must assign n, bounds are between 2 and 100")
			return
		}

		n, err := strconv.Atoi(r.Form["n"][0])
		if err != nil {
			sendError(w, http.StatusLengthRequired, "must assign n, bounds are between 2 and 100")
			return
		}

		if n < 2 || n > 100 {
			sendError(w, http.StatusRequestedRangeNotSatisfiable, "n out of bounds, must be between 2 and 100")
			return
		}

		file, err := ioutil.ReadFile(fileNamePath)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "can't access the database")
			return
		}
		names := strings.Split(strings.Join(strings.Split(string(file), "\n"), ""), "\r")

		if len(r.Form["start"]) > 0 {
			names = filter(names, r.Form["start"][0])
			if len(names) == 0 {
				sendError(w, http.StatusBadRequest, "none of the names start with the prefix you gave")
				return
			}
		}
		if n > len(names) {
			sendError(w, http.StatusBadRequest, fmt.Sprintf("requested more names than available (we have %d available names that satisfy your conditions)", len(names)))
			return
		}
		var randomNames []string
		for i := 0; i < n; i++ {
			for {
				random := rand.Intn(len(names))
				check := names[random]
				if !inSlice(randomNames, check) {
					randomNames = append(randomNames, check)
					break
				}
			}
		}
		var resp Resp
		resp.Error = false
		resp.Code = http.StatusOK
		resp.Content = randomNames
		toSend, _ := json.Marshal(resp)
		w.Write(toSend)
		return
	}
	sendError(w, http.StatusMethodNotAllowed, "must use get method")
}

func nameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		rand.Seed(time.Now().UnixNano())
		w.Header().Set("Content-Type", "application/json")

		file, err := ioutil.ReadFile(fileNamePath)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "can't access the database")
			return
		}
		names := strings.Split(strings.Join(strings.Split(string(file), "\n"), ""), "\r")

		r.ParseForm()
		if len(r.Form["start"]) > 0 {
			names = filter(names, r.Form["start"][0])
			if len(names) == 0 {
				sendError(w, http.StatusBadRequest, "none of the names start with the prefix you gave")
				return
			}
		}
		var resp Resp
		resp.Error = false
		resp.Code = http.StatusOK
		resp.Content = names[rand.Intn(len(names))]
		toSend, _ := json.Marshal(resp)
		w.Write(toSend)
		return
	}
	sendError(w, http.StatusMethodNotAllowed, "must use get method")
}

func existHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")

		r.ParseForm()
		if len(r.Form["toSearch"]) == 0 {
			sendError(w, http.StatusBadRequest, "'toSearch' must be defined")
			return
		}

		file, err := ioutil.ReadFile(fileNamePath)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "can't access the database")
			return
		}
		names := strings.Split(strings.Join(strings.Split(string(file), "\n"), ""), "\r")

		var resp Resp
		resp.Code = http.StatusOK
		resp.Content = inSlice(names, r.Form["toSearch"][0])
		resp.Error = false
		toSend, _ := json.Marshal(resp)
		w.Write(toSend)
		return
	}
	sendError(w, http.StatusMethodNotAllowed, "must use get method")
}

func main() {
	//handleFunc needs a route and a function that handle the request on that route
	http.HandleFunc("/name", nameHandler)

	http.HandleFunc("/names", namesHandler)

	http.HandleFunc("/exist", existHandler)

	//log fatal kill the program if listenAndServe returns an error

	//listen and serve needs a port and a http handler, in this case there is none
	//becuse we are using the default http package (http.HandleFunc) so we are giving nil (null) as parameter
	log.Fatal(http.ListenAndServe(":8080", nil))
}
