package addresses

import (
	"database/sql"
	"fmt"
	"net/http"
	"path"
	"strconv"

	"errors"

	"github.com/ifreddyrondon/address-resolver/database"
	"github.com/ifreddyrondon/address-resolver/gmap"
	"github.com/ifreddyrondon/address-resolver/gognar"
)

func Router(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		Get(w, r)
	case "POST":
		CreateAddress(w, r)
	case "PUT":
		UpdateProduct(w, r)
	case "DELETE":
		DeleteProduct(w, r)
	default:
		gognar.MethodNotAllowed(w, errors.New(fmt.Sprintf("Method %s not supported", r.Method)))
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
		gognar.InternalServerError(w, err)
		return
	}

	gognar.Send(w, list)
}

func GetAddress(w http.ResponseWriter, r *http.Request) {
	var addressID string
	fmt.Sscanf(r.URL.Path, "/address/%s", &addressID)
	id, err := strconv.Atoi(addressID)
	if err != nil {
		gognar.BadRequest(w, errors.New("Invalid address ID"))
		return
	}

	address := Address{ID: id}
	if err := address.getAddress(database.GetDB()); err != nil {
		switch err {
		case sql.ErrNoRows:
			gognar.NotFound(w, errors.New("Address not found"))
		default:
			gognar.InternalServerError(w, err)
		}
		return
	}

	gognar.Send(w, address)
}

func CreateAddress(w http.ResponseWriter, r *http.Request) {
	var address Address
	if err := gognar.ReadJSON(r.Body, &address); err != nil {
		gognar.BadRequest(w, errors.New("Invalid request payload"))
		return
	}

	if address.Address == "" {
		gognar.BadRequest(w, errors.New("Invalid request payload"))
		return
	}

	coordinate, err := gmap.GetLatLngFromAddress(address.Address)
	if err != nil {
		gognar.InternalServerError(w, err)
		return
	}

	address.Lat = coordinate.Lat
	address.Lng = coordinate.Lng

	if err := address.createAddress(database.GetDB()); err != nil {
		gognar.InternalServerError(w, err)
		return
	}

	gognar.Created(w, address)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	baseUrl := path.Base(r.URL.Path)
	if baseUrl == "address" {
		gognar.MethodNotAllowed(w, errors.New(fmt.Sprintf("Method %s not supported", r.Method)))
		return
	}

	var addressID string
	fmt.Sscanf(r.URL.Path, "/address/%s", &addressID)
	id, err := strconv.Atoi(addressID)
	if err != nil {
		gognar.BadRequest(w, errors.New("Invalid address ID"))
		return
	}

	var address Address
	if err := gognar.ReadJSON(r.Body, &address); err != nil {
		gognar.BadRequest(w, errors.New("Invalid resquest payload"))
		return
	}
	address.ID = id

	if err := address.updateAddress(database.GetDB()); err != nil {
		gognar.InternalServerError(w, err)
		return
	}

	gognar.Send(w, address)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	baseUrl := path.Base(r.URL.Path)
	if baseUrl == "address" {
		gognar.MethodNotAllowed(w, errors.New(fmt.Sprintf("Method %s not supported", r.Method)))
		return
	}

	var addressID string
	fmt.Sscanf(r.URL.Path, "/address/%s", &addressID)
	id, err := strconv.Atoi(addressID)
	if err != nil {
		gognar.BadRequest(w, errors.New("Invalid address ID"))
		return
	}

	address := Address{ID: id}
	if err := address.deleteAddress(database.GetDB()); err != nil {
		gognar.InternalServerError(w, err)
		return
	}

	gognar.NoContent(w)
}
