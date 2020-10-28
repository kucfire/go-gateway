package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

const (
	Addr = ":1210"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

func main() {
	// createServe()
	createServer2()
}

func createServe() {
	// create router
	mux := http.NewServeMux()

	fh := HandlerFunc(sayBye)

	// setting rule of router
	mux.HandleFunc("/bye", fh.ServeHTTP)
	// create a server
	server := &http.Server{
		Addr:         Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	// listen the port and support serve
	fmt.Println("Starting httpserver at" + Addr)
	log.Fatal(server.ListenAndServe())
}

func createServer2() {
	hf := HandlerFunc(sayBye)

	//
	resp := httptest.NewRecorder()
	res := httptest.NewRequest("GET", "/", bytes.NewBuffer([]byte("test")))
	hf.ServeHTTP(resp, res)
	bts, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bts))
}

func sayBye(w http.ResponseWriter, h *http.Request) {
	time.Sleep(time.Second)
	w.Write([]byte("bye bye, this is httpServer"))
}
