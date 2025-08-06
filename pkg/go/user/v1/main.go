package main

import (
	"github.com/Rabiann/weather-mailer/pkg/go/user/cmd"
)

func main() {
	var app cmd.App

	if err := app.Run(); err != nil {
		panic(err)
	}
}
