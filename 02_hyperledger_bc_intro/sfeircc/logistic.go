package main

import (
	"fmt"
	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type OrderContract struct {

}

type order struct {
	OrderId string //used
	ShopId string //used
	WarehouseId string //used
	TransporterId string //used
	Status string // used
	ReferenceId string // used	
	UnitPrice int //used
	QuantityPrice int //used
}

func (t *OrderContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}


func (t *OrderContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "order" {
		return t.order(stub, args)
	} else if function == "book" {
		return t.book(stub, args)
	} else if function == "deliver" {
		return t.deliver(stub, args)
	}else if function == "read" {
		return t.read(stub, args)
	}

	return shim.Error("Invalid function name")
}

func (t *OrderContract) order(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	orderId := args[0]
	shopId := args[1]
	referenceId := args[2]
	quantity,_ := strconv.Atoi(args[3])
	price,_ := strconv.Atoi(args[4])
		
	orderContract := order {
		OrderId: orderId, 
		ShopId: shopId, 
		ReferenceId: referenceId,
		UnitPrice: price,
		QuantityPrice: quantity,
		Status: "ORDERED"}

	tcBytes, _ := json.Marshal(orderContract)
	stub.PutState(orderContract.OrderId, tcBytes)
	return shim.Success(nil)
}

func (t *OrderContract) book(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	orderId := args[0]
	warehouseId := args[1]
	orBytes, _ := stub.GetState(orderId)
	or := order{}
	json.Unmarshal(orBytes, &or)

	if (or.Status == "ORDERED") {
		or.WarehouseId = warehouseId
		or.Status = "BOOKED"
	} else {
		fmt.Printf("Order not initiated yet")
	}

	orBytes, _ = json.Marshal(or)
	stub.PutState(orderId, orBytes)
	
	return shim.Success(nil)
}

func (t *OrderContract) deliver(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	orderId := args[0]
	transporterId := args[1]
	orBytes, _ := stub.GetState(orderId)
	or := order{}
	json.Unmarshal(orBytes, &or)

	if (or.Status == "BOOKED") {
		or.TransporterId = transporterId
		or.Status = "DELIVERED"
	} else {
		fmt.Printf("Order not initiated yet")
	}

	orBytes, _ = json.Marshal(or)
	stub.PutState(orderId, orBytes)
	
	return shim.Success(nil)
}

func (t *OrderContract) read(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("Query from chaincode ...")
	var A string // Entities
	var err error
	

	fmt.Printf("Query arguments: "+args[0])

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil trade for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(Avalbytes)
}

func main() {

	err := shim.Start(new(OrderContract))
	if err != nil {
		fmt.Printf("Error creating new Trade Contract: %s", err)
	}
}

