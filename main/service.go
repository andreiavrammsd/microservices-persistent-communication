package main

import (
	"encoding/json"
	"fmt"
	"errors"
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

func (s *Service) Call() {
	fmt.Println(s.Url)
}

func NewService(s []byte) (*Service, error) {
	var service *Service
	var err error
	err = json.Unmarshal(s, &service)
	return service, err
}
