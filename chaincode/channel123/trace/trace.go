package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/1uvu/Fabric-Demo/structures"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type TraceRecord = structures.TraceRecord
type TraceHistory = structures.TraceHistory

const TIME_LAYOUT = "2006-01-02 15:04:05"

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create trace chaincode: %s", err.Error())
		return
	}

	if err = chaincode.Start(); err != nil {
		fmt.Printf("Error starting trace chaincode: %s", err.Error())
	}
}

type SmartContract struct {
	contractapi.Contract
}

type QueryResult struct {
	Key    string `json:"Key"` // trace id
	Record *TraceRecord
}

type HistoryResult struct {
	Key     string         `json:"Key"`     // trace id
	History []TraceHistory `json:"history"` // history json str
}

//
// 提供的功能包括：初始化、登记、更新、查询、以及删除
//

//
// 调用示例: '{"function":"register","Args":["t1"]}'
//
func (contract *SmartContract) Register(ctx contractapi.TransactionContextInterface, traceID string) error {

	// 只有生产商才可注册

	if msp, _ := ctx.GetClientIdentity().GetMSPID(); msp != "Org1MSP" {
		return fmt.Errorf("Can not pass the identify with your MSP ID %s", msp)
	}

	// 注册时默认自动生成空对象
	_r, err := contract.Query(ctx, traceID)

	if err != nil {
		return err
	}
	_v, err := _r.GetTraceRecordValue("produceID")
	if _v != "" && err == nil {
		return fmt.Errorf("Can not register twice.")
	}

	record := structures.NewTraceRecord([]string{"", "", ""})
	recordAsBytes, _ := json.Marshal(&record)

	return ctx.GetStub().PutState(traceID, recordAsBytes)
}

//
// 调用示例: '{"function":"update","Args":["t1", "produceID", "p1"]}'
//
func (contract *SmartContract) Update(ctx contractapi.TransactionContextInterface, traceID, field, value string) error {

	// 只有生产商才可更新生产商的对应记录

	_msp := ""
	switch field {
	case "produceID":
		_msp = "Org1MSP"
	case "processID":
		_msp = "Org2MSP"
	case "transportID":
		_msp = "Org3MSP"
	default:
		return fmt.Errorf("Unknow field named %s", field)
	}

	if msp, _ := ctx.GetClientIdentity().GetMSPID(); msp != _msp {
		return fmt.Errorf("Can not pass the identify with your MSP ID %s", msp)
	}

	record, err := contract.Query(ctx, traceID)

	if err != nil {
		return err
	}

	// 只允许一次更新一个记录, 且相同角色只有一次更新的机会 (防篡改)

	_v, err := record.GetTraceRecordValue(field)
	if _v != "" && err == nil {
		return fmt.Errorf("Can not update twice.")
	}

	record.UpdateTraceRecordField(field, value)

	recordAsBytes, _ := json.Marshal(record)

	return ctx.GetStub().PutState(traceID, recordAsBytes)
}

//
// 调用示例: '{"function":"query","Args":["t1"]}'
//
func (contract *SmartContract) Query(ctx contractapi.TransactionContextInterface, traceID string) (*TraceRecord, error) {
	recordAsBytes, err := ctx.GetStub().GetState(traceID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	record := new(TraceRecord)
	_ = json.Unmarshal(recordAsBytes, record)

	return record, nil
}

//
// 调用示例: '{"function":"queryAll","Args":[]}'
//
func (s *SmartContract) QueryAll(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		record := new(TraceRecord)
		_ = json.Unmarshal(queryResponse.Value, record)

		queryResult := QueryResult{Key: queryResponse.Key, Record: record}
		results = append(results, queryResult)
	}

	return results, nil
}

//
// 调用示例: '{"function":"queryHistory","Args":["t1"]}'
//
func (s *SmartContract) QueryHistory(ctx contractapi.TransactionContextInterface, traceID string) (*HistoryResult, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(traceID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	// 通过迭代器对象遍历结果
	var history []TraceHistory
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		txID := queryResponse.TxId
		txValue := queryResponse.Value
		txTime := queryResponse.Timestamp
		txStatus := queryResponse.IsDelete
		txTimeStr := time.Unix(txTime.Seconds, 0).Format(TIME_LAYOUT)
		history = append(
			history,
			structures.NewTraceHistory([]string{
				txID, string(txValue),
				txTimeStr,
				fmt.Sprintf("%t", txStatus),
			},
			))
	}
	return &HistoryResult{traceID, history}, nil
}
