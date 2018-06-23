package main

import (
	"log"
	"net/http"
	"strings"
)

func StartServer(config ServerConfig, router http.Handler) {
	if config.TLS {
		if config.RedirectToTLS {
			go func() {
				log.Fatal(http.ListenAndServe(config.Address, http.HandlerFunc(redirect)))
			}()
		}

		log.Fatal(http.ListenAndServeTLS(
			config.AddressTLS,
			config.CertFile,
			config.KeyFile, router,
		))
	} else {
		log.Fatal(http.ListenAndServe(config.Address, router))
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	redirectURL := "https://" +
		strings.Replace(r.Host, config.Server.Address, config.Server.AddressTLS, 1) +
		r.URL.String()
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
