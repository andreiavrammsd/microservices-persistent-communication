package main

import (
	"encoding/json"
	"encoding/xml"
	"net/url"

	"gopkg.in/go-playground/validator.v9"
)

type Empty struct {
}

func NewValidate(config ValidationConfig) *validator.Validate {
	validate := validator.New()

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

		found := false
		for _, value := range config.Protocols {
			if u.Scheme == value {
				found = true
				break
			}
		}
		return found
	})

	validate.RegisterValidation("httpmethod", func(f validator.FieldLevel) bool {
		method := f.Field().String()
		found := false

		for _, value := range config.Methods {
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
