package app

import (
	"github.com/ifreddyrondon/address-resolver/addresses"
	"github.com/ifreddyrondon/address-resolver/database"
	"github.com/ifreddyrondon/address-resolver/gognar"
)

var app gognar.GogApp

func Initialize(dbConnectionUrl string) *gognar.GogApp {
	app.Initialize()
	database.CreateConnection(dbConnectionUrl)
	initializeRoutes()
	return &app
}

func Run(address string) {
	app.Run(address)
}

func initializeRoutes() {
	app.Router.HandleFunc("/address/", addresses.Router)
}
