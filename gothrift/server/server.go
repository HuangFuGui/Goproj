package main

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"os"
	"gothrift/com/huangfugui/rpc"
	"sync"
)

const (
	UserServiceAddr = "127.0.0.1:19001"
	UtilServiceAddr = "127.0.0.1:19002"

	USERNAME 	= "黄复贵"
	PASSWORD 	= "123456"

	OK 			= 200
	SUCCEED 	= "login succeed"
	FAILED 		= "username or password error"
)

type UserServiceImpl struct {

}

type UtilServiceImpl struct {

}

// UserService Login服务实现
func (this *UserServiceImpl) Login(callTime int64, cliInfo string, paramMap map[string]string) (r *rpc.Response, err error) {
	// 服务端打印请求日志
	fmt.Printf("go server Login()---cliInfo: %v, callTime: %v ms, params: %v\n", cliInfo, callTime, paramMap)

	username, password := paramMap["username"], paramMap["password"]
	if username == USERNAME && password == PASSWORD {
		return &rpc.Response{OK, SUCCEED}, nil
	}
	return &rpc.Response{OK, FAILED}, nil
}

// UtilService PrimeNumber服务实现
func (this *UtilServiceImpl) PrimeNumber(callTime int64, cliInfo string, threshold int32) (r []int32, err error) {
	// 服务端打印请求日志
	fmt.Printf("go server PrimeNumber()---cliInfo: %v, callTime: %v ms, threshold: %v\n", cliInfo, callTime, threshold)

	prime := primeChan(int(threshold))
	for v := range prime {
		r = append(r, int32(v))
	}
	return
}

func startUserServer() {
	// transport factory定义
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())

	// protocol定义
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	// server transport定义
	serverTransport, err := thrift.NewTServerSocket(UserServiceAddr)
	if err != nil {
		fmt.Println("Error!", err)
		os.Exit(1)
	}

	// UserService processor定义
	userServiceHandler := &UserServiceImpl{}
	processor := rpc.NewUserServiceProcessor(userServiceHandler)

	// server定义（设置上述参数）
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)

	// 开启服务
	fmt.Println("Start userServer on port", UserServiceAddr, "...")
	server.Serve()
}

func startUtilServer() {
	// transport factory定义
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())

	// protocol定义
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	// server transport定义
	serverTransport, err := thrift.NewTServerSocket(UtilServiceAddr)
	if err != nil {
		fmt.Println("Error!", err)
		os.Exit(1)
	}

	// UtilService processor定义
	utilServiceHandler := &UtilServiceImpl{}
	processor := rpc.NewUtilServiceProcessor(utilServiceHandler)

	// server定义（设置上述参数）
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)

	// 开启服务
	fmt.Println("Start utilServer on port", UtilServiceAddr, "...")
	server.Serve()
}

// 同时开启两个服务，一个端口对应一个服务
func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go startUserServer()
	go startUtilServer()
	wg.Wait()
}

//生成器，2，3，4，5，6，7，8....
func generator() chan int {
	ch := make(chan int)
	go func(){
		for i := 2; ; i++ {
			ch <- i
		}
	}()
	return ch
}

//过滤器，将通道能被num整除的数据过滤掉
func filter(ch chan int, num int) chan int{
	out := make(chan int)
	go func(){
		for {
			cur := <- ch
			if cur % num != 0 {
				out <- cur
			}
		}
	}()
	return out
}

//通道里是阈值内的所有素数
func primeChan(threshold int) chan int{
	prime := make(chan int)
	gen := generator()
	num := <- gen
	go func(){
		for num <= threshold {
			prime <- num
			gen = filter(gen, num)
			num = <- gen
		}
		close(prime)
	}()
	return prime
}
