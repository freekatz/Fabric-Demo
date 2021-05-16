package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/1uvu/Fabric-Demo/crypt"
	"github.com/1uvu/Fabric-Demo/structures"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type GenderType = structures.GenderType

type Patient = structures.PatientInHIB

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create patient chaincode: %s", err.Error())
		return
	}

	if err = chaincode.Start(); err != nil {
		fmt.Printf("Error starting patient chaincode: %s", err.Error())
	}
}

type SmartContract struct {
	contractapi.Contract
}

type QueryResult struct {
	Key    string `json:"Key"` // pat id
	Record *Patient
}

type DigestResult struct {
	Key    string `json:"Key"`
	Digest string `json: "digest"`
}

//
// 提供的功能包括：登记、更新、查询、以及删除
//

//
// 调用示例: '{"function":"registerPatient","Args":["ABCDEFGHIJKLMNOP", "{\"name\": \"ZJH-1\", \"gender\": \"male\", \"birth\": \"1998-10-01\", \"identifyID\": \"xxxxxx-xxxx-19981001-xxxx-xxxx\", \"phoneNumber\": \"151-2300-0000\", \"address\": \"ChongQing\", \"nativePlace\": \"NeiMengGu\", \"creditCard\": \"6217-0000-0000-0000\"}"]}'
//
func (contract *SmartContract) RegisterPatient(ctx contractapi.TransactionContextInterface, healthcareID string, patient Patient) error {
	// todo 实现数据检查逻辑
	patientAsBytes, _ := json.Marshal(patient)

	return ctx.GetStub().PutState(healthcareID, patientAsBytes)
}

//
// 调用示例: '{"function":"updatePatient","Args":["ABCDEFGHIJKLMNOP", "[\"name\", \"gender\"]", "[\"ZJH-2\", \"female\"]"]}'
//
func (contract *SmartContract) UpdatePatient(ctx contractapi.TransactionContextInterface, healthcareID string, fields []string, values []interface{}) error {
	patient, err := contract.QueryPatient(ctx, healthcareID)

	if err != nil {
		return err
	}

	if len(fields) != len(values) {
		return fmt.Errorf("len of fields and values are not equal.")
	}

	for i := range fields {
		f, v := fields[i], values[i]
		patient.UpdatePatientField(f, v)
	}

	patientAsBytes, _ := json.Marshal(patient)

	return ctx.GetStub().PutState(healthcareID, patientAsBytes)
}

//
// 调用示例: '{"function":"queryPatient","Args":["ABCDEFGHIJKLMNOP"]}'
//
func (contract *SmartContract) QueryPatient(ctx contractapi.TransactionContextInterface, healthcareID string) (*Patient, error) {
	patientAsBytes, err := ctx.GetStub().GetState(healthcareID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	patient := new(Patient)
	_ = json.Unmarshal(patientAsBytes, patient)

	if patient.Name == "" {
		return nil, fmt.Errorf("There no patient in ledger with healthcare id: %s", healthcareID)
	}

	return patient, nil
}

//
// 调用示例: '{"function":"queryAllPatients","Args":[]}'
// QueryAllCars returns all patients found in world state
func (s *SmartContract) QueryAllPatients(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
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

		patient := new(Patient)
		_ = json.Unmarshal(queryResponse.Value, patient)

		queryResult := QueryResult{Key: queryResponse.Key, Record: patient}
		results = append(results, queryResult)
	}

	return results, nil
}

//
// 调用示例: '{"function":"deletePatient","Args":["ABCDEFGHIJKLMNOP"]}'
//
func (contract *SmartContract) DeletePatient(ctx contractapi.TransactionContextInterface, healthcareID string) error {
	return ctx.GetStub().DelState(healthcareID)
}

//
// 调用示例: '{"function":"makeDigest","Args":["ABCDEFGHIJKLMNOP"]}'
//
func (contract *SmartContract) MakeDigest(ctx contractapi.TransactionContextInterface, healthcareID string) (*DigestResult, error) {
	patientAsBytes, err := ctx.GetStub().GetState(healthcareID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	patient := new(Patient)
	_ = json.Unmarshal(patientAsBytes, patient)

	// 计算密文和摘要
	cryptoAsBytes := crypt.AesEncryptCBC([]byte(patient.IdentifyID), []byte(healthcareID))
	patientDigest := base64.StdEncoding.Strict().EncodeToString(cryptoAsBytes)

	// 返回使用 hid 加密的 iid
	return &DigestResult{healthcareID, patientDigest}, nil
}
