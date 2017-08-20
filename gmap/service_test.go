package gmap_test

import (
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

func TestGetLatLngFromSimpleAddress(t *testing.T) {
	client := new(gmapfakes.FakeAddressClient)
	client.GetGeocodingReturns(getClientResponse("apoquindo_fixture.json"), nil)

	gmapService := gmap.GetService()
	coordinate, err := gmapService.GetLatLngFromAddress("Apoquindo")
	if err != nil {
		t.Fatal(err)
	}

	if coordinate.Lat == 0 {
		t.Error("Lat should not be zero")
	}

	if coordinate.Lng == 0 {
		t.Error("Lng should not be zero")
	}
}
