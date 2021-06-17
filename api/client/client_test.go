package client

import (
	"log"
	"testing"
)

var (
	clientConfigPath1 string = "./test_config/client-org1.yaml"
	clientConfigPath2 string = "./test_config/client-org2.yaml"
)

func TestGetApp(t *testing.T) {
	app2, err := GetApp("channel2", clientConfigPath2)

	if err != nil {
		t.Error(err)
	}

	_app2, err := GetApp("channel2", clientConfigPath2)

	if err != nil {
		t.Error(err)
	}

	app1, err := GetApp("channel1", clientConfigPath1)

	if err != nil {
		t.Error(err)
	}

	if app2 == nil {
		t.Error("app2 == nil")
	} else {
		log.Println("app2 == nil: ", app2 == nil)
	}

	if app2 != _app2 {
		t.Error("app2 != _app2")
	} else {
		log.Println("app2 == _app2: ", app2 == _app2)
	}

	if app2 == app1 {
		t.Error("app2 == app1")
	} else {
		log.Println("app2 == app1: ", app2 == app1)
	}
}

func TestGetAdmin(t *testing.T) {
	admin2, err := GetAdmin(clientConfigPath2)

	if err != nil {
		t.Error(err)
	}

	_admin2, err := GetAdmin(clientConfigPath2)

	if err != nil {
		t.Error(err)
	}

	admin1, err := GetAdmin(clientConfigPath1)

	if err != nil {
		t.Error(err)
	}

	log.Println("admin2 == nil: ", admin2 == nil)
	log.Println("admin2 == _admin2: ", admin2 == _admin2)
	log.Println("admin2 == admin1: ", admin2 == admin1)
}
