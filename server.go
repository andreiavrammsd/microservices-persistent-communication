package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"strings"
)

func StartServer(config ServerConfig, router *mux.Router) {
	if config.Tls {
		if config.RedirectToTls {
			go func() {
				log.Fatal(http.ListenAndServe(config.Address, http.HandlerFunc(redirect)))
			}()
		}

		log.Fatal(http.ListenAndServeTLS(
			config.AddressTls,
			config.CertFile,
			config.KeyFile, router,
		))
	} else {
		log.Fatal(http.ListenAndServe(config.Address, router))
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	redirectUrl := "https://" +
	strings.Replace(r.Host, config.Server.Address, config.Server.AddressTls, 1) +
	r.URL.String()
	http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
}
