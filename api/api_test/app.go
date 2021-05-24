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
	"github.com/1uvu/Fabric-Demo/api/app"
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
	orgName         string = "Org1"
	orgMSP          string = "Org1MSP"
	orgHost         string = "org1.example.com"
	appConfigName   string = "app-org1.yaml"
	adminConfigName string = "admin-org1.yaml"
	orgUser         string = "User1"
	orgAdmin        string = "Admin"
)

func main() {
	testApp()
	testAdmin()
}

func testApp() {
	fmt.Println("testing app client")

	app.SetEnv("true")

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
		appConfigName,
	)
	params := app.ClientParams{
		CredPath:   credPath,
		CertPath:   certPath,
		ConfigPath: configPath,
		OrgMSP:     orgMSP,
		OrgName:    orgName,
		OrgAdmin:   orgAdmin,
		OrgUser:    orgUser,
		OrgHost:    orgHost,
	}
	app.SetParams(&params)

	app1, err := app.GetAppClient("channel1")
	if err != nil {
		fmt.Printf("Failed to get app client: %s\n", err)
		os.Exit(1)
	}

	app12, err := app.GetAppClient("channel12")
	if err != nil {
		fmt.Printf("Failed to get app client: %s\n", err)
		os.Exit(1)
	}

	patientChaincode := app1.GetContract("patient")

	pid := "p3"
	patient := structures.NewPatientInHOS(
		[]string{
			"ZJH-3",
			"female",
			"2020-10-10",
			"abcdefghijklmnop",
			"151",
			"CQU",
			"NMG",
			"6674-1231-1000",
			"h3",
		},
	)

	// result, err := patientChaincode.SubmitTransaction("registerPatient", []string{pid, patient.String()}...)
	// if err != nil {
	// 	fmt.Printf("Failed to submit transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))

	result, err := patientChaincode.EvaluateTransaction("makeDigest", []string{pid}...)
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	digest := new(DigestResult)
	_ = json.Unmarshal(result, digest)
	fmt.Println(digest.Digest)

	bridgeChaincode := app12.GetContract("bridge")

	// _, err = bridgeChaincode.SubmitTransaction("register", []string{patient.HealthcareID, digest.Digest}...)
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }

	_, err = bridgeChaincode.EvaluateTransaction("verify", []string{patient.HealthcareID, digest.Digest}...)
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
}

func testAdmin() {
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
		adminConfigName,
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

	// 写账本
	// new channel request for invoke
	args := [][]byte{[]byte("t10")}
	req := channel.Request{
		ChaincodeID: "trace",
		Fcn:         "register",
		Args:        args,
	}

	// send request and handle response
	// peers is needed
	reqPeers := channel.WithTargetEndpoints(
		"peer0.org1.example.com",
		"peer0.org2.example.com",
		"peer0.org3.example.com",
	)
	// 可不指定 peers 使用默认，如需指定则需要符合 chaincode 的背书策略
	resp, err := app123.CC.Execute(req, reqPeers)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("invoke chaincode tx: %s", resp.TransactionID)

	// 读账本
	args = [][]byte{[]byte("t1")}
	req = channel.Request{
		ChaincodeID: "trace",
		Fcn:         "query",
		Args:        args,
	}
	resp, err = app123.CC.Query(req)

	if err != nil {
		fmt.Println(err)
	}
	log.Printf("invoke chaincode tx: %s", resp.TransactionID)
	log.Printf("resp content %s", string(resp.Payload))

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
	log.Printf("resp content %s", string(resp.Payload))
}
