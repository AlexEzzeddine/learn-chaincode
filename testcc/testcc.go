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
	Id         int
	ItemsId    []string
	CustomerId string
	Status     string
}

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

func (t *SimpleChaincode) SubmitOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Running write")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. name of the key and value to set")
	}

	id := args[0]
	var order Order
	order.Id, err = strconv.Atoi(id)
	order.ItemsId = strings.Split(args[1], ",")
	order.CustomerId = args[2]
	order.Status = "Issued"
	orderBytes, err := json.Marshal(order)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = stub.PutState(id, orderBytes) //write the variable into the chaincode state
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) EditOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var order Order
	fmt.Println("Changing order")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}
	id := args[0]
	orderBytes, err := stub.GetState(id)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + id + "\"}"
		return nil, errors.New(jsonResp)
	}

	err = json.Unmarshal(orderBytes, &order)

	order.ItemsId = strings.Split(args[1], ",")
	orderBytes, err = json.Marshal(order)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = stub.PutState(id, orderBytes) //write the variable into the chaincode state
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) ChangeStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var order Order
	fmt.Println("Changing status")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}
	id := args[0]
	orderBytes, err := stub.GetState(id)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + id + "\"}"
		return nil, errors.New(jsonResp)
	}

	err = json.Unmarshal(orderBytes, &order)

	order.Status = args[1]
	orderBytes, err = json.Marshal(order)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = stub.PutState(id, orderBytes) //write the variable into the chaincode state
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return nil, nil
}

// Deletes an entity from state
func (t *SimpleChaincode) CancelOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("Running delete")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
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
	function = strings.ToLower(function)
	// Handle different functions
	if function == "cancelorder" {
		return t.CancelOrder(stub, args)
	} else if function == "submitorder" {
		return t.SubmitOrder(stub, args)
	} else if function == "changestatus" {
		return t.ChangeStatus(stub, args)
	} else if function == "editorder" {
		return t.EditOrder(stub, args)
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
	var order Order

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	id = args[0]

	// Get the state from the ledger
	orderBytes, err := stub.GetState(id)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + id + "\"}"
		return nil, errors.New(jsonResp)
	}

	err = json.Unmarshal(orderBytes, &order)

	//if order == nil {
	//	jsonResp := "{\"Error\":\"Nil amount for " + id + "\"}"
	//	return nil, errors.New(jsonResp)
	//}

	jsonResp := "{\"Name\":\"" + id + "\",\"Customer Id\":\"" + order.CustomerId + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return orderBytes, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
