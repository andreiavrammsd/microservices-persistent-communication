package main

type Service struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Method string `json:"method"`
	Body   string `json:"body"`
}
