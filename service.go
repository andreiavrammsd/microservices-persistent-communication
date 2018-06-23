package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Service struct {
	Url     string                 `json:"url" validate:"validurl"`
	Method  string                 `json:"method" validate:"httpmethod"`
	Body    string                 `json:"body" validate:"validbody"`
	Headers map[string]interface{} `json:"headers"`
	Requeue bool                   `json:"requeue"`
}

func (s *Service) Call() bool {
	body := []byte(s.Body)
	req, _ := http.NewRequest(s.Method, s.Url, bytes.NewBuffer(body))

	for key, value := range s.Headers {
		stringValue := fmt.Sprintf("%v", value)
		req.Header.Set(key, stringValue)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer resp.Body.Close()

	return s.requestIsValid(resp)
}

func (s *Service) requestIsValid(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode <= 299
}

func NewService(s []byte) (Service, error) {
	service := Service{
		Requeue: true,
		Method:  "GET",
	}
	err := json.Unmarshal(s, &service)
	service.Method = strings.ToUpper(service.Method)
	return service, err
}
