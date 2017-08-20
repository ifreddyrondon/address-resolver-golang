// Entry point to address resolver server app
package main

import (
	"github.com/ifreddyrondon/address-resolver-golang/app"
	"os"
)

func main() {
	app.Initialize(os.Getenv("DATABASE_URL"))
	address := ":8080"
	app.Run(address)
}
