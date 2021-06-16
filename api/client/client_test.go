package client

import (
	"log"
	"testing"
)

var (
	clientConfigPath1 string = "./config/client-org1.yaml"
	clientConfigPath2 string = "./config/client-org2.yaml"
	clientConfigPath3 string = "./config/client-org3.yaml"
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

	log.Println("app2 == nil: ", app2 == nil)
	log.Println("app2 == _app2: ", app2 == _app2)
	log.Println("app2 == app1: ", app2 == app1)
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
