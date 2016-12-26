package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"io"
	"fmt"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type Response struct {
	Error bool `json:"error"`
	Message string `json:"message"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	var serviceRequest ServiceRequest
	var response Response
	var status int
	
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 100000))
	checkError(err)
	
	e := r.Body.Close()
	checkError(e)
	
	if err := json.Unmarshal(body, &serviceRequest); err != nil {
		response = Response{
			Error: true,
			Message: err.Error(),
		}
		status = 422
	} else {
		if serviceIsDefined(serviceRequest.Name) {
			message, _ := json.Marshal(serviceRequest)
			NewQueue(config.QueueName).Publish(string(message))

			response = Response{
				Error: false,
				Message: "Success",
			}
			status = 201
		} else {
			response = Response{
				Error: true,
				Message: fmt.Sprintf("Service \"%s\" is not defined.", serviceRequest.Name),
			}
			status = 400
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
