package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"bytes"
	"log"
)

type Service struct {
	Url    string `json:"url"`
	Method string `json:"method"`
	Body   string `json:"body"`
}

func (s *Service) Validate() error {
	if len(s.Url) == 0 {
		return errors.New("Invalid service")
	}
	
	return nil
}

func (s *Service) Call() bool {
	body := []byte(s.Body)
	req, _ := http.NewRequest(s.Method, s.Url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode <= 299
}

func NewService(s []byte) (*Service, error) {
	var service *Service
	var err error
	err = json.Unmarshal(s, &service)
	return service, err
}
