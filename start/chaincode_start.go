/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"strings"
)

type Order struct {
	id         int
	itemsId    []string
	customerId int
	status     OrderStatus
}
type OrderStatus int

const (
	Issued OrderStatus = iota
	Shipped
	ArrivedToPO
	Finished
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Running write")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	id := args[0]
	var order Order
	order.id, err = strconv.Atoi(id)
	order.itemsId = strings.Split(args[1], ",")
	order.customerId, err = strconv.Atoi(args[2])
	order.status = Issued
	orderBytes, err := json.Marshal(order)
	err = stub.PutState(id, orderBytes) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("Running delete")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	key := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(key)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "delete" {
		return t.delete(stub, args)
	} else if function == "write" {
		return t.write(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Query called, determining function")

	if function != "query" {
		fmt.Printf("Function is query")
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var id string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	id = args[0]

	// Get the state from the ledger
	order, err := stub.GetState(id)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + id + "\"}"
		return nil, errors.New(jsonResp)
	}

	if order == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + id + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + id + "\",\"Order\":\"" + string(order) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return order, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
