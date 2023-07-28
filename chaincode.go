package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type LogFile struct {
	FileID 		 string
	Content    []byte
	Authorized []string // List of authorized organizations/identities
}

func (sc *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	dataCSVBytes, err := ioutil.ReadFile("data.csv")
	if (err != nil) {
		return fmt.Errorf("Error reading data.csv: %v", err)
	}

	err = ctx.GetStub().PutState("1", dataCSVBytes)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	return nil
}

func (sc *SmartContract) UploadLogFile(ctx contractapi.TransactionContextInterface, fileID string, content []byte) error {

	currentMSPID, err := ctx.GetClientIdentity().GetMSPID()

	if (err != nil) {
		return err
	}

	logFile := &LogFile{
		FileID: 		fileID,
		Content:    content,
		Authorized: []string{currentMSPID}, // Add the MSP ID of the current client as an authorized organization
	}

	logFileBytes, err := json.Marshal(logFile)

	if (err != nil) {
		return err
	}

	return ctx.GetStub().PutState(fileID, logFileBytes)
}

func (sc *SmartContract) ReadLogFile(ctx contractapi.TransactionContextInterface, fileID string) (*LogFile, error) {
	logFileBytes, err := ctx.GetStub().GetState(fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	if logFileBytes == nil {
		return nil, fmt.Errorf("log file with ID '%s' does not exist", fileID)
	}

	logFile := &LogFile{}
	err = json.Unmarshal(logFileBytes, logFile)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal log file data: %w", err)
	}

	// Check if the client's MSP ID is in the authorized list
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if !contains(logFile.Authorized, clientMSPID) {
		return nil, fmt.Errorf("unauthorized access to log file")
	}

	return logFile, nil
}

// Helper function to check if a value exists in a string slice
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		fmt.Printf("Error creating log file chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting log file chaincode: %s", err.Error())
	}
}
