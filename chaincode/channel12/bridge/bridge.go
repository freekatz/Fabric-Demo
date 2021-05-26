package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create bridge chaincode: %s", err.Error())
		return
	}

	if err = chaincode.Start(); err != nil {
		fmt.Printf("Error starting bridge chaincode: %s", err.Error())
	}
}

type SmartContract struct {
	contractapi.Contract
}

type DigestResult struct {
	Key    string `json:"Key"` // h id
	Digest string `json:"digest"`
}

//
// 提供的功能包括：登记、更新、删除、验证
//

//
// 调用示例: '{"function":"register","Args":["h1","xxx"]}'
//
func (contract *SmartContract) Register(ctx contractapi.TransactionContextInterface, healthcareID, patientDigest string) error {

	if msp, _ := ctx.GetClientIdentity().GetMSPID(); msp != "Org2MSP" {
		return fmt.Errorf("Can not pass the identify with your MSP ID %s", msp)
	}

	digestAsBytes, err := ctx.GetStub().GetState(healthcareID)
	if err == nil && digestAsBytes != nil && len(digestAsBytes) > 0 {
		return fmt.Errorf("Info has been existed %s.", string(digestAsBytes))
	}
	digestAsBytes = []byte(patientDigest)

	return ctx.GetStub().PutState(healthcareID, digestAsBytes)
}

//
// 调用示例: '{"function":"update","Args":["h1","yyy"]}'
//
func (contract *SmartContract) Update(ctx contractapi.TransactionContextInterface, healthcareID, patientDigest string) error {

	if msp, _ := ctx.GetClientIdentity().GetMSPID(); msp != "Org2MSP" {
		return fmt.Errorf("Can not pass the identify with your MSP ID %s", msp)
	}

	digestAsBytes, err := ctx.GetStub().GetState(healthcareID)

	if err != nil || digestAsBytes == nil || len(digestAsBytes) == 0 {
		return err
	}

	digestAsBytes = []byte(patientDigest)
	return ctx.GetStub().PutState(healthcareID, digestAsBytes)
}

//
// 调用示例: '{"function":"query","Args":["h1"]}'
//
func (contract *SmartContract) Query(ctx contractapi.TransactionContextInterface, healthcareID string) (*DigestResult, error) {

	if msp, _ := ctx.GetClientIdentity().GetMSPID(); msp != "Org2MSP" {
		return nil, fmt.Errorf("Can not pass the identify with your MSP ID %s", msp)
	}

	digestAsBytes, err := ctx.GetStub().GetState(healthcareID)

	if err != nil || digestAsBytes == nil || len(digestAsBytes) == 0 {
		return nil, err
	}

	return &DigestResult{healthcareID, string(digestAsBytes)}, nil
}

//
// 调用示例: '{"function":"delete","Args":["h1"]}'
//
func (contract *SmartContract) Delete(ctx contractapi.TransactionContextInterface, healthcareID string) error {

	if msp, _ := ctx.GetClientIdentity().GetMSPID(); msp != "Org2MSP" {
		return fmt.Errorf("Can not pass the identify with your MSP ID %s", msp)
	}

	return ctx.GetStub().DelState(healthcareID)
}

//
// 调用示例: '{"function":"verify","Args":["h1","xxx"]}'
//
func (contract *SmartContract) Verify(ctx contractapi.TransactionContextInterface, healthcareID, patientDigest string) error {

	if msp, _ := ctx.GetClientIdentity().GetMSPID(); msp != "Org1MSP" {
		return fmt.Errorf("Can not pass the identify with your MSP ID %s", msp)
	}

	digestAsBytes, err := ctx.GetStub().GetState(healthcareID)

	if err != nil || digestAsBytes == nil || len(digestAsBytes) == 0 {
		return fmt.Errorf("Verify failed: no such info. %s", err.Error())
	}

	if string(digestAsBytes) != patientDigest {
		return fmt.Errorf("Verify failed: info can not match %s, %s.", string(digestAsBytes), patientDigest)
	}

	return nil
}
