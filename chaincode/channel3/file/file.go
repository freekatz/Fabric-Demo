package main

import (
	"encoding/json"
	"fmt"

	"github.com/1uvu/Fabric-Demo/structures"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type FileRecord = structures.FileRecord

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create file chaincode: %s", err.Error())
		return
	}

	if err = chaincode.Start(); err != nil {
		fmt.Printf("Error starting file chaincode: %s", err.Error())
	}
}

type SmartContract struct {
	contractapi.Contract
}

type QueryResult struct {
	Key    string `json:"Key"` // pat id
	Record *FileRecord
}

func (contract *SmartContract) Register(ctx contractapi.TransactionContextInterface, fileRecordID string, fileRecord FileRecord) error {
	// todo 实现数据检查逻辑
	recordAsBytes, _ := json.Marshal(fileRecord)

	return ctx.GetStub().PutState(fileRecordID, recordAsBytes)
}

func (contract *SmartContract) Update(ctx contractapi.TransactionContextInterface, fileRecordID string, field string, value string) error {
	record, err := contract.Query(ctx, fileRecordID)

	if err != nil {
		return err
	}

	record.UpdateRecordField(field, value)

	recordAsBytes, _ := json.Marshal(record)

	return ctx.GetStub().PutState(fileRecordID, recordAsBytes)
}

func (contract *SmartContract) Query(ctx contractapi.TransactionContextInterface, fileRecordID string) (*FileRecord, error) {
	recordAsBytes, err := ctx.GetStub().GetState(fileRecordID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	record := new(FileRecord)
	_ = json.Unmarshal(recordAsBytes, record)

	if record.Timestamp == "" {
		return nil, fmt.Errorf("There no file in ledger with file record id: %s", fileRecordID)
	}

	return record, nil
}

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

		record := new(FileRecord)
		_ = json.Unmarshal(queryResponse.Value, record)

		queryResult := QueryResult{Key: queryResponse.Key, Record: record}
		results = append(results, queryResult)
	}

	return results, nil
}

func (contract *SmartContract) Delete(ctx contractapi.TransactionContextInterface, fileRecordID string) error {
	return ctx.GetStub().DelState(fileRecordID)
}
