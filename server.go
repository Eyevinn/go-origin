package main

import (
	"fmt"
	"net/http"
	"log"
	"path/filepath"
	"os"
	"io"
	"strings"
)

var OriginStoreDir = filepath.FromSlash("./testmedia/")

func main() {
	if os.Getenv("MEDIAPATH") != "" {
		OriginStoreDir = filepath.FromSlash(os.Getenv("MEDIAPATH"))
	}

	log.Println("Starting Eyevinn simple origin store=" + OriginStoreDir)

	http.HandleFunc("/", healthCheckHandler)
	http.HandleFunc("/ingest/", pushHandler)
	http.HandleFunc("/live/", pullHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" || r.Method == "POST" {
		log.Println("PUT: " + r.URL.String())
		destpath := strings.Replace(r.URL.String(), "ingest", "live", 1)
		dir, file := filepath.Split(filepath.Join(OriginStoreDir, destpath))
		if dir != "" {
			os.MkdirAll(dir, os.ModePerm)
		}
		destfile, err := os.Create(filepath.Join(dir, file))
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		n, err := io.Copy(destfile, r.Body)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		log.Printf("%s: %d bytes written.\n", destpath, n)
		w.WriteHeader(204)
	}
}

func pullHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("GET: " + r.URL.String())
		log.Println("Serving file: " + filepath.Join(OriginStoreDir, r.URL.String()))
		http.ServeFile(w, r, filepath.Join(OriginStoreDir, r.URL.String()))
	}
}