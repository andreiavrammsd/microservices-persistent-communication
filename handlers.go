package main

import (
	"net/http"
	"log"
)

func Index(w http.ResponseWriter, r *http.Request) {
	body := NewRequest(r).GetBody()

	if config.FastPublish {
		if len(body) > 0 {
			serviceQueue.Publish(body)
		}
		return
	}

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
			serviceQueue.Publish(body)

			response.Body.Message = http.StatusText(http.StatusAccepted)
			response.Status = http.StatusAccepted
			log.Printf("Accepted: %s", string(body))
		} else {
			response.Body.Error = true
			response.Body.Message = err.Error()
			response.Status = http.StatusBadRequest
			log.Printf(err.Error())
		}
	}

	response.Write()
}
