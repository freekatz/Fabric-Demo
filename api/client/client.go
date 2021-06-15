package api

import (
	"fmt"
	"sync"

	"github.com/1uvu/fabric-sdk-client/client"
)

var (
	app   *client.AppClient
	admin *client.AdminClient

	appOnce   sync.Once
	adminOnce sync.Once

	clientConfigPath string = "client_config.yaml"
)

func GetApp() (*client.AppClient, error) {
	var err error
	appOnce.Do(func() {
		if app == nil {
			app, err = newApp()
		}
	})

	return app, err
}

func GetAdmin() (*client.AdminClient, error) {
	var err error
	adminOnce.Do(func() {
		if admin == nil {
			admin, err = newAdmin()
		}
	})

	return admin, err
}

func newApp() (*client.AppClient, error) {

	conf, err := newClientConfig(clientConfigPath)

	if err != nil {
		return nil, fmt.Errorf("failed to get app client: %s", err)
	}

	params := &conf.App.Params

	envPairs := conf.App.Envs

	app, err = client.GetAppClient("channel2", params, envPairs...)
	if err != nil {
		return nil, fmt.Errorf("failed to get app client: %s", err)
	}

	return app, nil
}

func newAdmin() (*client.AdminClient, error) {

	conf, err := newClientConfig(clientConfigPath)

	if err != nil {
		return nil, fmt.Errorf("failed to get app client: %s", err)
	}

	params := &conf.Admin.Params

	envPairs := conf.Admin.Envs

	admin, err = client.GetAdminClient(params, envPairs...)
	if err != nil {
		return nil, fmt.Errorf("failed to get app client: %s", err)
	}

	return admin, nil
}
