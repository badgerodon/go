package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

type (
	FrontEnd struct {
		backend *BackEnd
		proxy *httputil.ReverseProxy
		http net.Listener
		https net.Listener
	}

	FrontEndConfigEntry struct {
		Domain string `json:"domain"`		
		Ssl struct {
			Certificate string `json:"certificate"`
			Key         string `json:"key"`
		} `json:"ssl"`
	}
	FrontEndConfig []FrontEndConfigEntry
)

func NewFrontEnd(backend *BackEnd) *FrontEnd {
	this := &FrontEnd{
		backend: backend,
	}
	this.proxy = &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			this.Director(req)
		},
	}
	return this
}

func (this *FrontEnd) LoadFile(filename string) error {
	log.Println("Loading frontend config from", filename)
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	var cfg FrontEndConfig
	err = json.Unmarshal(contents, &cfg)
	if err != nil {
		return err
	}
	return this.Load(cfg)	
}

func (this *FrontEnd) Load(config FrontEndConfig) error {
	tcfg := &tls.Config{}
	tcfg.NextProtos = []string{"http/1.1"}

	var err error
	tcfg.Certificates = make([]tls.Certificate, 0)
	for _, entry := range config {
		if entry.Ssl.Certificate != "" && entry.Ssl.Key != "" {
			cert, err := tls.LoadX509KeyPair(
				entry.Ssl.Certificate,
				entry.Ssl.Key,
			)
			if err != nil {
				return err
			}
			tcfg.Certificates = append(tcfg.Certificates, cert)
		}
	}

	c1, err := net.Listen("tcp", ":https")
	if err != nil {
		return err
	}
	c2, err := net.Listen("tcp", ":http")
	if err != nil {
		return err
	}

	this.https = tls.NewListener(c1, tcfg)
	this.http = c2
	return nil
}
func (this *FrontEnd) Director(req *http.Request) {
	log.Println(req.Host)
	server, found := this.backend.Get(req.Host)
	if !found {
		log.Print("No Server Found")
	}
	req.URL.Scheme = "http"
	req.URL.Host = fmt.Sprint(server.Host, ":", server.Port)
}
func (this *FrontEnd) Run() error {
	err := make(chan error)
	go func() {
		err <- http.Serve(this.https, this.proxy)
	}()
	go func() {
		err <- http.Serve(this.http, this.proxy)
	}()
	return <- err
}
func (this *FrontEnd) Close() {
	this.http.Close()
	this.https.Close()
}