package test

import (
	"database/sql"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/patrickmn/go-cache"
	"net/http"
	"regexp"
	"testing"
	"time"
)

func startService() {
	//TODO fix me
	const dbType = "mysql"
	const regexString = "^BQY[0-9a-zA-Z]{29}$"
	// read serv config
	const tomlData = "config.toml"

	fmt.Println("Decoding Toml file")
	if _, err := toml.DecodeFile(tomlData, &Conf); err != nil {
		fmt.Println("Error: Config file  ", err)
	}
	// prepare regex
	Regex, _ = regexp.Compile(regexString)

	var err error
	fmt.Println("Connecting DB")
	DB, err = sql.Open(dbType, Conf.DBConf.DBUser+":"+Conf.DBConf.DBPass+"@tcp("+Conf.DBConf.DBAddr+":"+Conf.DBConf.DBPort+")/")
	if err != nil {
		fmt.Println("Error: DB  ", err)
	}

	Cache = cache.New(10*time.Second, 10*time.Second)
	fmt.Println("Starting network..")
	http.HandleFunc("/parseHeaders", parseHeaders)
	err = http.ListenAndServe(Conf.NetworkConf.BindAddr, nil)
	if err != nil {
		fmt.Println("Network error: ", err)
	}
	//defer DB.Close()
}

func sendHttpRequest(headerName string, HeaderValue string) (*http.Response, error) {
	client := http.Client{}
	//172.16.13.4
	request, err := http.NewRequest("GET", "http://localhost:8080/parseHeaders", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add(headerName, HeaderValue)
	response, err := client.Do(request)
	if err != nil {
		return response, err
	}
	return response, err
}
func TestServiceCorrectAnswer(t *testing.T) {
	//startService()
	response, err := sendHttpRequest("x-api-key", "BQY123456789a123456789b123456789")
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Error(" Expected Response Status 200 got: ", response.StatusCode)
	}
	if response.Header.Get("X-Accel-Redirect") == "" {
		t.Error(" Expected Header  X-Accel-Redirect = @default got: ", response.Header.Get("X-Accel-Redirect"))
	}
}

func TestDefaultRespNotInDB(t *testing.T) {
	response, err := sendHttpRequest("x-api-key", "BQY123456789a123456789b123456aaa")
	if err != nil {
		t.Fatal(err)
	}
	if response.Header.Get("X-Accel-Redirect") != "@default" {
		t.Error(" Expected Header  X-Accel-Redirect = @default got: ", response.Header.Get("X-Accel-Redirect"))
	}
}
func TestServiceWrongAnswer(t *testing.T) {
	// wrong header name
	response, err := sendHttpRequest("x", "BQY123456789a123456789b123456aaa")
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusBadRequest {
		t.Error(" Expected Response Status 400 got: ", response.StatusCode)
	}
	// wrong format of key
	response, err = sendHttpRequest("x-api-key", "23456789a123456789b123456aaa")
	if response.StatusCode != http.StatusBadRequest {
		t.Error(" Expected Response Status 400 got: ", response.StatusCode)
	}
}
