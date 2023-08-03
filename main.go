package main

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	chaincode "github.com/rohinsood/roblockain/m/v2/chaincode"
)

func main() {
	chaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		fmt.Printf("Error creating log file chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting log file chaincode: %s", err.Error())
	}
}