package main

import (
	"gocard/data"
	"net"
	"encoding/json"
	"fmt"
	"time"
)

func sayWelcome(conn net.Conn) {
	//欢迎信息，提示登录
	writeBytes, _ := json.Marshal(data.Segment{0, data.WELCOME})
	conn.Write(writeBytes)
}

func srvUserLogin(conn net.Conn, sgm *data.Segment) {
	username := sgm.Msg.(string)
	if userMap.Has(username) {
		writeBytes, _ := json.Marshal(data.Segment{0, data.USERNAMEREPEAT})
		conn.Write(writeBytes)
		return
	}
	//你好
	userMap.Set(username, conn)
	writeBytes, _ := json.Marshal(data.Segment{-1, username + data.HELLO})
	conn.Write(writeBytes)
	//提示创建房间还是查看所有房间
	writeBytes, _ = json.Marshal(data.Segment{1, data.BUILDORSEE})
	conn.Write(writeBytes)
}

func srvBuildRoom(conn net.Conn, sgm *data.Segment) {
	username := getUserByAddr(conn.RemoteAddr().String())
	roomMap.Set(username, make([]string, 0, data.ROOMLIMIT))
	//创建房间后输出所有房间信息
	srvSeeRoom(conn, sgm)
}

func srvSeeRoom(conn net.Conn, sgm *data.Segment) {
	roomMapTmp := roomMap.Items()
	writeBytes, _ := json.Marshal(data.Segment{3, roomMapTmp})
	conn.Write(writeBytes)
}

func srvChooseRoom(conn net.Conn, sgm *data.Segment) {
	owner, username, writesBytes, msg := sgm.Msg.(string), getUserByAddr(conn.RemoteAddr().String()), []byte(nil), ""
	mux.Lock()
	players, _ := roomMap.Get(owner)
	playerSlice, _ := players.([]string)
	if len(playerSlice) < data.ROOMLIMIT {
		//房间玩家未满，注意从map中读取出来的值是slice引用的副本（即：3域结构体）
		//若append，会生成一个新的slice引用，要把append后的新slice引用赋值给map才符合场景
		playerSlice = append(playerSlice, username)
		roomMap.Set(owner, playerSlice)
		mux.Unlock()
		msg = fmt.Sprintf(data.ENTERROOMANDWAIT, data.ROOMLIMIT - len(playerSlice))
		//提示当前有玩家进入房间
		for i, v := range playerSlice {
			if i < len(playerSlice) - 1 {
				joinWait := fmt.Sprintf(data.JOINROOMANDWAIT, username, data.ROOMLIMIT - len(playerSlice))
				writesBytes, _ = json.Marshal(data.Segment{-1, joinWait})
				conn, _ := userMap.Get(v)
				conn.(net.Conn).Write(writesBytes)
			}
		}
		writesBytes, _ = json.Marshal(data.Segment{4, msg})
		conn.Write(writesBytes)
		//游戏倒计时
		if data.ROOMLIMIT == len(playerSlice) {
			go func() {
				time.Sleep(time.Second)
				for i := data.COUNTDOWN; i > 0; i-- {
					writesBytes, _ = json.Marshal(data.Segment{-1, fmt.Sprintf(data.STARTGAMECOUNTDOWN, i)})
					time.Sleep(time.Second)
					for _, v := range playerSlice {
						conn, _ := userMap.Get(v)
						conn.(net.Conn).Write(writesBytes)
					}
				}
			}()
		}
	} else {
		//房间玩家已满
		mux.Unlock()
		writesBytes, _ = json.Marshal(data.Segment{4, msg})
		conn.Write(writesBytes)
	}
}

func srvLogout(addr string) {
	tmp := userMap.Items()
	for k, v := range tmp {
		if v.(net.Conn).RemoteAddr().String() == addr {
			userMap.Remove(k)
		}
	}
}


/***********************************************************************************************************************
上述为业务代码，下面为服务端的辅助函数
***********************************************************************************************************************/

//通过地址找到用户名
func getUserByAddr(addr string) string {
	tmp := userMap.Items()
	for k, v := range tmp {
		if v.(net.Conn).RemoteAddr().String() == addr {
			return k
		}
	}
	//编译需要，一般不会return ""
	return ""
}

//服务端log用户登录日志
func logUserLogin(conn net.Conn) {
	fmt.Println("[" + time.Now().Format(data.TIMEFORMAT) + "]" + "客户端" + conn.RemoteAddr().String() + "登录游戏")
}

//服务端log用户登出日志
func logUserLogout(conn net.Conn) {
	fmt.Println("[" + time.Now().Format(data.TIMEFORMAT) + "]" + "客户端" + conn.RemoteAddr().String() + "退出游戏")
}