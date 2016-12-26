package main

import (
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	body := NewRequest(r).GetBody()
	service, err := NewService(body)
	response := NewResponse(w)
	
	if err != nil {
		response.Body.Error = true
		response.Body.Message = err.Error()
		response.Status = http.StatusUnprocessableEntity
	} else {
		validationError := service.Validate()
		if validationError == nil {
			servicesQueue.Publish(body)
			response.Body.Message = "Success"
			response.Status = http.StatusCreated
		} else {
			response.Body.Error = true
			response.Body.Message = validationError.Error()
			response.Status = http.StatusBadRequest
		}
	}

	response.Write()
}
