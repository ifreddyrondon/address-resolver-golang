// Entry point to address resolver server app
package main

import (
	"github.com/ifreddyrondon/address-resolver/app"
	"os"
)

func main() {
	application := app.App{}
	application.Initialize(os.Getenv("DATABASE_URL"))
	application.Run(":8080")
}
