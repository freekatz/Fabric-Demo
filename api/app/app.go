package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

/*
channel client
用于编写 APP 的各种常用方法
使用 cc config，也可通过 rc 获取

基于 pkg/gateway
*/

//
// 封装定义数据类型
//

// 用于连接 Fabric 网络的参数
type ClientParams struct {
	CredPath   string
	CertPath   string
	ConfigPath string
	OrgName    string
	OrgAdmin   string
	OrgUser    string
	OrgMSP     string
	OrgHost    string
}

type AppClient struct {
	ClientParams
	*gateway.Network
}

var (
	params *ClientParams
)

func SetEnv(sw string) {
	os.Setenv("DISCOVERY_AS_LOCALHOST", sw)
}

func SetParams(p *ClientParams) {
	params = p
}

//
// 存储认证材料
//
func populateWallet(wallet *gateway.Wallet) error {
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(params.CertPath))
	if err != nil {
		return err
	}

	// there's a single file in this dir containing the private key
	keyDir := filepath.Join(params.CredPath, "keystore")
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity(params.OrgMSP, string(cert), string(key))

	err = wallet.Put("appUser", identity)
	if err != nil {
		return err
	}
	return nil
}

//
// 获取网络实例
//
func GetAppClient(channelID string) (*AppClient, error) {
	if params.ConfigPath == "" {
		return nil, fmt.Errorf("Please init the params by call SerParams() function.")
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			fmt.Printf("Failed to populate wallet contents: %s\n", err)
			os.Exit(1)
		}
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(params.ConfigPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
	defer gw.Close()

	network, err := gw.GetNetwork(channelID)
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	return &AppClient{ClientParams: *params, Network: network}, nil
}
