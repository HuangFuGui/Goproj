package data

const (
	//tcp每次读的最大字节数
	READLIMIT = 1024
	//每个房间的大小
	ROOMLIMIT = 4
	//游戏倒计时
	COUNTDOWN = 10

	HELLO = "，你好！"
	WELCOME = "欢迎来到gocard！请输入用户名："
	USERNAMEREPEAT = "用户名已存在，请重新输入："
	BUILDORSEE = "创建一个房间请发送2，查看当前所有房间请发送3"
	BUILDSUCCEEDANDSEE = "创建房间成功，正为你列出当前所有房间..."
	JUSTSEE = "正为你列出当前所有房间..."
	CHOOSEROOM = "输入房主的名字，进入房间："
	ENTERROOMANDWAIT = "成功进入房间，还需等待%d名玩家"
	REACHROOMLIMITANDRETRY = "房间人数已满，请重试"
	JOINROOMANDWAIT = "%s加入房间，还需等待%d名玩家"
	STARTGAMECOUNTDOWN = "进入游戏倒计时：%d"
	QUITING = "游戏退出中..."
	TIMEFORMAT = "2006-01-02 15:04:05"
)

//报文段：服务端与客户端之前传送数据的格式（命令，消息）
type Segment struct {
	Cmd int
	Msg	interface{}
}