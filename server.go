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
	http.HandleFunc("/ingest/", pushAndRemoveHandler)
	http.HandleFunc("/live/", pullHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func pushAndRemoveHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		upload(w, r)
	case "POST":
		upload(w, r)
	case "DELETE":
		delete(w, r)
	}
}

func pullHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("GET: " + r.URL.String())
		log.Println("Serving file: " + filepath.Join(OriginStoreDir, r.URL.String()))
		http.ServeFile(w, r, filepath.Join(OriginStoreDir, r.URL.String()))
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
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
		return
	}
	n, err := io.Copy(destfile, r.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	log.Printf("%s: %d bytes written.\n", destpath, n)
	w.WriteHeader(204)
}

func delete(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE: " + r.URL.String())
	path := strings.Replace(r.URL.String(), "ingest", "live", 1)
	file := filepath.Join(OriginStoreDir, path)
	err := os.Remove(file)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(204)
}