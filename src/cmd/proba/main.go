package main

import "proba"

func main() {
	app := proba.NewApplication()
	defer app.Close()
	app.Start()
}
