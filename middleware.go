package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(config.AuthorizationKey) > 0 {
			s := r.Header.Get(config.AuthorizationHeader)
			if len(s) == 0 {
				notAuthorizedHandler(w, r)
				return
			}

			b, err := base64.StdEncoding.DecodeString(s)
			if err != nil {
				notAuthorizedHandler(w, r)
				return
			}

			pair := strings.SplitN(string(b), ":", 2)
			if len(pair) != 2 {
				notAuthorizedHandler(w, r)
				return
			}

			h := sha256.New()
			h.Write([]byte(fmt.Sprintf("%s:%s", pair[1], config.AuthorizationKey)))
			key := hex.EncodeToString(h.Sum(nil))

			if key != pair[0] {
				notAuthorizedHandler(w, r)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func notAuthorizedHandler(w http.ResponseWriter, r *http.Request) {
	ip := strings.Trim(fmt.Sprintf("%s %s", r.RemoteAddr, r.Header.Get("X-Forwarded-For")), " ")
	log.Printf("Unauthorized request: %s, %s", ip, r.Header.Get(config.AuthorizationHeader))

	response := NewResponse(w)
	response.Body.Error = true
	response.Body.Message = http.StatusText(http.StatusUnauthorized)
	response.Status = http.StatusUnauthorized
	response.Headers["X-Content-Type-Options"] = "nosniff"
	response.Write()
}
