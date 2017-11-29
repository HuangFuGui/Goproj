package main

import (
	"net"
	"fmt"
	"encoding/json"
	"gocard/data"
	"strconv"
	"time"
)

func justPrintln(conn net.Conn, sgm *data.Segment) {
	fmt.Println(sgm.Msg)
}

func cliUserLogin(conn net.Conn, sgm *data.Segment) {
	fmt.Println(sgm.Msg)

	var username string
	fmt.Scanln(&username)
	writeBytes, _ := json.Marshal(data.Segment{0, username})
	conn.Write(writeBytes)
}

func cliBuildOrSee(conn net.Conn, sgm *data.Segment) {
	fmt.Println(sgm.Msg)

	var buildOrSee string
	fmt.Scanln(&buildOrSee)
	cmd, _ := strconv.Atoi(buildOrSee)
	writeBytes, _ := json.Marshal(data.Segment{cmd, ""})
	conn.Write(writeBytes)

	if cmd == 2 {
		fmt.Println(data.BUILDSUCCEEDANDSEE)
		time.Sleep(2 * time.Second)
	} else {
		fmt.Println(data.JUSTSEE)
		time.Sleep(2 * time.Second)
	}
}

func cliSeeRoom(conn net.Conn, sgm *data.Segment) {
	roomMap := sgm.Msg.(map[string]interface{})
	for owner, players := range roomMap {
		playerSlice := players.([]interface{})
		if len(playerSlice) == 0 {
			fmt.Printf("%s的房间（%d/%d），当前暂无玩家\n", owner, len(playerSlice), data.ROOMLIMIT)
		} else {
			fmt.Printf("%s的房间（%d/%d），玩家：", owner, len(playerSlice), data.ROOMLIMIT)
		}
		for index, player := range playerSlice {
			if index < len(playerSlice) - 1 {
				fmt.Print(player.(string) + "、")
			} else {
				fmt.Println(player.(string))
			}
		}
	}
	cliChooseRoom(conn, sgm)
}

func cliChooseRoom(conn net.Conn, sgm *data.Segment) {
	fmt.Println(data.CHOOSEROOM)

	var owner string
	fmt.Scanln(&owner)

	//TODO：客户端校验房间人数已满的情况

	writeBytes, _ := json.Marshal(data.Segment{4, owner})
	conn.Write(writeBytes)
}

func cliEnterRoom(conn net.Conn, sgm *data.Segment) {
	if sgm.Msg != "" {
		fmt.Println(sgm.Msg)
		//玩家进入某个房间后，如果房间未满人，应等待其他玩家直到人满后才能开始游戏，开始游戏后，倒计时10s钟等待发牌
		time.Sleep(5 * time.Minute)//模仿在等其他玩家加入
	} else {
		fmt.Println(data.REACHROOMLIMITANDRETRY)
		cliChooseRoom(conn, sgm)
	}
}