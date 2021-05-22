package ci

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

//
// 封装定义数据类型
//

// 包装为 Channel 类型
type Channel struct {
	*gateway.Network
}

// 包装为 Chaincode 类型
type Chaincode struct {
	*gateway.Contract
}

// 用于连接 Fabric 网络的参数
type CIParams struct {
	BasePath string
	OrgMSP   string
	OrgHost  string
	CCPName  string
	UserName string
}

// 用于调用合约的参数
type InvokeParams struct {
	ContractName string // 链码中的函数名
	Args         []string
}

//
// 定义全局变量
//

var (
	basePath string
	orgMSP   string
	orgHost  string
	ccpName  string
	userName string
)

var (
	credPath string
	certPath string
	ccpPath  string
)

func SetEnv(sw string) {
	os.Setenv("DISCOVERY_AS_LOCALHOST", sw)
}

func SetParams(params *CIParams) {
	basePath = params.BasePath
	orgMSP = params.OrgMSP
	orgHost = params.OrgHost
	ccpName = params.CCPName
	userName = params.UserName

	credPath = filepath.Join(
		basePath,
		orgHost,
		"users",
		fmt.Sprintf("%s@%s", userName, orgHost),
		"msp",
	)

	certPath = filepath.Join(
		credPath,
		"signcerts",
		fmt.Sprintf("%s@%s-cert.pem", userName, orgHost),
	)
	ccpPath = filepath.Join(
		basePath,
		orgHost,
		ccpName,
	)
}

//
// 存储认证材料
//
func populateWallet(wallet *gateway.Wallet) error {
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	// there's a single file in this dir containing the private key
	keyDir := filepath.Join(credPath, "keystore")
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

	identity := gateway.NewX509Identity(orgMSP, string(cert), string(key))

	err = wallet.Put("appUser", identity)
	if err != nil {
		return err
	}
	return nil
}

//
// 获取通道实例
//
func GetChannel(channelName string) (*Channel, error) {
	if basePath == "" {
		return nil, fmt.Errorf("please make sure you have invoked the set params functions!")
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
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
	defer gw.Close()

	network, err := gw.GetNetwork(channelName)
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	return &Channel{network}, nil
}

//
// 获取链码实例
//
func (channel *Channel) GetChaincode(chaincodeName string) *Chaincode {
	return &Chaincode{channel.GetContract(chaincodeName)}
}

//
// 调用链码中的合约
//
func (chaincode *Chaincode) InvokeContract(params *InvokeParams, rw bool) (result []byte, err error) {
	if rw {
		result, err = chaincode.SubmitTransaction(params.ContractName, params.Args...)
	} else {
		result, err = chaincode.EvaluateTransaction(params.ContractName, params.Args...)
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}
