package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/garyburd/redigo/redis"
	"strconv"
	"sync"
	"sync/atomic"
)

var Pool *redis.Pool

func init() {
	redisHost := ":6379"
	Pool = newPool(redisHost)
	close()
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func close() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}

func Get(key string) ([]byte, error) {
	conn := Pool.Get()
	defer conn.Close()

	data, err := redis.Bytes(conn.Do("GET", key))

	return data, err
}

func Set(key, value string) ([]byte, error) {
	conn := Pool.Get()
	defer conn.Close()

	resp, err := redis.Bytes(conn.Do("SET", key, value))

	return resp, err
}

// 流水线：一次性发送多个命令，减少客户端与redis服务器之间的网络通信次数来提升redis在执行多个命令时的性能
// 一般一个命令的执行结果并不会影响另一个命令的输入，且一般不需要receive读取结果
func Pipelining(){
	conn := Pool.Get()
	defer conn.Close()

	//Send writes the command to the connection's output buffer.
	//Flush flushes the connection's output buffer to the server.
	//Receive reads a single reply from the server.
	conn.Send("set", "k1", 100)
	conn.Send("incr", "k1")
	conn.Flush()

	// 一般没必要执行receive方法
	resp, err := conn.Receive()//resp:OK, err:<nil>
	fmt.Printf("resp:%v, err:%v\n", resp, err)

	resp, err = conn.Receive()//resp:101, err:<nil>
	fmt.Printf("resp:%v, err:%v\n", resp, err)
}

// redis的事务（multi、exec）要配合watch、unwatch、discard命令使用，否则没有意义
// 把watch去掉库存会出现负数造成超卖的异常情况
func Transaction(uid int, wg *sync.WaitGroup, casCount *int32) {
	conn := Pool.Get()
	defer conn.Close()

	for {
		// 统计所有用户cas次数
		atomic.AddInt32(casCount, 1)

		// cas开始监听库存
		conn.Do("watch", "inventory")

		// 如果已经没有库存，退出cas抢购
		resp, _ := redis.Bytes(conn.Do("get", "inventory"))
		inv, _ := strconv.Atoi(string(resp))
		if inv <= 0 {
			conn.Do("unwatch")
			fmt.Printf("当前用户（uid：%v），没有库存\n", uid)
			break
		}

		// 发送抢购事务命令
		conn.Send("multi")
		conn.Send("decr", "inventory")
		conn.Send("rpush", "buyers", uid)
		queue, _ := conn.Do("exec")

		// 若queue不为空，cas成功，当前用户不能再抢购，要break
		// 若queue为空，说明在开始监听库存到exec期间，库存被其他用户（协程）抢先修改，要重试
		if queue != nil {
			fmt.Printf("当前用户（uid：%v)，抢购成功，命令结果：%v\n", uid, queue)
			break
		} else {
			fmt.Printf("当前用户（uid：%v)，抢购失败，重试\n", uid)
		}
	}
	wg.Done()
}

func pub(channelKey string) {
	conn := Pool.Get()
	defer conn.Close()

	v := 1
	for {
		conn.Do("publish", channelKey, v)
		time.Sleep(3 * time.Second)
		v += 1
	}
}

func sub(channelKey string) {
	conn := Pool.Get()
	defer conn.Close()

	sub := redis.PubSubConn{conn}
	sub.Subscribe(channelKey)
	for {
		switch v := sub.Receive().(type) {
		case redis.Subscription:
			fmt.Printf("%v %v %v\n", v.Kind, v.Channel, v.Count)
		case redis.Message:
			fmt.Printf("%v %s\n", v.Channel, v.Data)
		case redis.PMessage:
			fmt.Printf("%v %v %v\n", v.Pattern, v.Channel, v.Data)
		case redis.Error:
			fmt.Printf("%v\n", v)
		}
	}
}

func main() {
	//resp, err := Set("key1", "lesliehuang")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("resp:%s\n", resp)
	//
	//data, err := Get("key1")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("data:%s\n", data)
	///*输出：
	//resp:OK
	//data:lesliehuang
	//*/

	//Pipelining()

	//var cliNum int = 10
	//var	casCount int32
	//var wg sync.WaitGroup
	//wg.Add(cliNum)
	//// 模拟n用户，内存cas抢购m个商品（n>m），并记录抢购成功的用户id
	//for uid := 1; uid <= cliNum; uid++ {
	//	go Transaction(uid, &wg, &casCount)
	//}
	//wg.Wait()
	//fmt.Printf("所有用户cas次数：%v\n", casCount)

	go pub("ch1")
	sub("ch1")
}
