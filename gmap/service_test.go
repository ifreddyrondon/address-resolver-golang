package gmap_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ifreddyrondon/address-resolver/gmap"
	"github.com/ifreddyrondon/address-resolver/gmap/gmapfakes"
)

var client = new(gmapfakes.FakeAddressClient)

func init() {
	gmap.SetClientInstance(client)
}

func TestGetLatLngFromSimpleAddress(t *testing.T) {
	*client = gmapfakes.FakeAddressClient{}
	client.GetGeocodingReturns(GetResponse("apoquindo.json"), nil)

	gmapService := gmap.GetService()
	coordinate, err := gmapService.GetLatLngFromAddress("Apoquindo")
	if err != nil {
		t.Fatal(err)
	}

	if coordinate.Lat != -33.410267 {
		t.Error("Lat should not be zero")
	}

	if coordinate.Lng != -70.5723 {
		t.Error("Lng should not be zero")
	}
}

func TestErrorGettingLatLngFromSimpleAddress(t *testing.T) {
	*client = gmapfakes.FakeAddressClient{}
	client.GetGeocodingReturns(nil, errors.New("failed!"))

	gmapService := gmap.GetService()
	coordinate, err := gmapService.GetLatLngFromAddress("Apoquindo")
	if err != nil {
		t.Fatal(err)
	}

	if coordinate.Lat != -33.410267 {
		t.Error("Lat should not be zero")
	}

	if coordinate.Lng != -70.5723 {
		t.Error("Lng should not be zero")
	}
}

func GetResponse(filename string) *http.Response {
	response := httptest.NewRecorder()
	content, err := ioutil.ReadFile("./responses/" + filename)
	if err != nil {
		panic(err)
	}
	response.Write(content)
	return response.Result()
}
