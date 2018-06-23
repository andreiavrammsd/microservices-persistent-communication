package main

import (
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	body := NewRequest(r).GetBody()

	if config.FastPublish {
		if len(body) == 0 {
			log.Println("could not read request body")
			return
		}

		if err := serviceQueue.Publish(body); err != nil {
			log.Println(err)
		}
		return
	}

	service, err := NewService(body)
	response := NewResponse(w)

	if err != nil {
		response.Body.Error = true
		response.Body.Message = err.Error()
		response.Status = http.StatusUnprocessableEntity
		log.Println(err.Error())
	} else {
		err := validate.Struct(service)
		if err == nil {
			if err := serviceQueue.Publish(body); err != nil {
				response.Body.Error = true
				response.Body.Message = http.StatusText(http.StatusInternalServerError)
				response.Status = http.StatusInternalServerError
				log.Printf("Error publishing: %s", err)
			} else {
				response.Body.Message = http.StatusText(http.StatusAccepted)
				response.Status = http.StatusAccepted
				log.Printf("Accepted: %s", string(body))
			}
		} else {
			response.Body.Error = true
			response.Body.Message = err.Error()
			response.Status = http.StatusBadRequest
			log.Println(err.Error())
		}
	}

	response.Write()
}
