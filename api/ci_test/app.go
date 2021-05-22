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

	"github.com/1uvu/Fabric-Demo/api/ci"
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

func main() {
	// test()
	ci.SetEnv("true")
	basePath := filepath.Join(
		"..",
		"..",
		"network",
		"orgs",
		"peerOrganizations",
	)
	ci.SetParams(&ci.CIParams{
		basePath,
		"Org1MSP",
		"org1.example.com",
		"connection-org1.yaml",
		"User1",
	})
	channel1, err := ci.GetChannel("channel1")
	if err != nil {
		fmt.Printf("Failed to get channel: %s\n", err)
		os.Exit(1)
	}

	channel12, err := ci.GetChannel("channel12")
	if err != nil {
		fmt.Printf("Failed to get channel: %s\n", err)
		os.Exit(1)
	}

	patientChaincode := channel1.GetChaincode("patient")

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
	result, err := patientChaincode.InvokeContract(&ci.InvokeParams{"registerPatient", []string{pid, patient.String()}}, true)
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))
	result, err = patientChaincode.InvokeContract(&ci.InvokeParams{"makeDigest", []string{pid}}, false)
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	digest := new(DigestResult)
	_ = json.Unmarshal(result, digest)
	fmt.Println(digest.Digest)

	bridgeChaincode := channel12.GetChaincode("bridge")

	_, err = bridgeChaincode.InvokeContract(&ci.InvokeParams{"register", []string{patient.HealthcareID, digest.Digest}}, true)
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}

	_, err = bridgeChaincode.InvokeContract(&ci.InvokeParams{"verify", []string{patient.HealthcareID, digest.Digest}}, false)
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
}
