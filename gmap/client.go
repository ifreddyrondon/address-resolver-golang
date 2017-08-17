package gmap

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

const (
	baseGoogleMapURI = "https://maps.googleapis.com"
	geocodingURI     = "/maps/api/geocode/json"
)

type AddressClient interface {
	GetGeocoding(address string) (*http.Response, error)
}

var (
	onceClient     sync.Once
	clientInstance *client
)

func GetClient() *client {
	onceClient.Do(func() {
		clientInstance = &client{}
	})
	return clientInstance
}

type client struct{}

func (c client) GetGeocoding(address string) (*http.Response, error) {
	uri := fmt.Sprintf("%s%s?address=%s", baseGoogleMapURI, geocodingURI, url.QueryEscape(address))
	resp, err := http.Get(uri)
	return resp, err
}
