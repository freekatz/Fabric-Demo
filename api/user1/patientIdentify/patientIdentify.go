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

type Patient = structures.PatientInHOS

type QueryResult struct {
	Key    string `json: "Key"` // pat id
	Record *Patient
}

type DigestResult struct {
	Key    string `json: "Key"`
	Digest string `json: "digest"`
}

func main() {
	// test()
	ci.SetEnv("true")
	ci.SetParams([]string{
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

	patientContract := channel1.GetContract("patient")

	pid := "p2"
	patient := structures.NewPatientInHOS(
		[]string{
			"ZJH-2",
			"female",
			"1998-10-10",
			"15323211",
			"139",
			"CQ",
			"NMG",
			"6674-1231",
			"abcdefghijklmnop",
		},
	)
	result, err := patientContract.SubmitTransaction("registerPatient", pid, patient.String())
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	patientAsBytes, err := patientContract.EvaluateTransaction("queryPatient", pid)
	_patient := new(Patient)
	_ = json.Unmarshal(patientAsBytes, _patient)
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}

	result, err = patientContract.EvaluateTransaction("makeDigest", pid)
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	digest := new(DigestResult)
	_ = json.Unmarshal(result, digest)
	fmt.Println(digest.Digest)

	bridgeContract := channel12.GetContract("bridge")

	_, err = bridgeContract.EvaluateTransaction("verify", _patient.HealthcareID, digest.Digest)
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
}

// func test() {
// 	ci.SetEnv("true")
// 	ci.SetParams([]string{
// 		"Org1MSP",
// 		"org1.example.com",
// 		"connection-org1.yaml",
// 		"User1",
// 	})
// 	channelName1 := "channel1"
// 	channel1, err := ci.GetChannel(channelName1)
// 	if err != nil {
// 		fmt.Printf("Failed to get channel %s: %s\n", channelName1, err)
// 		os.Exit(1)
// 	}
// 	channelTest(channel1)
// }

// func channelTest(channel *ci.Channel) {
// 	contract := channel.GetContract("patient")

// 	result, err := contract.EvaluateTransaction("queryAllPatients")
// 	if err != nil {
// 		fmt.Printf("Failed to evaluate transaction: %s\n", err)
// 		os.Exit(1)
// 	}
// 	fmt.Println(string(result))

// 	pid := "p2"
// 	patient := structures.NewPatientInHOS(
// 		[]string{
// 			"ZJH-2",
// 			"female",
// 			"1998-10-10",
// 			"15323211",
// 			"139",
// 			"CQ",
// 			"NMG",
// 			"6674-1231",
// 			"abcdefghijklmnop",
// 		},
// 	)
// 	result, err = contract.SubmitTransaction("registerPatient", pid, patient.String())
// 	if err != nil {
// 		fmt.Printf("Failed to submit transaction: %s\n", err)
// 		os.Exit(1)
// 	}
// 	fmt.Println(string(result))

// 	result, err = contract.EvaluateTransaction("queryPatient", pid)
// 	if err != nil {
// 		fmt.Printf("Failed to evaluate transaction: %s\n", err)
// 		os.Exit(1)
// 	}
// 	fmt.Println(string(result))

// 	result, err = contract.EvaluateTransaction("makeDigest", "p1")
// 	if err != nil {
// 		fmt.Printf("Failed to evaluate transaction: %s\n", err)
// 		os.Exit(1)
// 	}
// 	fmt.Println(string(result))

// 	result, err = contract.EvaluateTransaction("makeDigest", pid)
// 	if err != nil {
// 		fmt.Printf("Failed to evaluate transaction: %s\n", err)
// 		os.Exit(1)
// 	}
// 	fmt.Println(string(result))
// }
