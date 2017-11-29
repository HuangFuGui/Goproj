package main

import (
	"gothrift/com/huangfugui/rpc"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"net"
	"os"
	"time"
)

func main() {
	testLogin()
	testPrimeNumber()
}

func testLogin() {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transport, err := thrift.NewTSocket(net.JoinHostPort("127.0.0.1", "19000"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		os.Exit(1)
	}
	useTransport := transportFactory.GetTransport(transport)
	client := rpc.NewUserServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to 127.0.0.1:19000", " ", err)
		os.Exit(1)
	}
	defer transport.Close()

	paramMap := make(map[string]string)
	paramMap["username"], paramMap["password"] = "黄复贵", "123456"
	fmt.Println("form:", paramMap)

	startTime := currentTimeMillis()
	response, err := client.Login(startTime, "go client request Login", paramMap)
	endTime := currentTimeMillis()

	fmt.Println(response.Msg)
	fmt.Printf("Program exit, startTime: %v ms, endTime: %v ms, totalCost: %v ms\n", startTime, endTime, endTime - startTime)
}

func testPrimeNumber() {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transport, err := thrift.NewTSocket(net.JoinHostPort("127.0.0.1", "19002"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		os.Exit(1)
	}
	useTransport := transportFactory.GetTransport(transport)
	client := rpc.NewUtilServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to 127.0.0.1:19002", " ", err)
		os.Exit(1)
	}
	defer transport.Close()

	threshold := int32(100)
	fmt.Println("threshold:", threshold)

	startTime := currentTimeMillis()
	r, err := client.PrimeNumber(startTime, "go client request PrimeNumber", threshold)
	endTime := currentTimeMillis()
	fmt.Printf("Program exit, startTime: %v ms, endTime: %v ms, prime number within %v: %v", startTime, endTime, threshold, r)
}
// 纳秒转换成毫秒
func currentTimeMillis() int64 {
	return time.Now().UnixNano() / 1e6
}