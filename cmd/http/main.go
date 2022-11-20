package main

import (
	"email-svc/src/infrastructure/http"
)

func main() {
	app := http.NewServer()
	app.Run()
}
