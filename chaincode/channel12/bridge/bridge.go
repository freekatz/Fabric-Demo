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

//
// 提供的功能包括：登记、验证
//

//
// 调用示例: '{"function":"register","Args":["ABCDEFGHIJKLMNOP","xxx"]}'
//
func (contract *SmartContract) Register(ctx contractapi.TransactionContextInterface, healthcareID, patientDigest string) error {
	// todo 只能由 org2 调用
	// todo 不允许重复插入

	digestAsBytes := []byte(patientDigest)

	return ctx.GetStub().PutState(healthcareID, digestAsBytes)
}

//
// 调用示例: '{"function":"verify","Args":["ABCDEFGHIJKLMNOP","xxx"]}'
//
func (contract *SmartContract) Verify(ctx contractapi.TransactionContextInterface, healthcareID, patientDigest string) error {
	// todo 只能由 org1 调用

	digestAsBytes, err := ctx.GetStub().GetState(healthcareID)

	if err != nil || len(digestAsBytes) == 0 {
		return fmt.Errorf("Verify failed: no such info. %s", err.Error())
	}

	if string(digestAsBytes) != patientDigest {
		return fmt.Errorf("Verify failed: info can not match %s, %s.", string(digestAsBytes), patientDigest)
	}

	return nil
}
