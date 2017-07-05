package addresses

import (
	"net/http"
	"fmt"
	"path"
	"encoding/json"
	"github.com/ifreddyrondon/address-resolver/gmap"
)

func Router(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		Get(w, r)
	case "POST":
		CreateAddress(w, r)
	default:
		fmt.Fprintf(w, "Method %s not supported\n", r.Method)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	baseUrl := path.Base(r.URL.Path)
	if baseUrl == "address" {
		GetAddressList(w, r)
	} else {
		GetAddress(w, r)
	}
}

func GetAddressList(w http.ResponseWriter, r *http.Request) {
	list := Addresses{
		{Address: "Apoquindo 4800", Lat: 1.1, Lng: 2.4},
		{Address: "Jorge Matte 1481", Lat: 4.3, Lng: 7.5},
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(list)
	if err != nil {
		panic(err)
	}
}

func GetAddress(w http.ResponseWriter, r *http.Request) {
	// addressID := path.Base(r.URL.Path)
	var addressID string
	fmt.Sscanf(r.URL.Path, "/address/%s", &addressID)
	fmt.Fprintf(w, "GET address: %s\n", addressID)
}

func CreateAddress(w http.ResponseWriter, r *http.Request) {
	var address Address

	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	coodinate, err := gmap.GetLatLngFromAddress(address.Address)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	address.Lat = coodinate.Lat
	address.Lng = coodinate.Lng

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(address); err != nil {
		panic(err)
	}
}
