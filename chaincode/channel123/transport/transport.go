package main

import (
	"encoding/json"
	"fmt"

	"github.com/1uvu/Fabric-Demo/structures"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type TransportRecord = structures.TransportRecord

const CHANNEL_NAME = "channel123"
const CHAINCODE_NAME = "trace"

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create transport chaincode: %s", err.Error())
		return
	}

	if err = chaincode.Start(); err != nil {
		fmt.Printf("Error starting transport chaincode: %s", err.Error())
	}
}

type SmartContract struct {
	contractapi.Contract
}

//
// 提供的功能包括：登记、查询
//

//
// 调用示例: '{"function":"register","Args":["tt1","t1","{\"transporter\":\"Shunfeng\",\"originAddress\":\"chongqing\",\"targetAddress\":\"neimenggu\",\"startDate\":\"2021-05-11-7:00:00\",\"endDate\":\"2021-05-15-15:30:00\"}"]}'
//
func (contract *SmartContract) Register(ctx contractapi.TransactionContextInterface, transportID, traceID string, record TransportRecord) error {
	// 注册时更新 trace record

	if msp, _ := ctx.GetClientIdentity().GetMSPID(); msp != "Org3MSP" {
		return fmt.Errorf("Can not pass the identify with your MSP ID %s", msp)
	}

	Args := [][]byte{[]byte("update"), []byte(traceID), []byte("transportID"), []byte(transportID)}
	response := ctx.GetStub().InvokeChaincode(CHAINCODE_NAME, Args, CHANNEL_NAME)

	if response.Status != 200 {
		return fmt.Errorf("Error %s with trace id %s, or", string(response.Payload), traceID)
	}

	recordAsBytes, _ := json.Marshal(record)

	return ctx.GetStub().PutState(transportID, recordAsBytes)
}

//
// 调用示例: '{"function":"query","Args":["tt1"]}'
//
func (contract *SmartContract) Query(ctx contractapi.TransactionContextInterface, transportID string) (*TransportRecord, error) {
	recordAsBytes, err := ctx.GetStub().GetState(transportID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	record := new(TransportRecord)
	_ = json.Unmarshal(recordAsBytes, record)

	if record.Transporter == "" {
		return nil, fmt.Errorf("There have no transport record in ledger with transport id: %s", transportID)
	}

	return record, nil
}
