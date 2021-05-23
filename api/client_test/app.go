/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/1uvu/Fabric-Demo/api/client"
	"github.com/1uvu/Fabric-Demo/structures"
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

	client.SetEnv("true")

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
	appConfigPath := filepath.Join(
		basePath,
		"app",
		appConfigName,
	)
	adminConfigPath := ""
	params := client.ClientParams{
		CredPath:        credPath,
		CertPath:        certPath,
		AppConfigPath:   appConfigPath,
		AdminConfigPath: adminConfigPath,
		OrgMSP:          orgMSP,
		OrgName:         orgName,
		OrgAdmin:        orgAdmin,
		OrgUser:         orgUser,
		OrgHost:         orgHost,
	}
	client.SetParams(&params)

	app1, err := client.GetAppClient("channel1")
	if err != nil {
		fmt.Printf("Failed to get app client: %s\n", err)
		os.Exit(1)
	}

	app12, err := client.GetAppClient("channel12")
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

	client.SetEnv("true")

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
	appConfigPath := ""
	adminConfigPath := filepath.Join(
		basePath,
		"admin",
		adminConfigName,
	)
	params := client.ClientParams{
		CredPath:        credPath,
		CertPath:        certPath,
		AppConfigPath:   appConfigPath,
		AdminConfigPath: adminConfigPath,
		OrgMSP:          orgMSP,
		OrgName:         orgName,
		OrgAdmin:        orgAdmin,
		OrgUser:         orgUser,
		OrgHost:         orgHost,
	}
	client.SetParams(&params)

	admin, err := client.GetAdminClient()

	if err != nil {
		fmt.Printf("Failed to get admin client: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(admin)

	admin.InitAppClient("channel123")

	fmt.Println(admin)
}
