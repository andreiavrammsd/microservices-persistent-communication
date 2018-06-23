package main

import (
	"io"
	"io/ioutil"
	"net/http"
)

type Request struct {
	Request *http.Request
}

func (r *Request) GetBody() []byte {
	body, err := ioutil.ReadAll(io.LimitReader(r.Request.Body, 100000))
	checkError(err)

	e := r.Request.Body.Close()
	checkError(e)

	return body
}

func NewRequest(r *http.Request) *Request {
	return &Request{
		Request: r,
	}
}
