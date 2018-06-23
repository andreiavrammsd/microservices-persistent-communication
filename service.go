package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Service struct {
	URL     string                 `json:"url" validate:"validurl"`
	Method  string                 `json:"method" validate:"httpmethod"`
	Body    string                 `json:"body" validate:"validbody"`
	Headers map[string]interface{} `json:"headers"`
	Requeue bool                   `json:"requeue"`
}

func (s *Service) Call() (err error) {
	body := []byte(s.Body)
	req, err := http.NewRequest(s.Method, s.URL, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	for key, value := range s.Headers {
		stringValue := fmt.Sprintf("%v", value)
		req.Header.Set(key, stringValue)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer func() {
		if e := resp.Body.Close(); e != nil {
			err = e
		}
	}()

	if !s.requestIsValid(resp) {
		return errors.New("error status code received")
	}

	return
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
