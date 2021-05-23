package client

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
type ClientParams struct {
	CredPath        string
	CertPath        string
	AppConfigPath   string
	AdminConfigPath string
	OrgName         string
	OrgAdmin        string
	OrgUser         string
	OrgMSP          string
	OrgHost         string
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
	if params.AppConfigPath == "" {
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
		gateway.WithConfig(config.FromFile(filepath.Clean(params.AppConfigPath))),
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

type AdminClient struct {
	ClientParams

	// sdk clients
	SDK *fabsdk.FabricSDK
	RC  *resmgmt.Client
	MC  *msp.Client

	ChannelID string

	CC *channel.Client
	EC *event.Client
	LC *ledger.Client

	// for create channel
	// ChannelConfig string
	// OrdererID string
}

func GetAdminClient() (*AdminClient, error) {
	if params.AdminConfigPath == "" {
		return nil, fmt.Errorf("Please init the params by call SerParams() function.")
	}

	sdk, err := fabsdk.New(config.FromFile(params.AdminConfigPath))
	if err != nil {
		log.Panicf("failed to create fabric sdk: %s", err)
	}

	rcp := sdk.Context(fabsdk.WithUser(params.OrgAdmin), fabsdk.WithOrg(params.OrgName))
	rc, err := resmgmt.New(rcp)
	if err != nil {
		log.Panicf("failed to create resource client: %s", err)
	}

	mc, err := msp.New(sdk.Context())
	if err != nil {
		log.Panicf("failed to create msp client: %s", err)
	}

	admin := new(AdminClient)

	admin.ClientParams = *params
	admin.SDK = sdk
	admin.RC = rc
	admin.MC = mc

	return admin, nil
}

func (admin *AdminClient) InitAppClient(channelID string) error {
	ccp := admin.SDK.ChannelContext(channelID, fabsdk.WithUser(admin.OrgUser))
	// ccp := admin.SDK.ChannelContext(channelID, fabsdk.WithUser(admin.OrgAdmin), fabsdk.WithUser(admin.OrgUser))
	cc, err := channel.New(ccp)
	if err != nil {
		return fmt.Errorf("failed to create channel client: %s", err)
	}

	ec, err := event.New(ccp, event.WithSeekType(seek.Newest))
	if err != nil {
		return fmt.Errorf("failed to create event client: %s", err)
	}

	lc, err := ledger.New(ccp)
	if err != nil {
		return fmt.Errorf("failed to create ledger client: %s", err)
	}

	admin.ChannelID = channelID
	admin.CC = cc
	admin.EC = ec
	admin.LC = lc

	return nil
}
