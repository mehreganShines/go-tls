//simple HTTP server that uses TLS(also known as HTTPS) to secure the connection and simply logs the user agent when the user clicked the html button
//to generate certificates run "go-tls/self-cert.go" and save the results in cert.pem and key.pem files.
package main

import (
	"crypto/tls"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("User Agent: %s\n", r.UserAgent())
		w.Write([]byte("User Agent logged"))
	})

	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("Error loading cert and key: %v", err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
		ClientAuth:   tls.NoClientCert,
		CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256},
	}

	srv := &http.Server{
		Addr:      ":8443",
		TLSConfig: config,
	}

	log.Println("Listening on port 8443...")
	err = srv.ListenAndServeTLS("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
