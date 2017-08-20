package addresses_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"encoding/json"

	"fmt"

	"bytes"

	"github.com/ifreddyrondon/address-resolver-golang/addresses"
	"github.com/ifreddyrondon/address-resolver-golang/app"
	"github.com/ifreddyrondon/address-resolver-golang/database"
	"github.com/ifreddyrondon/address-resolver-golang/gognar"
)

var application *gognar.GogApp

func TestMain(m *testing.M) {
	application = app.Initialize(os.Getenv("DATABASE_URL_TEST"))
	code := m.Run()
	database.ClearTable()
	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	application.Router.ServeHTTP(rr, req)

	return rr
}

func addAddress(count int) {
	for i := 0; i < count; i++ {
		database.GetDB().Exec(addresses.CreateAddressQuery, "Address "+strconv.Itoa(i), (i+1.0)*10, (i-1.0)*10)
	}
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func checkErrorResponse(t *testing.T, expected, actual map[string]interface{}) {
	if actual["error"] != expected["error"] {
		t.Errorf("Expected the Error %v. Got %v", expected["error"], actual["error"])
	}

	if actual["message"] != expected["message"] {
		t.Errorf("Expected the Message '%v'. Got %v", expected["message"], actual["message"])
	}

	if actual["status"] != expected["status"] {
		t.Errorf("Expected the Status '%v'. Got %v", expected["status"], actual["status"])
	}
}

func TestGetAddressList(t *testing.T) {
	tt := []struct {
		name           string
		addressesCount int
	}{
		{"empty table", 0},
		{"list of two addresses", 2},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.ClearTable()
			addAddress(tc.addressesCount)
			req, _ := http.NewRequest("GET", "/address/", nil)
			response := executeRequest(req)

			checkResponseCode(t, http.StatusOK, response.Code)

			var body []interface{}
			json.Unmarshal(response.Body.Bytes(), &body)

			if len(body) != tc.addressesCount {
				t.Errorf("Expected an array with len %d. Got %d", tc.addressesCount, len(body))
			}
		})
	}
}

func TestGetAddress(t *testing.T) {
	tt := []struct {
		name     string
		value    string
		status   int
		response map[string]interface{}
	}{
		{
			name:   "get address",
			value:  "1",
			status: http.StatusOK,
			response: map[string]interface{}{
				"id":      1.0,
				"address": "Address 0",
			},
		},
		{
			name:   "missing address",
			value:  "2",
			status: http.StatusNotFound,
			response: map[string]interface{}{
				"status":  404.0,
				"error":   "Not Found",
				"message": "Address not found",
			},
		},
		{
			name:   "bad request",
			value:  "x",
			status: http.StatusBadRequest,
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "Invalid address ID",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.ClearTable()
			addAddress(1)

			req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s", tc.value), nil)
			response := executeRequest(req)

			checkResponseCode(t, tc.status, response.Code)

			var m map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &m)

			if tc.status != http.StatusOK {
				checkErrorResponse(t, tc.response, m)
			}

			if m["id"] != tc.response["id"] {
				t.Errorf("Expected the id %v. Got %v", tc.response["id"], m["id"])
			}

			if m["address"] != tc.response["address"] {
				t.Errorf("Expected the address '%v'. Got %v", tc.response["address"], m["address"])
			}
		})
	}
}

func TestCreateAddress(t *testing.T) {
	tt := []struct {
		name     string
		payload  []byte
		status   int
		response map[string]interface{}
	}{
		{
			name:    "create address",
			payload: []byte(`{"address": "ejido, manzano bajo"}`),
			status:  http.StatusCreated,
			response: map[string]interface{}{
				"id":      1.0,
				"address": "ejido, manzano bajo",
			},
		},
		{
			name:    "bad request, missing address field",
			payload: []byte(`{"a": "ejido, manzano bajo"}`),
			status:  http.StatusBadRequest,
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "Invalid request payload",
			},
		},
		{
			name:    "bad request, missing body",
			payload: nil,
			status:  http.StatusBadRequest,
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "Invalid request payload",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.ClearTable()

			req, _ := http.NewRequest("POST", "/address/", bytes.NewBuffer(tc.payload))
			response := executeRequest(req)

			checkResponseCode(t, tc.status, response.Code)

			var m map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &m)

			if tc.status != http.StatusOK {
				checkErrorResponse(t, tc.response, m)
			}

			if m["address"] != tc.response["address"] {
				t.Errorf("Expected address to be '%v'. Got '%v'", tc.response["address"], m["address"])
			}

			if m["id"] != tc.response["id"] {
				t.Errorf("Expected address ID to be '%v'. Got '%v'", tc.response["id"], m["id"])
			}
		})
	}
}

func TestUpdateAddress(t *testing.T) {
	tt := []struct {
		name      string
		updateUrl string
		payload   []byte
		status    int
		response  map[string]interface{}
	}{
		{
			name:      "update address",
			updateUrl: "/address/1",
			payload:   []byte(`{"address": "ejido, manzano bajo"}`),
			status:    http.StatusOK,
			response: map[string]interface{}{
				"id":      1.0,
				"address": "ejido, manzano bajo",
			},
		},
		{
			name:      "bad request",
			updateUrl: "/address/x",
			status:    http.StatusBadRequest,
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "Invalid address ID",
			},
		},
		{
			name:      "not allowed",
			updateUrl: "/address/",
			status:    http.StatusMethodNotAllowed,
			response: map[string]interface{}{
				"status":  405.0,
				"error":   "Method Not Allowed",
				"message": "Method PUT not supported",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.ClearTable()
			addAddress(1)

			req, _ := http.NewRequest("GET", "/address/1", nil)
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)

			var originalAddress map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &originalAddress)

			req, _ = http.NewRequest("PUT", tc.updateUrl, bytes.NewBuffer(tc.payload))
			response = executeRequest(req)

			checkResponseCode(t, tc.status, response.Code)

			var m map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &m)

			if tc.status != http.StatusOK {
				checkErrorResponse(t, tc.response, m)
			}

			if m["id"] != tc.response["id"] {
				t.Errorf("Expected the id to remain the same (%v). Got %v", tc.response["id"], m["id"])
			}

			if m["address"] == originalAddress["address"] {
				t.Errorf("Expected the address to change from '%v' to '%v'. Got '%v'", originalAddress["address"], tc.response["address"], m["address"])
			}
		})
	}
}

func TestDeleteAddress(t *testing.T) {
	tt := []struct {
		name      string
		deleteUrl string
		status    int
		response  map[string]interface{}
	}{
		{
			name:      "update address",
			deleteUrl: "/address/1",
			status:    http.StatusNoContent,
			response: map[string]interface{}{
				"id":      1.0,
				"address": "ejido, manzano bajo",
			},
		},
		{
			name:      "bad request",
			deleteUrl: "/address/x",
			status:    http.StatusBadRequest,
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "Invalid address ID",
			},
		},
		{
			name:      "not allowed",
			deleteUrl: "/address/",
			status:    http.StatusMethodNotAllowed,
			response: map[string]interface{}{
				"status":  405.0,
				"error":   "Method Not Allowed",
				"message": "Method DELETE not supported",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.ClearTable()
			addAddress(1)

			req, _ := http.NewRequest("GET", "/address/1", nil)
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)

			var originalAddress map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &originalAddress)

			req, _ = http.NewRequest("DELETE", tc.deleteUrl, nil)
			response = executeRequest(req)

			checkResponseCode(t, tc.status, response.Code)

			var m map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &m)

			if tc.status != http.StatusOK {
				checkErrorResponse(t, tc.response, m)
				return
			}

			req, _ = http.NewRequest("GET", "/address/1", nil)
			response = executeRequest(req)
			checkResponseCode(t, http.StatusNotFound, response.Code)
		})
	}
}

func TestPathMethodNotAllowed(t *testing.T) {
	req, _ := http.NewRequest("PATH", "/address/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusMethodNotAllowed, response.Code)

	var expected = map[string]interface{}{
		"status":  405.0,
		"error":   "Method Not Allowed",
		"message": "Method PATH not supported",
	}

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	checkErrorResponse(t, expected, m)
}
