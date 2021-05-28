/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/1uvu/Fabric-Demo/structures"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

type GenderType = structures.GenderType

type Patient = structures.PatientInHIB

type QueryResult struct {
	Key    string `json: "Key"` // pat id
	Record *Patient
}

type DigestResult struct {
	Key    string `json:"Key"` // h id
	Digest string `json:"digest"`
}

//
// 定义全局变量
//

var (
	basePath string = filepath.Join(
		"..",
		"..",
		"network",
		"orgs",
	)
	orgName    string = "Org2"
	orgMSP     string = "Org2MSP"
	orgHost    string = "org2.example.com"
	configName string = ""
	orgUser    string = "User1"
	orgAdmin   string = "Admin"
)

func Test(t *testing.T) {
	appTest()
	adminTest()
}

func appTest() {
	fmt.Println("testing app client")
	configName = "app-org2.yaml"

	SetAppEnv("true")

	credPath := filepath.Join(
		basePath,
		"peerOrganizations",
		orgHost,
		"users",
		fmt.Sprintf("%s@%s", orgUser, orgHost),
		"msp",
	)
	certPath := filepath.Join(
		credPath,
		"signcerts",
		fmt.Sprintf("%s@%s-cert.pem", orgUser, orgHost),
	)
	configPath := filepath.Join(
		basePath,
		"app",
		configName,
	)
	params := AppParams{ // todo 这里是指针，所以，初始化时直接变更即可
		CredPath:   credPath,
		CertPath:   certPath,
		ConfigPath: configPath,
		OrgMSP:     orgMSP,
		OrgName:    orgName,
		OrgAdmin:   orgAdmin,
		OrgUser:    orgUser,
		OrgHost:    orgHost,
	}
	SetAppParams(&params)

	app2, err := GetAppClient("channel2")
	if err != nil {
		fmt.Printf("Failed to get app client: %s\n", err)
		os.Exit(1)
	}

	app12, err := GetAppClient("channel12")
	if err != nil {
		fmt.Printf("Failed to get app client: %s\n", err)
		os.Exit(1)
	}

	hid := "h3"
	patient := structures.NewPatientInHIB(
		[]string{
			"ZJH-3",
			"female",
			"2020-10-10",
			"abcdefghijklmnop",
			"151",
			"CQU",
			"NMG",
			"6674-1231-1000",
		},
	)

	patientChaincode := app2.GetContract("patient")

	result, err := patientChaincode.SubmitTransaction("register", []string{hid, patient.String()}...)
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
	}

	result, err = patientChaincode.EvaluateTransaction("makeDigest", []string{hid}...)
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
	}
	digest := new(DigestResult)
	_ = json.Unmarshal(result, digest)
	log.Println(digest.Digest)

	bridgeChaincode := app12.GetContract("bridge")

	_, err = bridgeChaincode.SubmitTransaction("register", []string{hid, digest.Digest}...)
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
	}
}

func adminTest() {
	fmt.Println("testing admin client")
	configName = "admin-org2.yaml"

	SetAdminEnv("true")

	credPath := filepath.Join(
		basePath,
		"peerOrganizations",
		orgHost,
		"users",
		fmt.Sprintf("%s@%s", orgAdmin, orgHost),
		"msp",
	)
	certPath := filepath.Join(
		credPath,
		"signcerts",
		fmt.Sprintf("%s@%s-cert.pem", orgUser, orgHost),
	)
	configPath := filepath.Join(
		basePath,
		"admin",
		configName,
	)
	params := AdminParams{
		CredPath:   credPath,
		CertPath:   certPath,
		ConfigPath: configPath,
		OrgMSP:     orgMSP,
		OrgName:    orgName,
		OrgAdmin:   orgAdmin,
		OrgUser:    orgUser,
		OrgHost:    orgHost,
	}
	SetAdminParams(&params)

	admin1, err := GetAdminClient()

	if err != nil {
		fmt.Printf("Failed to get admin1 client: %s\n", err)
		os.Exit(1)
	}

	app123, _ := admin1.GetAppClient("channel123")

	// 读账本
	args := [][]byte{[]byte("t1")}
	req := channel.Request{
		ChaincodeID: "trace",
		Fcn:         "query",
		Args:        args,
	}
	resp, err := app123.CC.Query(req)

	if err != nil {
		fmt.Println(err)
	}
	log.Printf("invoke chaincode tx: %s", resp.TransactionID)
	log.Printf("resp content %s", string(resp.Payload))

	// 写账本
	// new channel request for invoke
	args = [][]byte{[]byte("t10")}
	req = channel.Request{
		ChaincodeID: "trace",
		Fcn:         "register",
		Args:        args,
	}

	// send request and handle response
	reqPeers := channel.WithTargetEndpoints(
		"peer0.org1.example.com",
		"peer0.org2.example.com", // 三者至少包括两个即可
		// "peer0.org3.example.com",
	)
	// 可不指定 peers 使用默认，如需指定则需要符合 chaincode 的背书策略
	resp, err = app123.CC.Execute(req, reqPeers)
	// resp, err := app123.CC.Execute(req)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("invoke chaincode tx: %s", resp.TransactionID)

	// 获取其他通道的 app client
	app2, _ := admin1.GetAppClient("channel2")
	args = [][]byte{[]byte("h3")}
	req = channel.Request{
		ChaincodeID: "patient",
		Fcn:         "query",
		Args:        args,
	}
	resp, err = app2.CC.Query(req)

	if err != nil {
		fmt.Println(err)
	}
	log.Printf("invoke chaincode tx: %s", resp.TransactionID)
	log.Printf("resp content:\n%s", string(resp.Payload))

	// 这里获取到的 _app2 与 app2 是同一个 app client
	_app2, _ := admin1.GetAppClient("channel2")
	args = [][]byte{}
	req = channel.Request{
		ChaincodeID: "patient",
		Fcn:         "queryAll",
		Args:        args,
	}
	resp, err = _app2.CC.Query(req)

	if err != nil {
		fmt.Println(err)
	}
	log.Printf("invoke chaincode tx: %s", resp.TransactionID)

	var patients []QueryResult
	_ = json.Unmarshal(resp.Payload, &patients)
	log.Println("All patients are as follows:")
	fmt.Println(patients)
}
