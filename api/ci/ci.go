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

type Channel = gateway.Network

var (
	orgMSP   string
	orgHost  string
	ccpName  string
	userName string
)

var (
	basePath string
	credPath string
	certPath string
	ccpPath  string
)

func SetEnv(sw string) {
	os.Setenv("DISCOVERY_AS_LOCALHOST", sw)
}

func SetParams(params []string) {
	orgMSP = params[0]
	orgHost = params[1]
	ccpName = params[2]
	userName = params[3]
	basePath = filepath.Join(
		"..",
		"..",
		"..",
		"network",
		"orgs",
		"peerOrganizations",
		orgHost,
	)

	credPath = filepath.Join(
		basePath,
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

	channel, err := gw.GetNetwork(channelName)
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	return channel, nil
}
