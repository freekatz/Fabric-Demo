/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/1uvu/Fabric-Demo/api/ci"
	"github.com/1uvu/Fabric-Demo/structures"
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

func main() {
	// test()
	ci.SetEnv("true")
	ci.SetParams([]string{
		"Org2MSP",
		"org2.example.com",
		"connection-org2.yaml",
		"User1",
	})
	channel2, err := ci.GetChannel("channel2")
	if err != nil {
		fmt.Printf("Failed to get channel: %s\n", err)
		os.Exit(1)
	}

	channel12, err := ci.GetChannel("channel12")
	if err != nil {
		fmt.Printf("Failed to get channel: %s\n", err)
		os.Exit(1)
	}

	patientContract := channel2.GetContract("patient")

	hid := "h2"
	patient := structures.NewPatientInHIB(
		[]string{
			"ZJH-2",
			"female",
			"1998-10-10",
			"abcdefghijklmnop",
			"139",
			"CQ",
			"NMG",
			"6674-1231",
		},
	)
	result, err := patientContract.SubmitTransaction("registerPatient", hid, patient.String())
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	result, err = patientContract.EvaluateTransaction("makeDigest", hid)
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	digest := new(DigestResult)
	_ = json.Unmarshal(result, digest)
	fmt.Println(digest.Digest)

	bridgeContract := channel12.GetContract("bridge")

	result, err = bridgeContract.EvaluateTransaction("register", hid, digest.Digest)
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))
}
