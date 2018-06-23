package main

import (
	"encoding/json"
	"encoding/xml"
	"net/url"

	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

type empty struct {
}

var customValidations = map[string]validator.Func{
	"validurl": func(f validator.FieldLevel) bool {
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
		for _, value := range config.Validation.Protocols {
			if u.Scheme == value {
				found = true
				break
			}
		}
		return found
	},

	"httpmethod": func(f validator.FieldLevel) bool {
		method := f.Field().String()
		found := false

		for _, value := range config.Validation.Methods {
			if value == method {
				found = true
				break
			}
		}

		return found
	},

	"validbody": func(f validator.FieldLevel) bool {
		text := f.Field().String()

		if len(text) == 0 {
			return true
		}

		v := &empty{}
		valid := true

		if err := json.Unmarshal([]byte(text), v); err != nil {
			if err := xml.Unmarshal([]byte(text), v); err != nil {
				valid = false
			}
		}

		return valid
	},
}

func NewValidate() (*validator.Validate, error) {
	v := validator.New()

	for tag, fn := range customValidations {
		err := v.RegisterValidation(tag, fn)

		if err != nil {
			return nil, fmt.Errorf("error at validation init (%s)", err)
		}
	}

	return v, nil
}
