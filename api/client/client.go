package api

import (
	"fmt"
	"sync"

	"github.com/1uvu/fabric-sdk-client/client"
)

var (
	apps   map[string]*client.AppClient   = make(map[string]*client.AppClient)
	admins map[string]*client.AdminClient = make(map[string]*client.AdminClient)

	appOnce   sync.Once
	adminOnce sync.Once
)

// todo 修改为 abstract factory
func GetApp(clientConfigPath string) (*client.AppClient, error) {
	var err error
	appOnce.Do(func() {
		if _, ok := apps[clientConfigPath]; !ok {
			apps[clientConfigPath], err = newApp(clientConfigPath)
		}
	})

	return apps[clientConfigPath], err
}

func GetAdmin(clientConfigPath string) (*client.AdminClient, error) {
	var err error
	adminOnce.Do(func() {
		if _, ok := admins[clientConfigPath]; !ok {
			admins[clientConfigPath], err = newAdmin(clientConfigPath)
		}
	})

	return admins[clientConfigPath], err
}

func newApp(clientConfigPath string) (*client.AppClient, error) {

	conf, err := newClientConfig(clientConfigPath)

	if err != nil {
		return nil, fmt.Errorf("failed to get app client: %s", err)
	}

	params := &conf.App.Params

	envPairs := conf.App.Envs

	app, err := client.GetAppClient("channel2", params, envPairs...)
	if err != nil {
		return nil, fmt.Errorf("failed to get app client: %s", err)
	}

	return app, nil
}

func newAdmin(clientConfigPath string) (*client.AdminClient, error) {

	conf, err := newClientConfig(clientConfigPath)

	if err != nil {
		return nil, fmt.Errorf("failed to get app client: %s", err)
	}

	params := &conf.Admin.Params

	envPairs := conf.Admin.Envs

	admin, err := client.GetAdminClient(params, envPairs...)
	if err != nil {
		return nil, fmt.Errorf("failed to get app client: %s", err)
	}

	return admin, nil
}
