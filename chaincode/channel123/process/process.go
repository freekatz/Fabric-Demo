package main

import (
	"encoding/json"
	"fmt"

	"github.com/1uvu/Fabric-Demo/structures"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ProcessRecord = structures.ProcessRecord

const CHANNEL_NAME = "channel123"
const CHAINCODE_NAME = "trace"

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create process chaincode: %s", err.Error())
		return
	}

	if err = chaincode.Start(); err != nil {
		fmt.Printf("Error starting process chaincode: %s", err.Error())
	}
}

type SmartContract struct {
	contractapi.Contract
}

//
// 提供的功能包括：登记、查询
//

//
// 调用示例: '{"function":"register","Args":["pp1","t1","{\"name\":\"cq-PP1\",\"type\":\"machine\",\"processor\":\"cq\",\"address\":\"chongqing\",\"date\":\"2021-05-10-23:00:00\",\"life\":\"-1\"}"]}'
//
func (contract *SmartContract) Register(ctx contractapi.TransactionContextInterface, processID, traceID string, record ProcessRecord) error {
	// 注册时更新 trace record
	// todo 权限控制

	Args := [][]byte{[]byte("update"), []byte(traceID), []byte("processID"), []byte(processID)}
	response := ctx.GetStub().InvokeChaincode(CHAINCODE_NAME, Args, CHANNEL_NAME)

	if response.Status != 200 {
		return fmt.Errorf("Error %s with trace id %s, or", string(response.Payload), traceID)
	}

	recordAsBytes, _ := json.Marshal(record)

	return ctx.GetStub().PutState(processID, recordAsBytes)
}

//
// 调用示例: '{"function":"query","Args":["pp1"]}'
//
func (contract *SmartContract) Query(ctx contractapi.TransactionContextInterface, processID string) (*ProcessRecord, error) {
	recordAsBytes, err := ctx.GetStub().GetState(processID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	record := new(ProcessRecord)
	_ = json.Unmarshal(recordAsBytes, record)

	if record.Name == "" {
		return nil, fmt.Errorf("There have no process record in ledger with process id: %s", processID)
	}

	return record, nil
}
