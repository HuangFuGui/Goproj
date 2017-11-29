package main

import (
	"net"
	"gocard/utils"
	"encoding/json"
	"gocard/data"
	"sync"
	"io"
	"github.com/orcaman/concurrent-map"
)


var (
	//userMap记录每个客户端用户名->连接（map[string]net.Conn）的映射
	userMap = cmap.New()
	//roomMap记录每个房间房主->玩家（map[string][]string）的映射
	roomMap = cmap.New()
	//全局锁，使得多个玩家加入同一个房间不会出现数据丢失问题
	mux sync.Mutex
)

func main(){
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":1234")
	utils.CheckError(err)

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	utils.CheckError(err)

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {

	logUserLogin(conn)

	var wg sync.WaitGroup
	wg.Add(1)

	go srvReadWorker(conn, &wg)

	sayWelcome(conn)

	wg.Wait()

	logUserLogout(conn)
}

func srvReadWorker(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	sgmCh := make(chan *data.Segment)
	defer close(sgmCh)

	go srvTaskDispatcher(conn, sgmCh)

	for {
		readBytes := make([]byte, data.READLIMIT)
		readLen, err := conn.Read(readBytes)
		//客户端主动关闭连接
		if err == io.EOF {
			srvLogout(conn.RemoteAddr().String())
			break
		}

		var sgm data.Segment
		json.Unmarshal(readBytes[:readLen], &sgm)

		sgmCh <- &sgm
	}
}

func srvTaskDispatcher(conn net.Conn, sgmCh <-chan *data.Segment) {
	for sgm := range sgmCh {
		if sgm.Cmd == 0 {
			srvUserLogin(conn, sgm)
		}
		if sgm.Cmd == 2 {
			srvBuildRoom(conn, sgm)
		}
		if sgm.Cmd == 3 {
			srvSeeRoom(conn, sgm)
		}
		if sgm.Cmd == 4 {
			srvChooseRoom(conn, sgm)
		}
	}
}