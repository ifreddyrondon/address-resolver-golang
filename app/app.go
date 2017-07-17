package app

import (
	"github.com/ifreddyrondon/address-resolver/addresses"
	"github.com/ifreddyrondon/address-resolver/database"
	"log"
	"net/http"
)

type App struct {
	Router *http.ServeMux
}

func (app *App) Initialize(dbConnectionUrl string) {
	database.CreateConnection(dbConnectionUrl)
	app.Router = http.NewServeMux()
	app.initializeRoutes()
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/address/", addresses.Router)
}
