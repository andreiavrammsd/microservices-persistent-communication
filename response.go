package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Body    *responseBody
	Status  int
	Headers map[string]string
	writer  http.ResponseWriter
}

type responseBody struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func (r *Response) Write() {
	for key, value := range r.Headers {
		r.writer.Header().Set(key, value)
	}
	r.writer.WriteHeader(r.Status)

	if err := json.NewEncoder(r.writer).Encode(r.Body); err != nil {
		log.Println(err)
	}
}

func NewResponse(w http.ResponseWriter) *Response {
	response := &Response{
		Body: &responseBody{
			Error:   false,
			Message: "",
		},
		Status:  200,
		Headers: make(map[string]string),
		writer:  w,
	}

	response.Headers["Content-Type"] = "application/json; charset=UTF-8"

	return response
}
