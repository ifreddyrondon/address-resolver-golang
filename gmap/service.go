package gmap

import (
	"encoding/json"
	"errors"
)

var serviceInstance = &AddressService{client: GetClient()}

func GetService() *AddressService {
	return serviceInstance
}

type AddressService struct {
	client AddressClient
}

func (a *AddressService) GetLatLngFromAddress(address string) (*Coordinate, error) {
	resp, err := a.client.GetGeocoding(address)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := new(geocodingResult)
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Results) < 1 {
		return nil, errors.New("Not found")
	}

	return &result.Results[0].Geometry.Coordinate, nil
}
