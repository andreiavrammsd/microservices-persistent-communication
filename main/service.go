package main

import (
	"encoding/json"
	"fmt"
	"errors"
	"time"
	"math/rand"
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
	fmt.Println(s.Url)

	rand.Seed(time.Now().Unix())
	return rand.Intn(10 - 1) + 1 <= 5
}

func NewService(s []byte) (*Service, error) {
	var service *Service
	var err error
	err = json.Unmarshal(s, &service)
	return service, err
}
