package gognar

import (
	"log"
	"net/http"
)

type GogApp struct {
	Router *http.ServeMux
}

func (app *GogApp) Initialize() {
	app.Router = http.NewServeMux()
}

func (app *GogApp) Run(address string) {
	log.Printf("Running on %s", address)
	log.Fatal(http.ListenAndServe(address, app.Router))
}
