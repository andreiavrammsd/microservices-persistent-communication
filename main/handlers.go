package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"io"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	var service Service
	
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 100000))
	checkError(err)
	
	e := r.Body.Close()
	checkError(e)
	
	if err := json.Unmarshal(body, &service); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		err := json.NewEncoder(w).Encode(err)
		checkError(err)
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(service)
}
