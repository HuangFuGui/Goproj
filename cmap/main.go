package main
import (
	"fmt"
	"github.com/orcaman/concurrent-map"
	"time"
)

func main(){
	//程序panic崩溃，fatal error: concurrent map read and map write
	//m := make(map[string]int)
	//m["key1"] = 1
	//m["key2"] = 2
	//go func() {
	//	for {
	//		_ = m["key1"]
	//	}
	//}()
	//go func() {
	//	for {
	//		m["key2"] = 2
	//	}
	//}()
	//time.Sleep(10 * time.Second)
	//fmt.Println("finished")

	//finished，并发安全，concurrent-map的源码值得学习
	m := cmap.New()
	m.Set("key1", 1)
	m.Set("key2", 2)
	go func() {
		for {
			if tmp, ok := m.Get("key1"); ok {
				_ = tmp.(int)
			}
		}
	}()
	go func() {
		for {
			m.Set("key2", 2)
		}
	}()
	time.Sleep(10 * time.Second)
	fmt.Println("finished")
}