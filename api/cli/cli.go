package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
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
type AppParams struct {
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
	AppParams
	*gateway.Network
}

var (
	appParams *AppParams
)

func SetAppEnv(sw string) {
	os.Setenv("DISCOVERY_AS_LOCALHOST", sw)
}

func SetAppParams(p *AppParams) {
	appParams = p
}

//
// 存储认证材料
//
func populateWallet(wallet *gateway.Wallet) error {
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(appParams.CertPath))
	if err != nil {
		return err
	}

	// there's a single file in this dir containing the private key
	keyDir := filepath.Join(appParams.CredPath, "keystore")
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

	identity := gateway.NewX509Identity(appParams.OrgMSP, string(cert), string(key))

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
	if appParams.ConfigPath == "" {
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
		gateway.WithConfig(config.FromFile(filepath.Clean(appParams.ConfigPath))),
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

	return &AppClient{AppParams: *appParams, Network: network}, nil
}

// 用于连接 Fabric 网络的参数
type AdminParams struct {
	CredPath   string
	CertPath   string
	ConfigPath string
	OrgName    string
	OrgAdmin   string
	OrgUser    string
	OrgMSP     string
	OrgHost    string
}

type appClient struct {
	ChannelID string

	CC *channel.Client
	EC *event.Client
	LC *ledger.Client
}

type AdminClient struct {
	AdminParams

	// sdk clients
	SDK *fabsdk.FabricSDK
	RC  *resmgmt.Client
	MC  *msp.Client

	acs map[string]*appClient

	// for create channel
	// ChannelConfig string
	// OrdererID string
}

var (
	adminParams *AdminParams
)

func SetAdminEnv(sw string) {
	os.Setenv("DISCOVERY_AS_LOCALHOST", sw)
}

func SetAdminParams(p *AdminParams) {
	adminParams = p
}

func GetAdminClient() (*AdminClient, error) {
	if adminParams.ConfigPath == "" {
		return nil, fmt.Errorf("Please init the params by call SerParams() function.")
	}

	sdk, err := fabsdk.New(config.FromFile(adminParams.ConfigPath))
	if err != nil {
		log.Panicf("failed to create fabric sdk: %s", err)
	}

	rcp := sdk.Context(fabsdk.WithUser(adminParams.OrgAdmin), fabsdk.WithOrg(adminParams.OrgName))
	rc, err := resmgmt.New(rcp)
	if err != nil {
		log.Panicf("failed to create resource client: %s", err)
	}

	mc, err := msp.New(sdk.Context())
	if err != nil {
		log.Panicf("failed to create msp client: %s", err)
	}

	admin := new(AdminClient)

	admin.AdminParams = *adminParams
	admin.SDK = sdk
	admin.RC = rc
	admin.MC = mc
	admin.acs = make(map[string]*appClient)

	return admin, nil
}

func (admin *AdminClient) GetAppClient(channelID string) (*appClient, error) {

	if app, ok := admin.acs[channelID]; ok {
		log.Printf("app client of %s has existed, return directly.", channelID)
		return app, nil
	}

	log.Printf("app client of %s do not existed, get it now.", channelID)
	// 这里是 admin 端的 channel client，因此为其指定 admin 用户
	// ccp := admin.SDK.ChannelContext(channelID, fabsdk.WithUser(admin.OrgUser))
	ccp := admin.SDK.ChannelContext(channelID, fabsdk.WithUser(admin.OrgAdmin), fabsdk.WithUser(admin.OrgUser))
	cc, err := channel.New(ccp)
	if err != nil {
		return nil, fmt.Errorf("failed to create channel client: %s", err)
	}

	ec, err := event.New(ccp, event.WithSeekType(seek.Newest))
	if err != nil {
		return nil, fmt.Errorf("failed to create event client: %s", err)
	}

	lc, err := ledger.New(ccp)
	if err != nil {
		return nil, fmt.Errorf("failed to create ledger client: %s", err)
	}

	app := new(appClient)
	app.ChannelID = channelID
	app.CC = cc
	app.EC = ec
	app.LC = lc

	admin.acs[channelID] = app

	return app, nil
}
