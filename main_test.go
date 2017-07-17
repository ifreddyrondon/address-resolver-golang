package main

import (
	"bytes"
	"encoding/json"
	"github.com/ifreddyrondon/address-resolver/app"
	"github.com/ifreddyrondon/address-resolver/database"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var application app.App

func TestMain(m *testing.M) {
	application.Initialize(os.Getenv("DATABASE_URL_TEST"))
	code := m.Run()
	database.ClearTable()
	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	application.Router.ServeHTTP(rr, req)

	return rr
}
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func addAddress(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		database.GetDB().Exec("INSERT INTO addresses(address, lat, lng) VALUES($1, $2, $3)", "Address "+strconv.Itoa(i), (i+1.0)*10, (i-1.0)*10)
	}
}

func TestEmptyTable(t *testing.T) {
	database.ClearTable()
	req, _ := http.NewRequest("GET", "/address/", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentAddress(t *testing.T) {
	database.ClearTable()

	req, _ := http.NewRequest("GET", "/address/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Address not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Address not found'. Got '%s'", m["error"])
	}
}

func TestCreateAddress(t *testing.T) {
	database.ClearTable()

	payload := []byte(`{"address": "ejido, manzano bajo"}`)

	req, _ := http.NewRequest("POST", "/address/", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["address"] != "ejido, manzano bajo" {
		t.Errorf("Expected address to be 'ejido, manzano bajo'. Got '%v'", m["address"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected address ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetAddress(t *testing.T) {
	database.ClearTable()
	addAddress(1)

	req, _ := http.NewRequest("GET", "/address/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateAddress(t *testing.T) {
	database.ClearTable()
	addAddress(1)

	req, _ := http.NewRequest("GET", "/address/1", nil)
	response := executeRequest(req)
	var originalAddress map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalAddress)

	payload := []byte(`{"address":"test address - updated name"}`)

	req, _ = http.NewRequest("PUT", "/address/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalAddress["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalAddress["id"], m["id"])
	}

	if m["address"] == originalAddress["address"] {
		t.Errorf("Expected the address to change from '%v' to '%v'. Got '%v'", originalAddress["address"], m["address"], m["address"])
	}
}

