package main

import (
	"encoding/json"
	"fmt"

	"github.com/1uvu/Fabric-Demo/structures"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ProduceRecord = structures.ProduceRecord

const CHANNEL_NAME = "channel123"
const CHAINCODE_NAME = "trace"

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create produce chaincode: %s", err.Error())
		return
	}

	if err = chaincode.Start(); err != nil {
		fmt.Printf("Error starting produce chaincode: %s", err.Error())
	}
}

type SmartContract struct {
	contractapi.Contract
}

//
// 提供的功能包括：登记、查询
//

//
// 调用示例: '{"function":"register","Args":["p1","t1","{\"name\":\"bj-P1\",\"producer\":\"bj\",\"address\":\"beijing\",\"date\":\"2021-05-01-12:00:00\",\"life\":\"-1\"}"]}'
//
func (contract *SmartContract) Register(ctx contractapi.TransactionContextInterface, produceID, traceID string, record ProduceRecord) error {
	// 注册时初始化并更新 trace record
	// todo 权限控制

	Args := [][]byte{[]byte("register"), []byte(traceID)}
	response := ctx.GetStub().InvokeChaincode(CHAINCODE_NAME, Args, CHANNEL_NAME)

	Args = [][]byte{[]byte("update"), []byte(traceID), []byte("produceID"), []byte(produceID)}
	response = ctx.GetStub().InvokeChaincode(CHAINCODE_NAME, Args, CHANNEL_NAME)

	if response.Status != 200 {
		return fmt.Errorf("Error %s with trace id %s, or", string(response.Payload), traceID)
	}

	recordAsBytes, _ := json.Marshal(record)

	return ctx.GetStub().PutState(produceID, recordAsBytes)
}

//
// 调用示例: '{"function":"query","Args":["p1"]}'
//
func (contract *SmartContract) Query(ctx contractapi.TransactionContextInterface, produceID string) (*ProduceRecord, error) {
	recordAsBytes, err := ctx.GetStub().GetState(produceID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	record := new(ProduceRecord)
	_ = json.Unmarshal(recordAsBytes, record)

	if record.Name == "" {
		return nil, fmt.Errorf("There have no produce record in ledger with produce id: %s", produceID)
	}

	return record, nil
}
