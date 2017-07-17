package addresses

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ifreddyrondon/address-resolver/database"
	"github.com/ifreddyrondon/address-resolver/gmap"
	"net/http"
	"path"
	"strconv"
)

func Router(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		Get(w, r)
	case "POST":
		CreateAddress(w, r)
	case "PUT":
		UpdateProduct(w, r)
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
	count := 10
	start := 0
	list, err := GetAddresses(database.GetDB(), start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, list)
}

func GetAddress(w http.ResponseWriter, r *http.Request) {
	var addressID string
	fmt.Sscanf(r.URL.Path, "/address/%s", &addressID)
	id, err := strconv.Atoi(addressID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid address ID")
		return
	}

	address := Address{ID: id}
	if err := address.getAddress(database.GetDB()); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Address not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, address)
}

func CreateAddress(w http.ResponseWriter, r *http.Request) {
	var address Address
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&address); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	coordinate, err := gmap.GetLatLngFromAddress(address.Address)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	address.Lat = coordinate.Lat
	address.Lng = coordinate.Lng

	if err := address.createAddress(database.GetDB()); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, address)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	baseUrl := path.Base(r.URL.Path)
	if baseUrl == "address" {
		fmt.Fprintf(w, "Method %s not supported\n", r.Method)
	}

	var addressID string
	fmt.Sscanf(r.URL.Path, "/address/%s", &addressID)
	id, err := strconv.Atoi(addressID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid address ID")
		return
	}

	var address Address
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&address); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	address.ID = id

	if err := address.updateAddress(database.GetDB()); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, address)
}


func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
