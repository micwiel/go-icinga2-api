// Package icinga2 provides a client for interacting with an Icinga2 Server
package icinga2api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

// Server ... Use to be ClientConfig
type Server struct {
	username           string
	password           string
	URL                string
	AllowUnverifiedSSL bool
	httpClient         *http.Client
}

// func config ...
func (server *Server) Config(username, password, url string, AllowUnverifiedSSL bool) (*Server, error) {

	// TODO : Add code to verify parameters
	return &Server{username, password, url, AllowUnverifiedSSL, nil}, nil

}

func (server *Server) Connect() error {

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: server.AllowUnverifiedSSL,
		},
	}

	server.httpClient = &http.Client{
		Transport: t,
	}

	request, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		server.httpClient = nil
	}

	request.SetBasicAuth(server.username, server.password)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	response, err := server.httpClient.Do(request)
	defer response.Body.Close()

	if (err != nil) || (response.StatusCode != 200) {
		server.httpClient = nil
		fmt.Printf("Failed to connect to %s : %s\n", server.URL, response.Status)
		return err
	}

	return nil

}

// NewAPIRequest ...
func (server *Server) NewAPIRequest(method, APICall string, jsonString []byte) (*APIResult, error) {

	fullURL := server.URL + APICall

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: server.AllowUnverifiedSSL,
		},
	}

	server.httpClient = &http.Client{
		Transport: t,
	}

	request, requestErr := http.NewRequest(method, fullURL, bytes.NewBuffer(jsonString))
	if requestErr != nil {
		return nil, requestErr
	}

	request.SetBasicAuth(server.username, server.password)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	//if Debug {
	//dump, _ := httputil.DumpRequestOut(request, true)
	//fmt.Printf("HTTP Request\n%s\n", dump)
	//}

	response, doErr := server.httpClient.Do(request)
	defer response.Body.Close()

	if doErr != nil {
		return nil, doErr
	}

	var results APIResult
	if decodeErr := json.NewDecoder(response.Body).Decode(&results); decodeErr != nil {
		return nil, decodeErr
	}

	if results.Code == 0 {
		//fmt.Println("Setting Result Code")
		results.Code = response.StatusCode
	}
	if results.Status == "" {
		//fmt.Println("Setting Result Status")
		results.Status = response.Status
	}

	return &results, nil

}