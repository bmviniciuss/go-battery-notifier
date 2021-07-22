package main

import (
	"github.com/bmviniciuss/go-battery-notifier/application"
	"github.com/bmviniciuss/go-battery-notifier/notification"
	"github.com/bmviniciuss/go-battery-notifier/reader"
	"log"
)

func main() {
	r := reader.NewUbuntuReader()
	n := notification.NewBeeepNotifier()
	a, err := application.NewApplication(r, n, 80, 20, 1, true)
	if err != nil {
		log.Fatal(err)
	}
	a.Run()
}
