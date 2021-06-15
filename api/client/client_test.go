package api

import (
	"log"
	"testing"
)

func TestGetApp(t *testing.T) {
	app, err := GetApp()

	if err != nil {
		t.Error(err)
	}

	app1, err := GetApp()

	if err != nil {
		t.Error(err)
	}

	log.Println("app == nil: ", app == nil)
	log.Println("app == app1: ", app == app1)
}

func TestGetAdmin(t *testing.T) {
	admin, err := GetAdmin()

	if err != nil {
		t.Error(err)
	}

	admin1, err := GetAdmin()

	if err != nil {
		t.Error(err)
	}

	log.Println("admin == nil: ", admin == nil)
	log.Println("admin == admin1: ", admin == admin1)
}
