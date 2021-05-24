package admin

import (
	"fmt"
	"log"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

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
	ChannelID string

	CC *channel.Client
	EC *event.Client
	LC *ledger.Client
}

type AdminClient struct {
	ClientParams

	// sdk clients
	SDK *fabsdk.FabricSDK
	RC  *resmgmt.Client
	MC  *msp.Client

	acs map[string]*AppClient

	// for create channel
	// ChannelConfig string
	// OrdererID string
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

func GetAdminClient() (*AdminClient, error) {
	if params.ConfigPath == "" {
		return nil, fmt.Errorf("Please init the params by call SerParams() function.")
	}

	sdk, err := fabsdk.New(config.FromFile(params.ConfigPath))
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
	admin.acs = make(map[string]*AppClient)

	return admin, nil
}

func (admin *AdminClient) GetAppClient(channelID string) (*AppClient, error) {

	if app, ok := admin.acs[channelID]; ok {
		return app, nil
	}

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

	app := new(AppClient)
	app.ChannelID = channelID
	app.CC = cc
	app.EC = ec
	app.LC = lc

	admin.acs[channelID] = app

	return app, nil
}
