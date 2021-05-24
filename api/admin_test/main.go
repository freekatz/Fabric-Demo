/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/1uvu/Fabric-Demo/api/admin"
	"github.com/1uvu/Fabric-Demo/structures"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

type GenderType = structures.GenderType

type Patient = structures.PatientInHOS

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
	orgName    string = "Org1"
	orgMSP     string = "Org1MSP"
	orgHost    string = "org1.example.com"
	configName string = "admin-org1.yaml"
	orgUser    string = "User1"
	orgAdmin   string = "Admin"
)

func main() {
	fmt.Println("testing admin client")

	admin.SetEnv("true")

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
	params := admin.ClientParams{
		CredPath:   credPath,
		CertPath:   certPath,
		ConfigPath: configPath,
		OrgMSP:     orgMSP,
		OrgName:    orgName,
		OrgAdmin:   orgAdmin,
		OrgUser:    orgUser,
		OrgHost:    orgHost,
	}
	admin.SetParams(&params)

	admin1, err := admin.GetAdminClient()

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
	app1, _ := admin1.GetAppClient("channel1")
	args = [][]byte{[]byte("p1")}
	req = channel.Request{
		ChaincodeID: "patient",
		Fcn:         "queryPatient",
		Args:        args,
	}
	resp, err = app1.CC.Query(req)

	if err != nil {
		fmt.Println(err)
	}
	log.Printf("invoke chaincode tx: %s", resp.TransactionID)
	log.Printf("resp content:\n%s", string(resp.Payload))

	// 这里获取到的 _app1 与 app1 是同一个 app client
	_app1, _ := admin1.GetAppClient("channel1")
	args = [][]byte{}
	req = channel.Request{
		ChaincodeID: "patient",
		Fcn:         "queryAllPatients",
		Args:        args,
	}
	resp, err = _app1.CC.Query(req)

	if err != nil {
		fmt.Println(err)
	}
	log.Printf("invoke chaincode tx: %s", resp.TransactionID)

	var patients []QueryResult
	_ = json.Unmarshal(resp.Payload, &patients)
	log.Println("All patients are as follows:")
	fmt.Println(patients)
}
