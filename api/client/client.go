package client

import (
	"fmt"
	"strings"
	"sync"

	"github.com/1uvu/fabric-sdk-client/client"
)

var (
	apps   sync.Map
	admins sync.Map
)

// todo 修改为 abstract factory
func GetApp(channelID, clientConfigPath string) (*client.AppClient, error) {
	key := strings.Join([]string{channelID, clientConfigPath}, "-")
	if _, ok := apps.Load(key); !ok {
		app, err := newApp(channelID, clientConfigPath)
		if err != nil {
			return nil, err
		}
		apps.Store(key, app)
	}

	app, ok := apps.Load(key)
	if !ok {
		return nil, fmt.Errorf("failed to get app client of %s", channelID)
	}

	return app.(*client.AppClient), nil
}

func GetAdmin(clientConfigPath string) (*client.AdminClient, error) {
	key := clientConfigPath
	if _, ok := admins.Load(key); !ok {
		admin, err := newAdmin(clientConfigPath)
		if err != nil {
			return nil, err
		}

		admins.Store(key, admin)
	}

	admin, ok := admins.Load(key)
	if !ok {
		return nil, fmt.Errorf("failed to get admin client")
	}

	return admin.(*client.AdminClient), nil
}

func newApp(channelID, clientConfigPath string) (*client.AppClient, error) {

	conf, err := newClientConfig(clientConfigPath)

	if err != nil {
		return nil, fmt.Errorf("failed to get app client: %s", err)
	}

	params := &conf.App.Params

	envPairs := conf.App.Envs

	app, err := client.GetAppClient(channelID, params, envPairs...)
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
