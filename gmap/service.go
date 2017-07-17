package gmap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	baseGoogleMapURI = "https://maps.googleapis.com"
	geocodingURI     = "/maps/api/geocode/json"
)

func GetLatLngFromAddress(address string) (*Coordinate, error) {
	uri := fmt.Sprintf("%s%s?address=%s", baseGoogleMapURI, geocodingURI, url.QueryEscape(address))
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := new(geocodingResult)
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Results[0].Geometry.Coordinate, nil
}
