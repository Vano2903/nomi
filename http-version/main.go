package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
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

//handle the requests from names endpoint
func namesHandler(w http.ResponseWriter, r *http.Request) {
	//check the method
	if r.Method == http.MethodGet {
		//generate the seed for the random
		rand.Seed(time.Now().UnixNano())
		//set the content type
		w.Header().Set("Content-Type", "application/json")

		//parse the url query as a map
		r.ParseForm()
		//check if n in the map
		if len(r.Form["n"]) <= 0 {
			sendError(w, http.StatusLengthRequired, "must assign n, bounds are between 2 and 100")
			return
		}

		//convert n from the map (string) to int and check for exceptions
		n, err := strconv.Atoi(r.Form["n"][0])
		if err != nil {
			sendError(w, http.StatusLengthRequired, "n must be a number, bounds are between 2 and 100")
			return
		}

		//check the bounds of n
		if n < 2 || n > 100 {
			sendError(w, http.StatusRequestedRangeNotSatisfiable, "n out of bounds, must be between 2 and 100")
			return
		}

		//read the file
		file, err := ioutil.ReadFile(fileNamePath)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "can't access the database")
			return
		}
		//create a slice of strings with all the names (use \n and \r as sep)
		names := strings.Split(strings.Join(strings.Split(string(file), "\n"), ""), "\r")

		//check if start is defined
		if len(r.Form["start"]) > 0 {
			//if defined filter the names slice by all the names that begin with "start" string
			names = filter(names, r.Form["start"][0])
			//check the length of the name slice, if 0 it means that nothing was found
			if len(names) == 0 {
				sendError(w, http.StatusBadRequest, "none of the names start with the prefix you gave")
				return
			}
		}

		//check if the user requested more names that available based from his request
		if n > len(names) {
			sendError(w, http.StatusBadRequest, fmt.Sprintf("requested more names than available (we have %d available names that satisfy your conditions)", len(names)))
			return
		}
		var randomNames []string
		//generate the slice with the random names
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

		//respond with the json
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

//handle the requests from name endpoint
func nameHandler(w http.ResponseWriter, r *http.Request) {
	//most the code is the same as the namesHandler function, check that one out
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

//handle the requests from exist endpoint
func existHandler(w http.ResponseWriter, r *http.Request) {
	//most of the documentation of this code is in namesHandler
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")

		r.ParseForm()
		//check if toSearch is defined
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

	//read port from enviroment, if not found will assing 8080 by default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//listen and serve needs a port and a http handler, in this case there is none
	//becuse we are using the default http package (http.HandleFunc) so we are giving nil (null) as parameter
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
