// Entry point to address resolver server app
package main

import (
	"net/http"
	"log"
	"github.com/ifreddyrondon/address-resolver/addresses"
)

func main() {
	http.HandleFunc("/address/", addresses.Router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
