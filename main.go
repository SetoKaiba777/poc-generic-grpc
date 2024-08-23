package main

import genericclient "grpc-generic-client/generic-client"

func main() {
	c, err :=genericclient.NewClient("localhost:50051","./example.proto")
	if err != nil {
		panic(err)
	}
	input := map[string]interface{}{
		"name": "John Doe",
	}
	err = c.NewRequest("ExampleService","SayHello",input)
	if err != nil {
		panic(err)
	}
	c.PrintResponse()
}