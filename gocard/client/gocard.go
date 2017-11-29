package main

import (
	"net"
	"gocard/utils"
	"gocard/data"
	"encoding/json"
	"sync"
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"
)

var wg sync.WaitGroup

func main(){
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":1234")
	utils.CheckError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	utils.CheckError(err)

	go cliReadWorker(conn)

	wg.Add(1)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)
		<- signals
		wg.Done()
	}()

	wg.Wait()
	conn.Close()
	fmt.Println(data.QUITING)
	time.Sleep(2 * time.Second)
	os.Exit(0)//程序正常退出
}

func cliReadWorker(conn net.Conn) {
	sgmCh := make(chan *data.Segment)
	defer close(sgmCh)

	go cliTaskDispatcher(conn, sgmCh)

	for {
		readBytes := make([]byte, data.READLIMIT)
		readLen, err := conn.Read(readBytes)
		if err != nil {
			break
		}

		var sgm data.Segment
		json.Unmarshal(readBytes[:readLen], &sgm)

		sgmCh <- &sgm
	}
}

func cliTaskDispatcher(conn net.Conn, sgmCh <-chan *data.Segment) {
	for sgm := range sgmCh {
		if sgm.Cmd == -1 {
			justPrintln(conn, sgm)
		}
		if sgm.Cmd == 0 {
			cliUserLogin(conn, sgm)
		}
		if sgm.Cmd == 1 {
			cliBuildOrSee(conn, sgm)
		}
		if sgm.Cmd == 3 {
			cliSeeRoom(conn, sgm)
		}
		if sgm.Cmd == 4 {
			go cliEnterRoom(conn, sgm)
		}
	}
}