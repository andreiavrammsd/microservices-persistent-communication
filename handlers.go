package main

import (
	"net/http"
	"log"
)

func Index(w http.ResponseWriter, r *http.Request) {
	body := NewRequest(r).GetBody()
	service, err := NewService(body)
	response := NewResponse(w)
	
	if err != nil {
		response.Body.Error = true
		response.Body.Message = err.Error()
		response.Status = http.StatusUnprocessableEntity
		log.Printf(err.Error())
	} else {
		err := validate.Struct(service)
		if err == nil {
			servicesQueue.Publish(body)
			response.Body.Message = "Success"
			response.Status = http.StatusCreated
			log.Printf("Received: %s", string(body))
		} else {
			response.Body.Error = true
			response.Body.Message = err.Error()
			response.Status = http.StatusBadRequest
			log.Printf(err.Error())
		}
	}

	response.Write()
}
