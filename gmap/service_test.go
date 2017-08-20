package gmap_test

import (
	"errors"
	"fmt"
	"testing"

	"net/http"

	"io/ioutil"
	"net/http/httptest"
	"path/filepath"

	"github.com/ifreddyrondon/address-resolver-golang/gmap"
	"github.com/ifreddyrondon/address-resolver-golang/gmap/gmapfakes"
)

func getClientResponse(fixtureFile string) *http.Response {
	response := httptest.NewRecorder()
	fileName := filepath.Join("testdata", fixtureFile)
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	response.Write(content)
	return response.Result()
}

func TestGetLatLngAddress(t *testing.T) {
	client := new(gmapfakes.FakeAddressClient)
	gmap.SetClientInstance(client)
	client.GetGeocodingReturns(getClientResponse("apoquindo_fixture.json"), nil)

	gmapService := gmap.GetService()
	coordinate, err := gmapService.GetLatLngFromAddress("Apoquindo")
	if err != nil {
		t.Fatal(err)
	}

	if coordinate.Lat != -33.410267 {
		t.Error(fmt.Sprintf("Lat should not be zero. Got '%v", coordinate.Lat))
	}

	if coordinate.Lng != -70.5723 {
		t.Error(fmt.Sprintf("Lng should not be zero. Got '%v", coordinate.Lng))
	}

	if client.GetGeocodingCallCount() != 1 {
		t.Error(fmt.Sprintf("Client should be called 1. Got '%v", client.GetGeocodingCallCount()))
	}
}

func TestGetErrorFromAddressWhenNotResult(t *testing.T) {
	client := new(gmapfakes.FakeAddressClient)
	gmap.SetClientInstance(client)
	client.GetGeocodingReturns(getClientResponse("zero_results_fixture.json"), nil)

	gmapService := gmap.GetService()
	coordinate, err := gmapService.GetLatLngFromAddress("Apoquindo")
	if coordinate != nil {
		t.Error(fmt.Sprintf("coordinate should be nil. Got '%v", coordinate))
	}

	if err.Error() != "Not found" {
		t.Error(fmt.Sprintf("Error message should be 'Not found'. Got '%v", err.Error()))
	}

	if client.GetGeocodingCallCount() != 1 {
		t.Error(fmt.Sprintf("Client should be called 1. Got '%v", client.GetGeocodingCallCount()))
	}
}

func TestGetErrorFromAddressWhenClientFails(t *testing.T) {
	client := new(gmapfakes.FakeAddressClient)
	gmap.SetClientInstance(client)
	client.GetGeocodingReturns(nil, errors.New("failed!"))

	gmapService := gmap.GetService()
	coordinate, err := gmapService.GetLatLngFromAddress("Apoquindo")
	if coordinate != nil {
		t.Error(fmt.Sprintf("coordinate should be nil. Got '%v", coordinate))
	}

	if err.Error() != "failed!" {
		t.Error(fmt.Sprintf("Error message should be 'failed!'. Got '%v", err.Error()))
	}

	if client.GetGeocodingCallCount() != 1 {
		t.Error(fmt.Sprintf("Client should be called 1. Got '%v", client.GetGeocodingCallCount()))
	}
}
