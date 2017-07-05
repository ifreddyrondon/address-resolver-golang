package gmap_test

import (
	"testing"
	"github.com/ifreddyrondon/address-resolver/gmap"
)

func TestGetLatLngFromSimpleAddress(t *testing.T) {
	coordinate, err := gmap.GetLatLngFromAddress("Apoquindo")
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

func TestGetLatLngFromComplexAddress(t *testing.T) {
	coordinate, err := gmap.GetLatLngFromAddress("Apoquindo")
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
