package main

import (
	"gopkg.in/go-playground/validator.v9"
	"encoding/json"
	"net/url"
	"encoding/xml"
)

type Empty struct {
}

func NewValidate() *validator.Validate {
	validate = validator.New()

	validate.RegisterValidation("validurl", func(f validator.FieldLevel) bool {
		rawurl := f.Field().String()

		_, err := url.ParseRequestURI(rawurl)
		if err != nil {
			return false
		}

		u, err := url.Parse(rawurl)
		if err != nil {
			return false
		}

		if len(u.Host) == 0 {
			return false
		}

		protocols := [9]string{"http", "https"}

		found := false
		for _, value := range protocols {
			if u.Scheme == value {
				found = true
				break
			}
		}
		return found
	})

	validate.RegisterValidation("httpmethod", func(f validator.FieldLevel) bool {
		method := f.Field().String()
		methods := [9]string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "CONNECT", "TRACE"}
		found := false

		for _, value := range methods {
			if value == method {
				found = true
				break
			}
		}

		return found
	})

	validate.RegisterValidation("validbody", func(f validator.FieldLevel) bool {
		text := f.Field().String()

		if len(text) == 0 {
			return true
		}

		v := &Empty{}
		errJson := json.Unmarshal([]byte(text), v)
		errXml := xml.Unmarshal([]byte(text), v)

		return errJson == nil || errXml == nil
	})

	return validate
}
